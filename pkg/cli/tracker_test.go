package cli

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func restoreThemeAndColors(t *testing.T) {
	t.Helper()

	prevTheme := currentTheme
	prevColor := ColorEnabled()
	t.Cleanup(func() {
		currentTheme = prevTheme
		SetColorEnabled(prevColor)
	})
}

func TestTaskTracker_Good(t *testing.T) {
	t.Run("add and complete tasks", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{} // non-TTY

		t1 := tr.Add("repo-a")
		t2 := tr.Add("repo-b")

		t1.Update("pulling...")
		t2.Update("pulling...")

		t1.Done("up to date")
		t2.Done("3 commits behind")

		out := tr.String()
		assert.Contains(t, out, "repo-a")
		assert.Contains(t, out, "repo-b")
		assert.Contains(t, out, "up to date")
		assert.Contains(t, out, "3 commits behind")
	})

	t.Run("task states", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		task := tr.Add("test")

		// Pending
		_, _, state := task.snapshot()
		assert.Equal(t, taskPending, state)

		// Running
		task.Update("working")
		_, status, state := task.snapshot()
		assert.Equal(t, taskRunning, state)
		assert.Equal(t, "working", status)

		// Done
		task.Done("finished")
		_, status, state = task.snapshot()
		assert.Equal(t, taskDone, state)
		assert.Equal(t, "finished", status)
	})

	t.Run("task fail", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		task := tr.Add("bad-repo")
		task.Fail("connection refused")

		_, status, state := task.snapshot()
		assert.Equal(t, taskFailed, state)
		assert.Equal(t, "connection refused", status)
	})

	t.Run("concurrent updates", func(t *testing.T) {
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

		assert.True(t, tr.allDone())
	})

	t.Run("summary all passed", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("a").Done("ok")
		tr.Add("b").Done("ok")
		tr.Add("c").Done("ok")

		assert.Equal(t, "3/3 passed", tr.Summary())
	})

	t.Run("summary with failures", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("a").Done("ok")
		tr.Add("b").Fail("error")
		tr.Add("c").Done("ok")

		assert.Equal(t, "2/3 passed, 1 failed", tr.Summary())
	})

	t.Run("wait completes for non-TTY", func(t *testing.T) {
		var buf bytes.Buffer
		tr := NewTaskTracker().WithOutput(&buf)

		task := tr.Add("quick")
		go func() {
			time.Sleep(10 * time.Millisecond)
			task.Done("done")
		}()

		tr.Wait()
		assert.Contains(t, buf.String(), "quick")
		assert.Contains(t, buf.String(), "done")
	})

	t.Run("WithOutput sets output writer", func(t *testing.T) {
		var buf bytes.Buffer
		tr := NewTaskTracker().WithOutput(&buf)

		tr.Add("quick").Done("done")
		tr.Wait()

		assert.Contains(t, buf.String(), "quick")
		assert.Contains(t, buf.String(), "done")
	})

	t.Run("name width alignment", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("short")
		tr.Add("very-long-repo-name")

		w := tr.nameWidth()
		assert.Equal(t, 19, w)
	})

	t.Run("name width counts visible width", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("東京")
		tr.Add("repo")

		w := tr.nameWidth()
		assert.Equal(t, 4, w)
	})

	t.Run("String output format", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("repo-a").Done("clean")
		tr.Add("repo-b").Fail("dirty")
		tr.Add("repo-c").Update("pulling")

		out := tr.String()
		assert.Contains(t, out, "✓")
		assert.Contains(t, out, "✗")
		assert.Contains(t, out, "⠋")
	})

	t.Run("glyph shortcodes render in names and statuses", func(t *testing.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add(":check: repo").Done("done :warn:")

		out := tr.String()
		assert.Contains(t, out, "[OK] repo")
		assert.Contains(t, out, "[WARN]")
	})

	t.Run("ASCII theme uses ASCII symbols", func(t *testing.T) {
		restoreThemeAndColors(t)
		UseASCII()

		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}

		tr.Add("repo-a").Done("clean")
		tr.Add("repo-b").Fail("dirty")
		tr.Add("repo-c").Update("pulling")

		out := tr.String()
		assert.Contains(t, out, "[OK]")
		assert.Contains(t, out, "[FAIL]")
		assert.Contains(t, out, "-")
		assert.NotContains(t, out, "✓")
		assert.NotContains(t, out, "✗")
	})

	t.Run("iterators tolerate mutation during iteration", func(t *testing.T) {
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

		require.Eventually(t, func() bool {
			select {
			case <-done:
				return true
			default:
				return false
			}
		}, time.Second, 10*time.Millisecond)

		for name, status := range tr.Snapshots() {
			assert.Equal(t, "visited", status, name)
		}
	})
}

func TestTaskTracker_Bad(t *testing.T) {
	t.Run("allDone with no tasks", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		assert.True(t, tr.allDone())
	})

	t.Run("allDone incomplete", func(t *testing.T) {
		tr := NewTaskTracker()
		tr.out = &bytes.Buffer{}
		tr.Add("pending")
		assert.False(t, tr.allDone())
	})
}

func TestTrackedTask_Good(t *testing.T) {
	t.Run("thread safety", func(t *testing.T) {
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
		require.Equal(t, taskRunning, state)
		require.Equal(t, "running", status)
	})
}
