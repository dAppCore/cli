package cli

import (
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// Spinner frames for the live tracker.
var spinnerFramesUnicode = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var spinnerFramesASCII = []string{"-", "\\", "|", "/"}

// taskState tracks the lifecycle of a tracked task.
type taskState int

const (
	taskPending taskState = iota
	taskRunning
	taskDone
	taskFailed
)

// TrackedTask represents a single task in a TaskTracker.
// Safe for concurrent use — call Update, Done, or Fail from any goroutine.
type TrackedTask struct {
	name    string
	status  string
	state   taskState
	tracker *TaskTracker
	mu      sync.Mutex
}

// Update sets the task status message and marks it as running.
func (t *TrackedTask) Update(status string) {
	t.mu.Lock()
	t.status = status
	t.state = taskRunning
	t.mu.Unlock()
}

// Done marks the task as successfully completed with a final message.
func (t *TrackedTask) Done(message string) {
	t.mu.Lock()
	t.status = message
	t.state = taskDone
	t.mu.Unlock()
}

// Fail marks the task as failed with an error message.
func (t *TrackedTask) Fail(message string) {
	t.mu.Lock()
	t.status = message
	t.state = taskFailed
	t.mu.Unlock()
}

func (t *TrackedTask) snapshot() (string, string, taskState) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.name, t.status, t.state
}

// TaskTracker displays multiple concurrent tasks with individual spinners.
//
//	tracker := cli.NewTaskTracker()
//	for _, repo := range repos {
//	    t := tracker.Add(repo.Name)
//	    go func(t *cli.TrackedTask) {
//	        t.Update("pulling...")
//	        // ...
//	        t.Done("up to date")
//	    }(t)
//	}
//	tracker.Wait()
type TaskTracker struct {
	tasks   []*TrackedTask
	out     io.Writer
	mu      sync.Mutex
	started bool
}

// Tasks returns an iterator over the tasks in the tracker.
func (tr *TaskTracker) Tasks() iter.Seq[*TrackedTask] {
	return func(yield func(*TrackedTask) bool) {
		tr.mu.Lock()
		tasks := make([]*TrackedTask, len(tr.tasks))
		copy(tasks, tr.tasks)
		tr.mu.Unlock()

		for _, t := range tasks {
			if !yield(t) {
				return
			}
		}
	}
}

// Snapshots returns an iterator over snapshots of tasks in the tracker.
func (tr *TaskTracker) Snapshots() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		tr.mu.Lock()
		tasks := make([]*TrackedTask, len(tr.tasks))
		copy(tasks, tr.tasks)
		tr.mu.Unlock()

		for _, t := range tasks {
			name, status, _ := t.snapshot()
			if !yield(name, status) {
				return
			}
		}
	}
}

// NewTaskTracker creates a new parallel task tracker.
func NewTaskTracker() *TaskTracker {
	return &TaskTracker{out: stderrWriter()}
}

// WithOutput sets the destination writer for tracker output.
// Pass nil to keep the current writer unchanged.
func (tr *TaskTracker) WithOutput(out io.Writer) *TaskTracker {
	if out != nil {
		tr.out = out
	}
	return tr
}

// Add registers a task and returns it for goroutine use.
func (tr *TaskTracker) Add(name string) *TrackedTask {
	t := &TrackedTask{
		name:    name,
		status:  "waiting",
		state:   taskPending,
		tracker: tr,
	}
	tr.mu.Lock()
	tr.tasks = append(tr.tasks, t)
	tr.mu.Unlock()
	return t
}

// Wait renders the task display and blocks until all tasks complete.
// Uses ANSI cursor manipulation for live updates when connected to a terminal.
// Falls back to line-by-line output for non-TTY.
func (tr *TaskTracker) Wait() {
	if !tr.isTTY() {
		tr.waitStatic()
		return
	}
	tr.waitLive()
}

