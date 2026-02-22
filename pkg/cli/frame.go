package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

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
		out:     os.Stdout,
		done:    make(chan struct{}),
	}
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
	select {
	case <-f.done:
	default:
		close(f.done)
	}
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

	w, h := f.termSize()
	var sb strings.Builder

	order := []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}
	for _, r := range order {
		if _, exists := f.layout.regions[r]; !exists {
			continue
		}
		m, ok := f.models[r]
		if !ok {
			continue
		}
		rw, rh := f.regionSize(r, w, h)
		view := m.View(rw, rh)
		if view != "" {
			sb.WriteString(view)
			if !strings.HasSuffix(view, "\n") {
				sb.WriteByte('\n')
			}
		}
	}

	return sb.String()
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

func (f *Frame) regionSize(r Region, totalW, totalH int) (int, int) {
	// Simple allocation: Header/Footer get 1 line, sidebars get 1/4 width,
	// Content gets the rest.
	switch r {
	case RegionHeader, RegionFooter:
		return totalW, 1
	case RegionLeft, RegionRight:
		return totalW / 4, totalH - 2 // minus header + footer
	case RegionContent:
		sideW := 0
		if _, ok := f.models[RegionLeft]; ok {
			sideW += totalW / 4
		}
		if _, ok := f.models[RegionRight]; ok {
			sideW += totalW / 4
		}
		return totalW - sideW, totalH - 2
	}
	return totalW, totalH
}

func (f *Frame) runLive() {
	// Enter alt-screen.
	fmt.Fprint(f.out, "\033[?1049h")
	// Hide cursor.
	fmt.Fprint(f.out, "\033[?25l")

	defer func() {
		// Show cursor.
		fmt.Fprint(f.out, "\033[?25h")
		// Leave alt-screen.
		fmt.Fprint(f.out, "\033[?1049l")
	}()

	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		f.renderFrame()

		select {
		case <-f.done:
			return
		case <-ticker.C:
		}
	}
}

func (f *Frame) renderFrame() {
	f.mu.Lock()
	defer f.mu.Unlock()

	w, h := f.termSize()

	// Move to top-left.
	fmt.Fprint(f.out, "\033[H")
	// Clear screen.
	fmt.Fprint(f.out, "\033[2J")

	order := []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}
	for _, r := range order {
		if _, exists := f.layout.regions[r]; !exists {
			continue
		}
		m, ok := f.models[r]
		if !ok {
			continue
		}
		rw, rh := f.regionSize(r, w, h)
		view := m.View(rw, rh)
		if view != "" {
			fmt.Fprint(f.out, view)
			if !strings.HasSuffix(view, "\n") {
				fmt.Fprintln(f.out)
			}
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Built-in Region Components
// ─────────────────────────────────────────────────────────────────────────────

// statusLineModel renders a "title  key:value  key:value" bar.
type statusLineModel struct {
	title string
	pairs []string
}

// StatusLine creates a header/footer bar with a title and key:value pairs.
//
//	frame.Header(cli.StatusLine("core dev", "18 repos", "main"))
func StatusLine(title string, pairs ...string) Model {
	return &statusLineModel{title: title, pairs: pairs}
}

func (s *statusLineModel) View(width, _ int) string {
	parts := []string{BoldStyle.Render(s.title)}
	for _, p := range s.pairs {
		parts = append(parts, DimStyle.Render(p))
	}
	line := strings.Join(parts, "  ")
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// keyHintsModel renders keyboard shortcut hints.
type keyHintsModel struct {
	hints []string
}

// KeyHints creates a footer showing keyboard shortcuts.
//
//	frame.Footer(cli.KeyHints("↑/↓ navigate", "enter select", "q quit"))
func KeyHints(hints ...string) Model {
	return &keyHintsModel{hints: hints}
}

func (k *keyHintsModel) View(width, _ int) string {
	parts := make([]string, len(k.hints))
	for i, h := range k.hints {
		parts[i] = DimStyle.Render(h)
	}
	line := strings.Join(parts, "  ")
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// breadcrumbModel renders a navigation path.
type breadcrumbModel struct {
	parts []string
}

// Breadcrumb creates a navigation breadcrumb bar.
//
//	frame.Header(cli.Breadcrumb("core", "dev", "health"))
func Breadcrumb(parts ...string) Model {
	return &breadcrumbModel{parts: parts}
}

func (b *breadcrumbModel) View(width, _ int) string {
	styled := make([]string, len(b.parts))
	for i, p := range b.parts {
		if i == len(b.parts)-1 {
			styled[i] = BoldStyle.Render(p)
		} else {
			styled[i] = DimStyle.Render(p)
		}
	}
	line := strings.Join(styled, DimStyle.Render(" > "))
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// staticModel wraps a plain string as a Model.
type staticModel struct {
	text string
}

// StaticModel wraps a static string as a Model, for use in Frame regions.
func StaticModel(text string) Model {
	return &staticModel{text: text}
}

func (s *staticModel) View(_, _ int) string {
	return s.text
}
