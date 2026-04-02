package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
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
// Uses HLCRF variant strings for region layout — same as the static Layout system,
// but with live-updating Model components instead of static strings.
//
//	frame := cli.NewFrame("HCF")
//	frame.Header(cli.StatusLine("core dev", "18 repos", "main"))
//	frame.Content(myTableModel)
//	frame.Footer(cli.KeyHints("↑/↓ navigate", "enter select", "q quit"))
//	frame.Run()
type Frame struct {
	variant string
	layout  *Composite
	models  map[Region]Model
	history []Model // content region stack for Navigate/Back
	out     io.Writer
	done    chan struct{}
	mu      sync.Mutex

	// Focus management (bubbletea upgrade)
	focused Region
	keyMap  KeyMap
	width   int
	height  int
	program *tea.Program
}

// NewFrame creates a new Frame with the given HLCRF variant string.
//
//	frame := cli.NewFrame("HCF")      // header, content, footer
//	frame := cli.NewFrame("H[LC]F")   // header, [left + content], footer
func NewFrame(variant string) *Frame {
	return &Frame{
		variant: variant,
		layout:  Layout(variant),
		models:  make(map[Region]Model),
		out:     os.Stderr,
		done:    make(chan struct{}),
		focused: RegionContent,
		keyMap:  DefaultKeyMap(),
		width:   80,
		height:  24,
	}
}

// WithOutput sets the destination writer for rendered output.
// Pass nil to keep the current writer unchanged.
func (f *Frame) WithOutput(out io.Writer) *Frame {
	if out != nil {
		f.out = out
	}
	return f
}

// Header sets the Header region model.
func (f *Frame) Header(m Model) *Frame { f.setModel(RegionHeader, m); return f }

// Left sets the Left sidebar region model.
func (f *Frame) Left(m Model) *Frame { f.setModel(RegionLeft, m); return f }

// Content sets the Content region model.
func (f *Frame) Content(m Model) *Frame { f.setModel(RegionContent, m); return f }

// Right sets the Right sidebar region model.
func (f *Frame) Right(m Model) *Frame { f.setModel(RegionRight, m); return f }

// Footer sets the Footer region model.
func (f *Frame) Footer(m Model) *Frame { f.setModel(RegionFooter, m); return f }

func (f *Frame) setModel(r Region, m Model) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.models[r] = m
}

// Navigate replaces the Content region with a new model, pushing the current one
// onto the history stack for Back().
func (f *Frame) Navigate(m Model) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if current, ok := f.models[RegionContent]; ok {
		f.history = append(f.history, current)
	}
	f.models[RegionContent] = m
}

// Back pops the content history stack, restoring the previous Content model.
// Returns false if the history is empty.
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

// Stop signals the Frame to exit its Run loop.
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

// Send injects a message into the Frame's tea.Program.
// Safe to call before Run() (message is discarded).
func (f *Frame) Send(msg tea.Msg) {
	if f.program != nil {
		f.program.Send(msg)
	}
}

// WithKeyMap sets custom key bindings for Frame navigation.
func (f *Frame) WithKeyMap(km KeyMap) *Frame {
	f.keyMap = km
	return f
}

// Focused returns the currently focused region.
func (f *Frame) Focused() Region {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.focused
}

// Focus sets focus to a specific region.
// Ignores the request if the region is not in this Frame's variant.
func (f *Frame) Focus(r Region) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.layout.regions[r]; exists {
		f.focused = r
	}
}

// buildFocusRing returns the ordered list of regions in this Frame's variant.
// Order follows HLCRF convention.
func (f *Frame) buildFocusRing() []Region {
	order := []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}
	var ring []Region
	for _, r := range order {
		if _, exists := f.layout.regions[r]; exists {
			ring = append(ring, r)
		}
	}
	return ring
}

// Init implements tea.Model. Collects Init() from all FrameModel regions.
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

