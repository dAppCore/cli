package cli

import (
	"iter"
	"time"

	"dappco.re/go"
	"dappco.re/go/cli/internal/term"
)

var spinnerFramesUnicode = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var spinnerFramesASCII = []string{"-", "\\", "|", "/"}

type taskState int

const (
	taskPending taskState = iota
	taskRunning
	taskDone
	taskFailed
)

type TrackedTask struct {
	name    string
	status  string
	state   taskState
	tracker *TaskTracker
	mu      core.Mutex
}

func (t *TrackedTask) Update(status string) {
	t.mu.Lock()
	t.status = status
	t.state = taskRunning
	t.mu.Unlock()
}

func (t *TrackedTask) Done(message string) {
	t.mu.Lock()
	t.status = message
	t.state = taskDone
	t.mu.Unlock()
}

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

type TaskTracker struct {
	tasks   []*TrackedTask
	out     Writer
	mu      core.Mutex
	started bool
}

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

func NewTaskTracker() *TaskTracker {
	return &TaskTracker{out: stderrWriter()}
}

func (tr *TaskTracker) WithOutput(out Writer) *TaskTracker {
	if out != nil {
		tr.out = out
	}
	return tr
}

func (tr *TaskTracker) Add(name string) *TrackedTask {
	t := &TrackedTask{name: name, status: "waiting", state: taskPending, tracker: tr}
	tr.mu.Lock()
	tr.tasks = append(tr.tasks, t)
	tr.mu.Unlock()
	return t
}

func (tr *TaskTracker) Wait() {
	if !tr.isTTY() {
		tr.waitStatic()
		return
	}
	tr.waitLive()
}

func (tr *TaskTracker) isTTY() bool {
	if fd, ok := writerFileDescriptor(tr.out); ok {
		return term.IsTerminal(fd)
	}
	return false
}

func (tr *TaskTracker) waitStatic() {
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
			core.Print(tr.out, "%s %-20s %s", icon, name, status)
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

	frame := 0
	for i := range n {
		tr.renderLine(i, frame)
	}
	if n == 0 || tr.allDone() {
		return
	}

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		frame++

		tr.mu.Lock()
		count := len(tr.tasks)
		tr.mu.Unlock()

		if r := core.WriteString(tr.out, core.Sprintf("\033[%dA", count)); !r.OK {
			LogWarn("failed to move tracker cursor", "err", r.Error())
		}
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

	core.Print(tr.out, "\033[2K%s %s %s", icon, Pad(name, nameW), styledStatus)
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
		return core.Sprintf("%d/%d passed, %d failed", passed, total, failed)
	}
	return core.Sprintf("%d/%d passed", passed, total)
}

func (tr *TaskTracker) String() string {
	tr.mu.Lock()
	tasks := tr.tasks
	tr.mu.Unlock()

	nameW := tr.nameWidth()
	sb := core.NewBuilder()
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
		core.Print(sb, "%s %s %s", icon, Pad(name, nameW), status)
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
