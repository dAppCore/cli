package frame

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Model is the interface for components that slot into Frame regions.
// View receives the allocated width and height and returns rendered text.
type Model interface {
	View(width, height int) string
}

// ModelFunc is a convenience adapter for using a function as a Model.
type ModelFunc func(width, height int) string

// View implements Model.
func (f ModelFunc) View(width, height int) string { return f(width, height) }

// Frame is a live compositional AppShell for TUI.
type Frame struct {
	variant string
	layout  *Composite
	models  map[Region]Model
	history []Model
	out     Writer
	done    chan struct{}
	mu      sync.Mutex

	focused Region
	keyMap  KeyMap
	width   int
	height  int
	program *tea.Program
}

// NewFrame creates a new Frame with the given HLCRF variant string.
func NewFrame(variant string) *Frame {
	return &Frame{
		variant: variant,
		layout:  Layout(variant),
		models:  make(map[Region]Model),
		out:     stderrWriter(),
		done:    make(chan struct{}),
		focused: RegionContent,
		keyMap:  DefaultKeyMap(),
		width:   80,
		height:  24,
	}
}

func (f *Frame) WithOutput(out Writer) *Frame {
	if out != nil {
		f.out = out
	}
	return f
}

func (f *Frame) Header(m Model) *Frame  { f.setModel(RegionHeader, m); return f }
func (f *Frame) Left(m Model) *Frame    { f.setModel(RegionLeft, m); return f }
func (f *Frame) Content(m Model) *Frame { f.setModel(RegionContent, m); return f }
func (f *Frame) Right(m Model) *Frame   { f.setModel(RegionRight, m); return f }
func (f *Frame) Footer(m Model) *Frame  { f.setModel(RegionFooter, m); return f }

func (f *Frame) setModel(r Region, m Model) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.models[r] = m
}

func (f *Frame) Navigate(m Model) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if current, ok := f.models[RegionContent]; ok {
		f.history = append(f.history, current)
	}
	f.models[RegionContent] = m
}

func (f *Frame) Back() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.history) == 0 {
		return false
	}
	f.models[RegionContent] = f.history[len(f.history)-1]
	f.history = f.history[:len(f.history)-1]
	return true
}

func (f *Frame) Stop() {
	if f.program != nil {
		f.program.Quit()
		return
	}
	select {
	case <-f.done:
	default:
		close(f.done)
	}
}

func (f *Frame) Send(msg tea.Msg) {
	if f.program != nil {
		f.program.Send(msg)
	}
}

func (f *Frame) WithKeyMap(km KeyMap) *Frame {
	f.keyMap = km
	return f
}

func (f *Frame) Focused() Region {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.focused
}

func (f *Frame) Focus(r Region) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if layoutHasRegion(f.layout, r) {
		f.focused = r
	}
}

func (f *Frame) buildFocusRing() []Region {
	order := []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}
	var ring []Region
	for _, r := range order {
		if layoutHasRegion(f.layout, r) {
			ring = append(ring, r)
		}
	}
	return ring
}

func (f *Frame) Init() tea.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()
	var cmds []tea.Cmd
	for _, m := range f.models {
		fm := adaptModel(m)
		if cmd := fm.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (f *Frame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	f.mu.Lock()
	defer f.mu.Unlock()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		return f, f.broadcastLocked(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case f.keyMap.Quit:
			return f, tea.Quit
		case f.keyMap.Back:
			f.backLocked()
			return f, nil
		case f.keyMap.FocusNext:
			f.cycleFocusLocked(1)
			return f, nil
		case f.keyMap.FocusPrev:
			f.cycleFocusLocked(-1)
			return f, nil
		case f.keyMap.FocusUp:
			f.spatialFocusLocked(RegionHeader)
			return f, nil
		case f.keyMap.FocusDown:
			f.spatialFocusLocked(RegionFooter)
			return f, nil
		case f.keyMap.FocusLeft:
			f.spatialFocusLocked(RegionLeft)
			return f, nil
		case f.keyMap.FocusRight:
			f.spatialFocusLocked(RegionRight)
			return f, nil
		default:
			return f, f.updateFocusedLocked(msg)
		}
	default:
		return f, f.broadcastLocked(msg)
	}
}

func (f *Frame) View() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.viewLocked()
}