func (tr *TaskTracker) isTTY() bool {
	if f, ok := tr.out.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

func (tr *TaskTracker) waitStatic() {
	// Non-TTY: print final state of each task when it completes.
	reported := make(map[int]bool)
	for {
		tr.mu.Lock()
		tasks := tr.tasks
		tr.mu.Unlock()

		allDone := true
		for i, t := range tasks {
			name, status, state := t.snapshot()
			name = compileGlyphs(name)
			status = compileGlyphs(status)
			if state != taskDone && state != taskFailed {
				allDone = false
				continue
			}
			if reported[i] {
				continue
			}
			reported[i] = true
			icon := Glyph(":check:")
			if state == taskFailed {
				icon = Glyph(":cross:")
			}
			fmt.Fprintf(tr.out, "%s %-20s %s\n", icon, name, status)
		}
		if allDone {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (tr *TaskTracker) waitLive() {
	tr.mu.Lock()
	n := len(tr.tasks)
	tr.mu.Unlock()

	// Print initial lines.
	frame := 0
	for i := range n {
		tr.renderLine(i, frame)
	}

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		frame++

		tr.mu.Lock()
		count := len(tr.tasks)
		tr.mu.Unlock()

		// Move cursor up to redraw all lines.
		fmt.Fprintf(tr.out, "\033[%dA", count)
		for i := range count {
			tr.renderLine(i, frame)
		}

		if tr.allDone() {
			return
		}
	}
}

func (tr *TaskTracker) renderLine(idx, frame int) {
	tr.mu.Lock()
	t := tr.tasks[idx]
	tr.mu.Unlock()

	name, status, state := t.snapshot()
	name = compileGlyphs(name)
	status = compileGlyphs(status)
	nameW := tr.nameWidth()

	var icon string
	switch state {
	case taskPending:
		icon = DimStyle.Render(Glyph(":pending:"))
	case taskRunning:
		icon = InfoStyle.Render(trackerSpinnerFrame(frame))
	case taskDone:
		icon = SuccessStyle.Render(Glyph(":check:"))
	case taskFailed:
		icon = ErrorStyle.Render(Glyph(":cross:"))
	}

	var styledStatus string
	switch state {
	case taskDone:
		styledStatus = SuccessStyle.Render(status)
	case taskFailed:
		styledStatus = ErrorStyle.Render(status)
	default:
		styledStatus = DimStyle.Render(status)
	}

	fmt.Fprintf(tr.out, "\033[2K%s %s %s\n", icon, Pad(name, nameW), styledStatus)
}

func (tr *TaskTracker) nameWidth() int {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	w := 0
	for _, t := range tr.tasks {
		if nameW := displayWidth(compileGlyphs(t.name)); nameW > w {
			w = nameW
		}
	}
	return w
}

func (tr *TaskTracker) allDone() bool {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	for _, t := range tr.tasks {
		_, _, state := t.snapshot()
		if state != taskDone && state != taskFailed {
			return false
		}
	}
	return true
}

// Summary returns a one-line summary of task results.
func (tr *TaskTracker) Summary() string {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	var passed, failed int
	for _, t := range tr.tasks {
		_, _, state := t.snapshot()
		switch state {
		case taskDone:
			passed++
		case taskFailed:
			failed++
		}
	}

	total := len(tr.tasks)
	if failed > 0 {
		return fmt.Sprintf("%d/%d passed, %d failed", passed, total, failed)
	}
	return fmt.Sprintf("%d/%d passed", passed, total)
}

// String returns the current state of all tasks as plain text (no ANSI).
func (tr *TaskTracker) String() string {
	tr.mu.Lock()
	tasks := tr.tasks
	tr.mu.Unlock()

	nameW := tr.nameWidth()
	var sb strings.Builder
	for _, t := range tasks {
		name, status, state := t.snapshot()
		name = compileGlyphs(name)
		status = compileGlyphs(status)
		icon := Glyph(":pending:")
		switch state {
		case taskDone:
			icon = Glyph(":check:")
		case taskFailed:
			icon = Glyph(":cross:")
		case taskRunning:
			icon = Glyph(":spinner:")
		}
		fmt.Fprintf(&sb, "%s %s %s\n", icon, Pad(name, nameW), status)
	}
	return sb.String()
}

func trackerSpinnerFrame(frame int) string {
	frames := spinnerFramesUnicode
	if currentTheme == ThemeASCII {
		frames = spinnerFramesASCII
	}
	return frames[frame%len(frames)]
}
