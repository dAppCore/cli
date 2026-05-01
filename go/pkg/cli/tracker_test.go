package cli

import (
	core "dappco.re/go"
)

func TestTracker_TrackedTask_Update_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Update("running")

	core.AssertContains(t, task.tracker.String(), "running")
	core.AssertContains(t, task.tracker.Summary(), "0/1")
}

func TestTracker_TrackedTask_Update_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Update("running") })
	core.AssertNil(t, task)
}

func TestTracker_TrackedTask_Update_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Update("")

	core.AssertContains(t, task.tracker.String(), "build")
	core.AssertContains(t, task.tracker.Summary(), "0/1")
}

func TestTracker_TrackedTask_Done_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Done("done")

	core.AssertContains(t, task.tracker.String(), "done")
	core.AssertEqual(t, "1/1 passed", task.tracker.Summary())
}

func TestTracker_TrackedTask_Done_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Done("done") })
	core.AssertNil(t, task)
}

func TestTracker_TrackedTask_Done_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Done("")

	core.AssertEqual(t, "1/1 passed", task.tracker.Summary())
	core.AssertContains(t, task.tracker.String(), "build")
}

func TestTracker_TrackedTask_Fail_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Fail("failed")

	core.AssertContains(t, task.tracker.String(), "failed")
	core.AssertEqual(t, "0/1 passed, 1 failed", task.tracker.Summary())
}

func TestTracker_TrackedTask_Fail_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Fail("failed") })
	core.AssertNil(t, task)
}

func TestTracker_TrackedTask_Fail_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Fail("")

	core.AssertEqual(t, "0/1 passed, 1 failed", task.tracker.Summary())
	core.AssertContains(t, task.tracker.String(), "build")
}

func TestTracker_NewTaskTracker_Good(t *core.T) {
	tr := NewTaskTracker()

	core.AssertNotNil(t, tr)
	core.AssertNotNil(t, tr.out)
}

func TestTracker_NewTaskTracker_Bad(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEmpty(t, tr.tasks)
	core.AssertFalse(t, tr.started)
}

func TestTracker_NewTaskTracker_Ugly(t *core.T) {
	tr := NewTaskTracker().WithOutput(core.NewBuilder())

	core.AssertNotNil(t, tr.out)
	core.AssertEqual(t, "", tr.String())
}

func TestTracker_TaskTracker_WithOutput_Good(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)

	core.AssertEqual(t, out, tr.out)
	core.AssertEqual(t, tr, tr.WithOutput(out))
}

func TestTracker_TaskTracker_WithOutput_Bad(t *core.T) {
	tr := NewTaskTracker()
	original := tr.out

	core.AssertEqual(t, tr, tr.WithOutput(nil))
	core.AssertEqual(t, original, tr.out)
}

func TestTracker_TaskTracker_WithOutput_Ugly(t *core.T) {
	tr := NewTaskTracker().WithOutput(core.Discard)

	core.AssertEqual(t, core.Discard, tr.out)
	core.AssertFalse(t, tr.isTTY())
}

func TestTracker_TaskTracker_Add_Good(t *core.T) {
	tr := NewTaskTracker()
	task := tr.Add("build")

	core.AssertNotNil(t, task)
	core.AssertContains(t, tr.String(), "build")
}

func TestTracker_TaskTracker_Add_Bad(t *core.T) {
	tr := NewTaskTracker()
	task := tr.Add("")

	core.AssertNotNil(t, task)
	core.AssertLen(t, tr.tasks, 1)
}

func TestTracker_TaskTracker_Add_Ugly(t *core.T) {
	tr := NewTaskTracker()
	first := tr.Add("same")
	second := tr.Add("same")

	core.AssertTrue(t, first != second)
	core.AssertLen(t, tr.tasks, 2)
}

func TestTracker_TaskTracker_Wait_Good(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)
	tr.Add("build").Done("done")

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertContains(t, out.String(), "done")
}

func TestTracker_TaskTracker_Wait_Bad(t *core.T) {
	tr := NewTaskTracker().WithOutput(core.NewBuilder())

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertEqual(t, "0/0 passed", tr.Summary())
}

func TestTracker_TaskTracker_Wait_Ugly(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)
	tr.Add("build").Fail("failed")

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertContains(t, out.String(), "failed")
}

func TestTracker_TaskTracker_Tasks_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build")
	var names []string
	for task := range tr.Tasks() {
		names = append(names, task.name)
	}

	core.AssertEqual(t, []string{"build"}, names)
}

func TestTracker_TaskTracker_Tasks_Bad(t *core.T) {
	tr := NewTaskTracker()
	var count int
	for range tr.Tasks() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestTracker_TaskTracker_Tasks_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("first")
	tr.Add("second")
	var count int
	for range tr.Tasks() {
		count++
	}

	core.AssertEqual(t, 2, count)
}

func TestTracker_TaskTracker_Snapshots_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build").Update("running")
	var got []string
	for name, status := range tr.Snapshots() {
		got = append(got, name+":"+status)
	}

	core.AssertEqual(t, []string{"build:running"}, got)
}

func TestTracker_TaskTracker_Snapshots_Bad(t *core.T) {
	tr := NewTaskTracker()
	var count int
	for range tr.Snapshots() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestTracker_TaskTracker_Snapshots_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("").Done("")
	var got []string
	for name, status := range tr.Snapshots() {
		got = append(got, name+":"+status)
	}

	core.AssertEqual(t, []string{":"}, got)
}

func TestTracker_TaskTracker_Summary_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("ok").Done("done")

	core.AssertEqual(t, "1/1 passed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "passed")
}

func TestTracker_TaskTracker_Summary_Bad(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("bad").Fail("failed")

	core.AssertEqual(t, "0/1 passed, 1 failed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "failed")
}

func TestTracker_TaskTracker_Summary_Ugly(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEqual(t, "0/0 passed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "0/0")
}

func TestTracker_TaskTracker_String_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build").Done("done")

	core.AssertContains(t, tr.String(), "build")
	core.AssertContains(t, tr.String(), "done")
}

func TestTracker_TaskTracker_String_Bad(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEqual(t, "", tr.String())
	core.AssertEmpty(t, tr.String())
}

func TestTracker_TaskTracker_String_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add(":check:").Update(":warn:")

	core.AssertContains(t, tr.String(), "✓")
	core.AssertContains(t, tr.String(), "⚠")
}
