package cli

import (
	"bytes"
	"sync"

	"dappco.re/go"
	"time"
)

func restoreThemeAndColors(t *core.T) {
	t.Helper()

	prevTheme := currentTheme
	prevColor := ColorEnabled()
	t.Cleanup(func() {
		currentTheme = prevTheme
		SetColorEnabled(prevColor)
	})
}

func TestTaskTracker_Good(t *core.T) {
	t.Run("add and complete tasks", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{} // non-TTY

		t1 := tr.Add("repo-a")
		t2 := tr.Add("repo-b")

		t1.Update("pulling...")
		t2.Update("pulling...")

		t1.Done("up to date")
		t2.Done("3 commits behind")

		out := tr.String()
		core.AssertContains(t, out, "repo-a")
		core.AssertContains(t, out, "repo-b")
		core.AssertContains(t, out, "up to date")
		core.AssertContains(t, out, "3 commits behind")
	})

	t.Run("task states", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		task := tr.Add("test")

		// Pending
		_, _, state := task.snapshot()
		core.AssertEqual(t, taskPending, state)

		// Running
		task.Update("working")
		_, status, state := task.snapshot()
		core.AssertEqual(t, taskRunning, state)
		core.AssertEqual(t, "working", status)

		// Done
		task.Done("finished")
		_, status, state = task.snapshot()
		core.AssertEqual(t, taskDone, state)
		core.AssertEqual(t, "finished", status)
	})

	t.Run("task fail", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		task := tr.Add("bad-repo")
		task.Fail("connection refused")

		_, status, state := task.snapshot()
		core.AssertEqual(t, taskFailed, state)
		core.AssertEqual(t, "connection refused", status)
	})

	t.Run("concurrent updates", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		var wg sync.WaitGroup
		for i := range 10 {
			task := tr.Add("task-" + string(rune('a'+i)))
			wg.Add(1)
			go func(t *TrackedTask) {
				defer wg.Done()
				t.Update("running")
				time.Sleep(5 * time.Millisecond)
				t.Done("ok")
			}(task)
		}
		wg.Wait()
		core.AssertTrue(t, tr.allDone())
	})

	t.Run("summary all passed", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("a").Done("ok")
		tr.Add("b").Done("ok")
		tr.Add("c").Done("ok")
		core.AssertEqual(t, "3/3 passed", tr.Summary())
	})

	t.Run("summary with failures", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("a").Done("ok")
		tr.Add("b").Fail("error")
		tr.Add("c").Done("ok")
		core.AssertEqual(t, "2/3 passed, 1 failed", tr.Summary())
	})

	t.Run("wait completes for non-TTY", func(t *core.T) {
		var buf bytes.Buffer
		tr := NewTaskTracker().WithOutput(&buf)

		task := tr.Add("quick")
		go func() {
			time.Sleep(10 * time.Millisecond)
			task.Done("done")
		}()

		tr.Wait()
		core.AssertContains(t, buf.String(), "quick")
		core.AssertContains(t, buf.String(), "done")
	})

	t.Run("WithOutput sets output writer", func(t *core.T) {
		var buf bytes.Buffer
		tr := NewTaskTracker().WithOutput(&buf)

		tr.Add("quick").Done("done")
		tr.Wait()
		core.AssertContains(t, buf.String(), "quick")
		core.AssertContains(t, buf.String(), "done")
	})

	t.Run("name width alignment", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("short")
		tr.Add("very-long-repo-name")

		w := tr.nameWidth()
		core.AssertEqual(t, 19, w)
	})

	t.Run("name width counts visible width", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("東京")
		tr.Add("repo")

		w := tr.nameWidth()
		core.AssertEqual(t, 4, w)
	})

	t.Run("String output format", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("repo-a").Done("clean")
		tr.Add("repo-b").Fail("dirty")
		tr.Add("repo-c").Update("pulling")

		out := tr.String()
		core.AssertContains(t, out, "✓")
		core.AssertContains(t, out, "✗")
		core.AssertContains(t, out, "⠋")
	})

	t.Run("glyph shortcodes render in names and statuses", func(t *core.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add(":check: repo").Done("done :warn:")

		out := tr.String()
		core.AssertContains(t, out, "[OK] repo")
		core.AssertContains(t, out, "[WARN]")
	})

	t.Run("ASCII theme uses ASCII symbols", func(t *core.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("repo-a").Done("clean")
		tr.Add("repo-b").Fail("dirty")
		tr.Add("repo-c").Update("pulling")

		out := tr.String()
		core.AssertContains(t, out, "[OK]")
		core.AssertContains(t, out, "[FAIL]")
		core.AssertContains(t, out, "-")
		core.AssertNotContains(t, out, "✓")
		core.AssertNotContains(t, out, "✗")
	})

	t.Run("iterators tolerate mutation during iteration", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("first")
		tr.Add("second")

		done := make(chan struct{})
		go func() {
			defer close(done)
			for task := range tr.Tasks() {
				task.Update("visited")
			}
		}()

		// TODO(assertx): no core/go equivalent for require.Eventually(t, func() bool, time.Second, 10*time.Millisecond).
		eventuallyDone := false
		deadline := time.Now().Add(time.Second)
		for time.Now().Before(deadline) {
			select {
			case <-done:
				eventuallyDone = true
			default:
			}
			if eventuallyDone {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !eventuallyDone {
			select {
			case <-done:
				eventuallyDone = true
			default:
			}
		}
		core.AssertTrue(t, eventuallyDone, "iterator did not finish before timeout")

		for name, status := range tr.Snapshots() {
			core.AssertEqual(t, "visited", status, name)
		}
	})
}

func TestTaskTracker_Bad(t *core.T) {
	t.Run("allDone with no tasks", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		core.AssertTrue(t, tr.allDone())
	})

	t.Run("allDone incomplete", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		tr.Add("pending")
		core.AssertFalse(t, tr.allDone())
	})
}

func TestTrackedTask_Good(t *core.T) {
	t.Run("thread safety", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		task := tr.Add("concurrent")

		var wg sync.WaitGroup
		for range 100 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				task.Update("running")
			}()
		}
		wg.Wait()

		_, status, state := task.snapshot()
		core.RequireTrue(t, core.DeepEqual(taskRunning, state), core.Sprintf("want=%#v got=%#v", taskRunning, state))
		core.RequireTrue(t, core.DeepEqual("running", status), core.Sprintf("want=%#v got=%#v", "running", status))
	})
}

func TestTaskTracker_Ugly(t *core.T) {
	t.Run("empty task name does not panic", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		core.AssertNotPanics(t, func() {
			task := tr.Add("")
			task.Done("ok")
		})
	})

	t.Run("Done called twice does not panic", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		task := tr.Add("double-done")
		core.AssertNotPanics(t, func() {
			task.Done("first")
			task.Done("second")
		})
	})

	t.Run("Fail after Done does not panic", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		task := tr.Add("already-done")
		core.AssertNotPanics(t, func() {
			task.Done("completed")
			task.Fail("too late")
		})
	})

	t.Run("String on empty tracker does not panic", func(t *core.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		core.AssertNotPanics(t, func() {
			_ = tr.String()
		})
	})
}