func (f *Frame) viewLocked() string {
	w, h := f.width, f.height
	if w == 0 || h == 0 {
		w, h = f.termSize()
	}

	headerH, footerH := 0, 0
	if layoutHasRegion(f.layout, RegionHeader) {
		if _, ok := f.models[RegionHeader]; ok {
			headerH = 1
		}
	}
	if layoutHasRegion(f.layout, RegionFooter) {
		if _, ok := f.models[RegionFooter]; ok {
			footerH = 1
		}
	}
	middleH := max(h-headerH-footerH, 1)

	header := f.renderRegionLocked(RegionHeader, w, headerH)
	footer := f.renderRegionLocked(RegionFooter, w, footerH)

	leftW, rightW := 0, 0
	if layoutHasRegion(f.layout, RegionLeft) {
		if _, ok := f.models[RegionLeft]; ok {
			leftW = w / 4
		}
	}
	if layoutHasRegion(f.layout, RegionRight) {
		if _, ok := f.models[RegionRight]; ok {
			rightW = w / 4
		}
	}
	contentW := max(w-leftW-rightW, 1)

	left := f.renderRegionLocked(RegionLeft, leftW, middleH)
	right := f.renderRegionLocked(RegionRight, rightW, middleH)
	content := f.renderRegionLocked(RegionContent, contentW, middleH)

	var middleParts []string
	if leftW > 0 {
		middleParts = append(middleParts, left)
	}
	middleParts = append(middleParts, content)
	if rightW > 0 {
		middleParts = append(middleParts, right)
	}

	middle := content
	if len(middleParts) > 1 {
		middle = lipgloss.JoinHorizontal(lipgloss.Top, middleParts...)
	}

	var verticalParts []string
	if headerH > 0 {
		verticalParts = append(verticalParts, header)
	}
	verticalParts = append(verticalParts, middle)
	if footerH > 0 {
		verticalParts = append(verticalParts, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, verticalParts...)
}

func (f *Frame) renderRegionLocked(r Region, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	m, ok := f.models[r]
	if !ok {
		return ""
	}
	fm := adaptModel(m)
	return fm.View(w, h)
}

func (f *Frame) cycleFocusLocked(dir int) {
	ring := f.buildFocusRing()
	if len(ring) == 0 {
		return
	}
	idx := 0
	for i, r := range ring {
		if r == f.focused {
			idx = i
			break
		}
	}
	idx = (idx + dir + len(ring)) % len(ring)
	f.focused = ring[idx]
}

func (f *Frame) spatialFocusLocked(target Region) {
	if layoutHasRegion(f.layout, target) {
		f.focused = target
	}
}

func (f *Frame) backLocked() {
	if len(f.history) == 0 {
		return
	}
	f.models[RegionContent] = f.history[len(f.history)-1]
	f.history = f.history[:len(f.history)-1]
}

func (f *Frame) broadcastLocked(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for r, m := range f.models {
		fm := adaptModel(m)
		updated, cmd := fm.Update(msg)
		f.models[r] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (f *Frame) updateFocusedLocked(msg tea.Msg) tea.Cmd {
	m, ok := f.models[f.focused]
	if !ok {
		return nil
	}
	fm := adaptModel(m)
	updated, cmd := fm.Update(msg)
	f.models[f.focused] = updated
	return cmd
}

// Run renders the frame and blocks.
func (f *Frame) Run() {
	if !f.isTTY() {
		_, _ = f.out.Write([]byte(f.String()))
		return
	}
	f.runLive()
}

// RunFor runs the frame for a fixed duration, then stops.
func (f *Frame) RunFor(d time.Duration) {
	go func() {
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-timer.C:
			f.Stop()
		case <-f.done:
		}
	}()
	f.Run()
}

// String renders the frame as a static string (no ANSI, no live updates).
func (f *Frame) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	view := f.viewLocked()
	if view == "" {
		return ""
	}
	view = ansi.Strip(view)
	if !strings.HasSuffix(view, "\n") {
		view += "\n"
	}
	return view
}

func (f *Frame) isTTY() bool {
	if fd, ok := writerFileDescriptor(f.out); ok {
		return isTerminal(fd)
	}
	return false
}

func (f *Frame) termSize() (int, int) {
	if fd, ok := writerFileDescriptor(f.out); ok {
		w, h, err := terminalSize(fd)
		if err == nil {
			return w, h
		}
	}
	return 80, 24
}

func (f *Frame) runLive() {
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}
	if f.out != stdoutWriter() {
		opts = append(opts, tea.WithOutput(f.out))
	}

	p := tea.NewProgram(f, opts...)
	f.program = p

	if _, err := p.Run(); err != nil {
		_, _ = fmt.Fprintln(stderrWriter(), err.Error())
	}
}