// Update implements tea.Model. Routes messages based on type and focus.
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
			// Forward to focused region
			return f, f.updateFocusedLocked(msg)
		}

	default:
		// Broadcast non-key messages to all regions
		return f, f.broadcastLocked(msg)
	}
}

// View implements tea.Model. Composes region views using lipgloss.
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

	// Calculate region dimensions
	headerH, footerH := 0, 0
	if _, ok := f.layout.regions[RegionHeader]; ok {
		if _, ok := f.models[RegionHeader]; ok {
			headerH = 1
		}
	}
	if _, ok := f.layout.regions[RegionFooter]; ok {
		if _, ok := f.models[RegionFooter]; ok {
			footerH = 1
		}
	}
	middleH := max(h-headerH-footerH, 1)

	// Render each region
	header := f.renderRegionLocked(RegionHeader, w, headerH)
	footer := f.renderRegionLocked(RegionFooter, w, footerH)

	// Calculate sidebar widths
	leftW, rightW := 0, 0
	if _, ok := f.layout.regions[RegionLeft]; ok {
		if _, ok := f.models[RegionLeft]; ok {
			leftW = w / 4
		}
	}
	if _, ok := f.layout.regions[RegionRight]; ok {
		if _, ok := f.models[RegionRight]; ok {
			rightW = w / 4
		}
	}
	contentW := max(w-leftW-rightW, 1)

	left := f.renderRegionLocked(RegionLeft, leftW, middleH)
	right := f.renderRegionLocked(RegionRight, rightW, middleH)
	content := f.renderRegionLocked(RegionContent, contentW, middleH)

	// Compose middle row
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

	// Compose full layout
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

// cycleFocusLocked moves focus forward (+1) or backward (-1) in the focus ring.
// Must be called with f.mu held.
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

// spatialFocusLocked moves focus to a specific region if it exists in the layout.
// Must be called with f.mu held.
func (f *Frame) spatialFocusLocked(target Region) {
	if _, exists := f.layout.regions[target]; exists {
		f.focused = target
	}
}

// backLocked pops the content history. Must be called with f.mu held.
func (f *Frame) backLocked() {
	if len(f.history) == 0 {
		return
	}
	f.models[RegionContent] = f.history[len(f.history)-1]
	f.history = f.history[:len(f.history)-1]
}

// broadcastLocked sends a message to all FrameModel regions.
// Must be called with f.mu held.
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

// updateFocusedLocked sends a message to only the focused region.
// Must be called with f.mu held.
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

// Run renders the frame and blocks. In TTY mode, it live-refreshes at ~12fps.
// In non-TTY mode, it renders once and returns immediately.
func (f *Frame) Run() {
	if !f.isTTY() {
		fmt.Fprint(f.out, f.String())
		return
	}
	f.runLive()
}

// RunFor runs the frame for a fixed duration, then stops.
// Useful for dashboards that refresh periodically.
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
// This is the non-TTY fallback path.
func (f *Frame) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	view := f.viewLocked()
	if view == "" {
		return ""
	}
	view = ansi.Strip(view)
	// Ensure trailing newline for non-TTY consistency
	if !strings.HasSuffix(view, "\n") {
		view += "\n"
	}
	return view
}

func (f *Frame) isTTY() bool {
	if file, ok := f.out.(*os.File); ok {
		return term.IsTerminal(int(file.Fd()))
	}
	return false
}

func (f *Frame) termSize() (int, int) {
	if file, ok := f.out.(*os.File); ok {
		w, h, err := term.GetSize(int(file.Fd()))
		if err == nil {
			return w, h
		}
	}
	return 80, 24 // sensible default
}

func (f *Frame) runLive() {
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}
	if f.out != os.Stdout {
		opts = append(opts, tea.WithOutput(f.out))
	}

	p := tea.NewProgram(f, opts...)
	f.program = p

	if _, err := p.Run(); err != nil {
		Error(err.Error())
	}
}
