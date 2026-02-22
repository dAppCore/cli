# Frame Bubbletea Upgrade Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Upgrade `cli.Frame` from raw ANSI + `golang.org/x/term` to bubbletea internally, adding keyboard focus management and lipgloss layout, while preserving the existing public API.

**Architecture:** Frame implements `tea.Model` internally and owns a single `tea.Program`. A dual interface pattern keeps the existing `Model` (view-only) working alongside a new `FrameModel` (interactive). Lipgloss replaces manual ANSI escape codes for layout composition.

**Tech Stack:** Go 1.26, bubbletea v1.3.10, lipgloss v1.1.0, existing HLCRF layout parser

---

## Important Context

**Repo:** `~/Code/core/cli` (module `forge.lthn.ai/core/cli`)

**Workspace:** `~/Code/go.work` — Go workspace with 29 modules. Run all commands from `~/Code/core/cli/`.

**Run tests:** `go test -race ./pkg/cli/` (always from `~/Code/core/cli/`)

**Design doc:** `docs/plans/2026-02-22-frame-bubbletea-design.md`

**Key files you'll touch:**
- `pkg/cli/frame.go` — current Frame (359 lines, raw ANSI rendering)
- `pkg/cli/frame_model.go` — **new** file for FrameModel interface, KeyMap, adapter
- `pkg/cli/frame_test.go` — existing 14 tests (must all keep passing)
- `go.mod` — add bubbletea + lipgloss deps

**Key files to read (don't modify):**
- `pkg/cli/layout.go` — HLCRF variant parser (`Region` type, `Composite` struct)
- `pkg/cli/ansi.go` — `AnsiStyle`, `SetColorEnabled()`, `ColorEnabled()`
- `pkg/cli/styles.go` — `BoldStyle`, `DimStyle`, `Truncate()`, `Pad()`

**bubbletea API (v1.3.10) cheatsheet:**
- `tea.Model` interface: `Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)`, `View() string`
- `tea.NewProgram(model, opts...)` — creates program
- `tea.WithAltScreen()` — fullscreen mode
- `tea.WithOutput(io.Writer)` — custom output
- `tea.Batch(cmds...)` — combine commands
- `tea.Quit()` — exit command
- `tea.KeyMsg` — has `.Type` (KeyType) and `.String()` method
- Key constants: `tea.KeyTab`, `tea.KeyShiftTab`, `tea.KeyUp`, `tea.KeyDown`, `tea.KeyLeft`, `tea.KeyRight`, `tea.KeyEsc`, `tea.KeyCtrlC`
- `tea.WindowSizeMsg` — has `.Width`, `.Height`
- `program.Send(msg)` — inject message from outside
- `program.Quit()` — stop program

**lipgloss API (v1.1.0) cheatsheet:**
- `lipgloss.JoinVertical(pos, strs...)` — stack strings vertically
- `lipgloss.JoinHorizontal(pos, strs...)` — join strings side-by-side
- `lipgloss.Place(w, h, hPos, vPos, str)` — place string in box
- Constants: `lipgloss.Left`, `lipgloss.Right`, `lipgloss.Center`, `lipgloss.Top`, `lipgloss.Bottom`
- `lipgloss.NewStyle().Width(n).Height(n).Render(str)` — constrain to dimensions

---

### Task 1: Add bubbletea and lipgloss dependencies

**Files:**
- Modify: `go.mod`

**Step 1: Add the dependencies**

Run from `~/Code/core/cli/`:

```bash
go get github.com/charmbracelet/bubbletea@v1.3.10
go get github.com/charmbracelet/lipgloss@v1.1.0
```

**Step 2: Tidy**

```bash
go mod tidy
```

**Step 3: Verify existing tests still pass**

Run: `go test -race ./pkg/cli/`
Expected: PASS (all 14 existing tests unchanged)

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: add bubbletea and lipgloss for Frame upgrade"
```

---

### Task 2: Create FrameModel interface and modelAdapter

**Files:**
- Create: `pkg/cli/frame_model.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the failing test**

Add to `pkg/cli/frame_test.go` at the bottom:

```go
func TestFrameModel_Good(t *testing.T) {
	t.Run("modelAdapter wraps plain Model", func(t *testing.T) {
		m := StaticModel("hello")
		adapted := adaptModel(m)

		// Should return nil cmd from Init
		cmd := adapted.Init()
		assert.Nil(t, cmd)

		// Should return itself from Update
		updated, cmd := adapted.Update(nil)
		assert.Equal(t, adapted, updated)
		assert.Nil(t, cmd)

		// Should delegate View to wrapped model
		assert.Equal(t, "hello", adapted.View(80, 24))
	})

	t.Run("FrameModel passes through without wrapping", func(t *testing.T) {
		fm := &testFrameModel{viewText: "interactive"}
		adapted := adaptModel(fm)

		// Should be the same object, not wrapped
		_, ok := adapted.(*testFrameModel)
		assert.True(t, ok, "FrameModel should not be wrapped")
		assert.Equal(t, "interactive", adapted.View(80, 24))
	})
}

// testFrameModel is a mock FrameModel for testing.
type testFrameModel struct {
	viewText     string
	initCalled   bool
	updateCalled bool
	lastMsg      tea.Msg
}

func (m *testFrameModel) View(w, h int) string { return m.viewText }

func (m *testFrameModel) Init() tea.Cmd {
	m.initCalled = true
	return nil
}

func (m *testFrameModel) Update(msg tea.Msg) (FrameModel, tea.Cmd) {
	m.updateCalled = true
	m.lastMsg = msg
	return m, nil
}
```

You'll need to add `tea "github.com/charmbracelet/bubbletea"` to the test file's imports.

**Step 2: Run test to verify it fails**

Run: `go test -race -run TestFrameModel ./pkg/cli/`
Expected: FAIL — `adaptModel` undefined, `FrameModel` undefined, `testFrameModel` can't satisfy unwritten interface

**Step 3: Write the implementation**

Create `pkg/cli/frame_model.go`:

```go
package cli

import tea "github.com/charmbracelet/bubbletea"

// FrameModel extends Model with bubbletea lifecycle methods.
// Use this for interactive components that handle input.
// Plain Model components work unchanged — Frame wraps them automatically.
type FrameModel interface {
	Model
	Init() tea.Cmd
	Update(tea.Msg) (FrameModel, tea.Cmd)
}

// adaptModel wraps a plain Model as a FrameModel via modelAdapter.
// If the model already implements FrameModel, it is returned as-is.
func adaptModel(m Model) FrameModel {
	if fm, ok := m.(FrameModel); ok {
		return fm
	}
	return &modelAdapter{m: m}
}

// modelAdapter wraps a plain Model to satisfy FrameModel.
// Init returns nil, Update is a no-op, View delegates to the wrapped Model.
type modelAdapter struct {
	m Model
}

func (a *modelAdapter) View(w, h int) string               { return a.m.View(w, h) }
func (a *modelAdapter) Init() tea.Cmd                      { return nil }
func (a *modelAdapter) Update(tea.Msg) (FrameModel, tea.Cmd) { return a, nil }
```

**Step 4: Run test to verify it passes**

Run: `go test -race -run TestFrameModel ./pkg/cli/`
Expected: PASS

**Step 5: Run all tests to verify no regressions**

Run: `go test -race ./pkg/cli/`
Expected: PASS (all existing + new tests)

**Step 6: Commit**

```bash
git add pkg/cli/frame_model.go pkg/cli/frame_test.go
git commit -m "feat(frame): add FrameModel interface and modelAdapter"
```

---

### Task 3: Add KeyMap struct with defaults

**Files:**
- Modify: `pkg/cli/frame_model.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the failing test**

Add to `pkg/cli/frame_test.go`:

```go
func TestKeyMap_Good(t *testing.T) {
	t.Run("default keymap has expected bindings", func(t *testing.T) {
		km := DefaultKeyMap()
		assert.Equal(t, tea.KeyTab, km.FocusNext)
		assert.Equal(t, tea.KeyShiftTab, km.FocusPrev)
		assert.Equal(t, tea.KeyUp, km.FocusUp)
		assert.Equal(t, tea.KeyDown, km.FocusDown)
		assert.Equal(t, tea.KeyLeft, km.FocusLeft)
		assert.Equal(t, tea.KeyRight, km.FocusRight)
		assert.Equal(t, tea.KeyEsc, km.Back)
		assert.Equal(t, tea.KeyCtrlC, km.Quit)
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -race -run TestKeyMap ./pkg/cli/`
Expected: FAIL — `DefaultKeyMap` undefined, `KeyMap` undefined

**Step 3: Write the implementation**

Add to `pkg/cli/frame_model.go`:

```go
// KeyMap defines key bindings for Frame navigation.
// Use DefaultKeyMap() for sensible defaults, or build your own.
type KeyMap struct {
	FocusNext  tea.KeyType // Tab — cycle focus forward
	FocusPrev  tea.KeyType // Shift-Tab — cycle focus backward
	FocusUp    tea.KeyType // Up — spatial: move to Header
	FocusDown  tea.KeyType // Down — spatial: move to Footer
	FocusLeft  tea.KeyType // Left — spatial: move to Left sidebar
	FocusRight tea.KeyType // Right — spatial: move to Right sidebar
	Back       tea.KeyType // Esc — Navigate back
	Quit       tea.KeyType // Ctrl-C — quit
}

// DefaultKeyMap returns the standard Frame key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		FocusNext:  tea.KeyTab,
		FocusPrev:  tea.KeyShiftTab,
		FocusUp:    tea.KeyUp,
		FocusDown:  tea.KeyDown,
		FocusLeft:  tea.KeyLeft,
		FocusRight: tea.KeyRight,
		Back:       tea.KeyEsc,
		Quit:       tea.KeyCtrlC,
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test -race -run TestKeyMap ./pkg/cli/`
Expected: PASS

**Step 5: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS

**Step 6: Commit**

```bash
git add pkg/cli/frame_model.go pkg/cli/frame_test.go
git commit -m "feat(frame): add KeyMap with default bindings"
```

---

### Task 4: Add focus management fields to Frame

**Files:**
- Modify: `pkg/cli/frame.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the failing tests**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameFocus_Good(t *testing.T) {
	t.Run("default focus is Content", func(t *testing.T) {
		f := NewFrame("HCF")
		assert.Equal(t, RegionContent, f.Focused())
	})

	t.Run("Focus sets focused region", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Focus(RegionHeader)
		assert.Equal(t, RegionHeader, f.Focused())
	})

	t.Run("Focus ignores invalid region", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Focus(RegionLeft) // Left not in "HCF"
		assert.Equal(t, RegionContent, f.Focused()) // unchanged
	})

	t.Run("WithKeyMap returns frame for chaining", func(t *testing.T) {
		km := DefaultKeyMap()
		km.Quit = tea.KeyCtrlQ
		f := NewFrame("HCF").WithKeyMap(km)
		assert.Equal(t, tea.KeyCtrlQ, f.keyMap.Quit)
	})

	t.Run("focusRing builds from variant", func(t *testing.T) {
		f := NewFrame("HLCRF")
		ring := f.buildFocusRing()
		assert.Equal(t, []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}, ring)
	})

	t.Run("focusRing respects variant order", func(t *testing.T) {
		f := NewFrame("HCF")
		ring := f.buildFocusRing()
		assert.Equal(t, []Region{RegionHeader, RegionContent, RegionFooter}, ring)
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -race -run TestFrameFocus ./pkg/cli/`
Expected: FAIL — `Focused()`, `Focus()`, `WithKeyMap()`, `buildFocusRing()` all undefined

**Step 3: Write the implementation**

Modify `pkg/cli/frame.go`. Add new fields to the `Frame` struct:

```go
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
```

Add `tea "github.com/charmbracelet/bubbletea"` to frame.go imports (alongside existing ones). You do NOT need to import lipgloss yet.

Update `NewFrame` to initialise new fields:

```go
func NewFrame(variant string) *Frame {
	return &Frame{
		variant: variant,
		layout:  Layout(variant),
		models:  make(map[Region]Model),
		out:     os.Stdout,
		done:    make(chan struct{}),
		focused: RegionContent,
		keyMap:  DefaultKeyMap(),
		width:   80,
		height:  24,
	}
}
```

Add the new public methods:

```go
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
```

**Step 4: Run test to verify it passes**

Run: `go test -race -run TestFrameFocus ./pkg/cli/`
Expected: PASS

**Step 5: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS

**Step 6: Commit**

```bash
git add pkg/cli/frame.go pkg/cli/frame_test.go
git commit -m "feat(frame): add focus management fields, Focused(), Focus(), WithKeyMap()"
```

---

### Task 5: Implement tea.Model on Frame (Init, Update, View)

This is the core task. Frame becomes a `tea.Model`. The existing `runLive()` and `renderFrame()` methods get replaced.

**Files:**
- Modify: `pkg/cli/frame.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the failing tests**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameTeaModel_Good(t *testing.T) {
	t.Run("Init collects FrameModel inits", func(t *testing.T) {
		f := NewFrame("HCF")
		fm := &testFrameModel{viewText: "x"}
		f.Content(fm)

		cmd := f.Init()
		// Should produce a batch command (non-nil if any FrameModel has Init)
		// fm.Init returns nil, so batch of nils = nil
		_ = cmd // no panic = success
		assert.True(t, fm.initCalled)
	})

	t.Run("Update routes key to focused region", func(t *testing.T) {
		f := NewFrame("HCF")
		header := &testFrameModel{viewText: "h"}
		content := &testFrameModel{viewText: "c"}
		f.Header(header)
		f.Content(content)

		// Focus is Content by default
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		f.Update(keyMsg)

		assert.True(t, content.updateCalled, "focused region should receive key")
		assert.False(t, header.updateCalled, "unfocused region should not receive key")
	})

	t.Run("Update broadcasts WindowSizeMsg to all", func(t *testing.T) {
		f := NewFrame("HCF")
		header := &testFrameModel{viewText: "h"}
		content := &testFrameModel{viewText: "c"}
		footer := &testFrameModel{viewText: "f"}
		f.Header(header)
		f.Content(content)
		f.Footer(footer)

		sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
		f.Update(sizeMsg)

		assert.True(t, header.updateCalled, "header should get resize")
		assert.True(t, content.updateCalled, "content should get resize")
		assert.True(t, footer.updateCalled, "footer should get resize")
		assert.Equal(t, 120, f.width)
		assert.Equal(t, 40, f.height)
	})

	t.Run("Update handles quit key", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Content(StaticModel("c"))

		quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
		_, cmd := f.Update(quitMsg)

		// cmd should be tea.Quit
		assert.NotNil(t, cmd)
	})

	t.Run("Update handles back key", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Content(StaticModel("page-1"))
		f.Navigate(StaticModel("page-2"))

		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		f.Update(escMsg)

		assert.Contains(t, f.String(), "page-1")
	})

	t.Run("Update cycles focus with Tab", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		assert.Equal(t, RegionContent, f.Focused())

		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		f.Update(tabMsg)
		assert.Equal(t, RegionFooter, f.Focused())

		f.Update(tabMsg)
		assert.Equal(t, RegionHeader, f.Focused()) // wraps around

		shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
		f.Update(shiftTabMsg)
		assert.Equal(t, RegionFooter, f.Focused()) // back
	})

	t.Run("View produces non-empty output", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.Header(StaticModel("HEAD"))
		f.Content(StaticModel("BODY"))
		f.Footer(StaticModel("FOOT"))

		view := f.View()
		assert.Contains(t, view, "HEAD")
		assert.Contains(t, view, "BODY")
		assert.Contains(t, view, "FOOT")
	})

	t.Run("View lipgloss layout: header before content before footer", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.Header(StaticModel("AAA"))
		f.Content(StaticModel("BBB"))
		f.Footer(StaticModel("CCC"))
		f.width = 80
		f.height = 24

		view := f.View()
		posA := indexOf(view, "AAA")
		posB := indexOf(view, "BBB")
		posC := indexOf(view, "CCC")
		assert.Greater(t, posA, -1, "header should be present")
		assert.Greater(t, posB, -1, "content should be present")
		assert.Greater(t, posC, -1, "footer should be present")
		assert.Less(t, posA, posB, "header before content")
		assert.Less(t, posB, posC, "content before footer")
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -race -run TestFrameTeaModel ./pkg/cli/`
Expected: FAIL — `Init()`, `Update(tea.Msg)`, `View()` don't exist on Frame (wrong signatures from what tea.Model needs)

**Step 3: Write the implementation**

This is the biggest change. Modify `pkg/cli/frame.go`:

**Add lipgloss import:**

```go
import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)
```

**Add the three tea.Model methods:**

```go
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
	middleH := h - headerH - footerH
	if middleH < 1 {
		middleH = 1
	}

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
	contentW := w - leftW - rightW
	if contentW < 1 {
		contentW = 1
	}

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
```

**Add internal focus helpers (inside frame.go):**

```go
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
```

**Step 4: Run test to verify it passes**

Run: `go test -race -run TestFrameTeaModel ./pkg/cli/`
Expected: PASS

**Step 5: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS — existing tests use `String()` (non-TTY path) which is unchanged

**Step 6: Commit**

```bash
git add pkg/cli/frame.go pkg/cli/frame_test.go
git commit -m "feat(frame): implement tea.Model (Init, Update, View) with lipgloss layout"
```

---

### Task 6: Replace runLive() with tea.Program

**Files:**
- Modify: `pkg/cli/frame.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the failing test**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameSend_Good(t *testing.T) {
	t.Run("Send is safe before Run", func(t *testing.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))

		// Should not panic when program is nil
		assert.NotPanics(t, func() {
			f.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
		})
	})
}
```

**Step 2: Run test to verify it fails**

Run: `go test -race -run TestFrameSend ./pkg/cli/`
Expected: FAIL — `Send` method doesn't exist

**Step 3: Write the implementation**

Modify `pkg/cli/frame.go`:

**Replace `runLive()`:**

```go
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
```

**Delete the old `renderFrame()` method** — it's no longer used (View() replaces it).

**Update `Stop()`:**

```go
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
```

**Add `Send()`:**

```go
// Send injects a message into the Frame's tea.Program.
// Safe to call before Run() (message is discarded).
func (f *Frame) Send(msg tea.Msg) {
	if f.program != nil {
		f.program.Send(msg)
	}
}
```

**Update `RunFor()`** to work with tea.Program:

```go
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
```

**Step 4: Run test to verify it passes**

Run: `go test -race -run TestFrameSend ./pkg/cli/`
Expected: PASS

**Step 5: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS

**Step 6: Commit**

```bash
git add pkg/cli/frame.go pkg/cli/frame_test.go
git commit -m "feat(frame): replace raw ANSI runLive with tea.Program"
```

---

### Task 7: Clean up dead code

**Files:**
- Modify: `pkg/cli/frame.go`

**Step 1: Review and remove unused code**

The old `renderFrame()` method should have been removed in Task 6. Verify it's gone.

Also check if `f.done` channel is still needed. It's used by `RunFor()` as a fallback and by the non-TTY `Stop()` path. Keep it for now.

Check if `golang.org/x/term` import is still needed. `isTTY()` and `termSize()` still use it for the non-TTY fallback path. Keep it.

**Step 2: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS

**Step 3: Run go vet**

Run: `go vet ./pkg/cli/`
Expected: no warnings

**Step 4: Commit (only if there were changes)**

```bash
git add pkg/cli/frame.go
git commit -m "refactor(frame): remove unused renderFrame method"
```

---

### Task 8: Update String() to use viewLocked()

The `String()` method currently has its own rendering logic separate from `View()`. Unify them so `String()` delegates to the same lipgloss-based layout.

**Files:**
- Modify: `pkg/cli/frame.go`
- Test: `pkg/cli/frame_test.go`

**Step 1: Verify existing String() tests still define the contract**

The existing tests assert:
- `"static render HCF"` — contains header, content, footer
- `"region order preserved"` — header before content before footer
- `"empty regions skipped"` — only content → "only content\n"
- `"empty frame"` — no models → ""

These must still pass after the change.

**Step 2: Update String()**

Replace the existing `String()` method:

```go
func (f *Frame) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	view := f.viewLocked()
	if view == "" {
		return ""
	}
	// Ensure trailing newline for non-TTY consistency
	if !strings.HasSuffix(view, "\n") {
		view += "\n"
	}
	return view
}
```

**Step 3: Run all tests**

Run: `go test -race ./pkg/cli/`
Expected: PASS — the existing test assertions should still hold because `viewLocked()` produces the same output structure (header, content, footer in order)

**Important:** If `"empty regions skipped"` fails (expected `"only content\n"` but lipgloss adds padding), you may need to adjust `viewLocked()` to not use `lipgloss.Place()` for content — just return the raw view string when there's only one region. The fix:

In `viewLocked()`, if only one vertical part exists (no header, no footer, just content), return it directly without lipgloss wrapping.

**Step 4: Commit**

```bash
git add pkg/cli/frame.go
git commit -m "refactor(frame): unify String() with View() via viewLocked()"
```

---

### Task 9: Add spatial focus navigation tests

**Files:**
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the tests**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameSpatialFocus_Good(t *testing.T) {
	t.Run("arrow keys move to target region", func(t *testing.T) {
		f := NewFrame("HLCRF")
		f.Header(StaticModel("h"))
		f.Left(StaticModel("l"))
		f.Content(StaticModel("c"))
		f.Right(StaticModel("r"))
		f.Footer(StaticModel("f"))

		// Start at Content
		assert.Equal(t, RegionContent, f.Focused())

		// Up → Header
		f.Update(tea.KeyMsg{Type: tea.KeyUp})
		assert.Equal(t, RegionHeader, f.Focused())

		// Down → Footer
		f.Update(tea.KeyMsg{Type: tea.KeyDown})
		assert.Equal(t, RegionFooter, f.Focused())

		// Left → Left sidebar
		f.Update(tea.KeyMsg{Type: tea.KeyLeft})
		assert.Equal(t, RegionLeft, f.Focused())

		// Right → Right sidebar
		f.Update(tea.KeyMsg{Type: tea.KeyRight})
		assert.Equal(t, RegionRight, f.Focused())
	})

	t.Run("spatial focus ignores missing regions", func(t *testing.T) {
		f := NewFrame("HCF") // no Left or Right
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		assert.Equal(t, RegionContent, f.Focused())

		// Left arrow → no Left region, focus stays
		f.Update(tea.KeyMsg{Type: tea.KeyLeft})
		assert.Equal(t, RegionContent, f.Focused())
	})
}
```

**Step 2: Run tests**

Run: `go test -race -run TestFrameSpatialFocus ./pkg/cli/`
Expected: PASS (these should work with the Update logic from Task 5)

**Step 3: Commit**

```bash
git add pkg/cli/frame_test.go
git commit -m "test(frame): add spatial focus navigation tests"
```

---

### Task 10: Add Navigate/Back FrameModel focus transfer tests

**Files:**
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the tests**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameNavigateFrameModel_Good(t *testing.T) {
	t.Run("Navigate with FrameModel preserves focus on Content", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(&testFrameModel{viewText: "page-1"})
		f.Footer(StaticModel("f"))

		// Focus something else
		f.Focus(RegionHeader)
		assert.Equal(t, RegionHeader, f.Focused())

		// Navigate replaces Content, focus should remain where it was
		f.Navigate(&testFrameModel{viewText: "page-2"})
		assert.Equal(t, RegionHeader, f.Focused())
		assert.Contains(t, f.String(), "page-2")
	})

	t.Run("Back restores FrameModel", func(t *testing.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		fm1 := &testFrameModel{viewText: "page-1"}
		fm2 := &testFrameModel{viewText: "page-2"}
		f.Header(StaticModel("h"))
		f.Content(fm1)
		f.Footer(StaticModel("f"))

		f.Navigate(fm2)
		assert.Contains(t, f.String(), "page-2")

		ok := f.Back()
		assert.True(t, ok)
		assert.Contains(t, f.String(), "page-1")
	})
}
```

**Step 2: Run tests**

Run: `go test -race -run TestFrameNavigateFrameModel ./pkg/cli/`
Expected: PASS

**Step 3: Commit**

```bash
git add pkg/cli/frame_test.go
git commit -m "test(frame): add Navigate/Back tests with FrameModel"
```

---

### Task 11: Add message routing edge case tests

**Files:**
- Test: `pkg/cli/frame_test.go`

**Step 1: Write the tests**

Add to `pkg/cli/frame_test.go`:

```go
func TestFrameMessageRouting_Good(t *testing.T) {
	t.Run("custom message broadcasts to all FrameModels", func(t *testing.T) {
		f := NewFrame("HCF")
		header := &testFrameModel{viewText: "h"}
		content := &testFrameModel{viewText: "c"}
		footer := &testFrameModel{viewText: "f"}
		f.Header(header)
		f.Content(content)
		f.Footer(footer)

		// Send a custom message (not KeyMsg, not WindowSizeMsg)
		type customMsg struct{ data string }
		f.Update(customMsg{data: "hello"})

		assert.True(t, header.updateCalled, "header should receive custom msg")
		assert.True(t, content.updateCalled, "content should receive custom msg")
		assert.True(t, footer.updateCalled, "footer should receive custom msg")
	})

	t.Run("plain Model regions ignore messages gracefully", func(t *testing.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		// Should not panic — modelAdapter ignores all messages
		assert.NotPanics(t, func() {
			f.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			f.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
		})
	})
}
```

**Step 2: Run tests**

Run: `go test -race -run TestFrameMessageRouting ./pkg/cli/`
Expected: PASS

**Step 3: Commit**

```bash
git add pkg/cli/frame_test.go
git commit -m "test(frame): add message routing edge case tests"
```

---

### Task 12: Final verification and cleanup

**Files:**
- All modified files

**Step 1: Run full test suite with race detector**

Run: `go test -race ./pkg/cli/`
Expected: PASS with no race conditions

**Step 2: Run go vet**

Run: `go vet ./pkg/cli/`
Expected: no warnings

**Step 3: Check for unused imports**

Run: `go build ./pkg/cli/`
Expected: clean build, no errors

**Step 4: Count test coverage**

Run: `go test -race -cover ./pkg/cli/`
Expected: coverage reported, aim for >85% on frame.go and frame_model.go

**Step 5: Verify all existing tests still pass with exact same assertions**

Run: `go test -race -v -run "TestFrame_Good|TestFrame_Bad|TestStatusLine|TestKeyHints|TestBreadcrumb|TestStaticModel" ./pkg/cli/`
Expected: all 14 original tests PASS

**Step 6: Final commit if any cleanup was needed**

```bash
git add -A
git commit -m "chore(frame): final cleanup after bubbletea upgrade"
```

---

## Summary of Files

| File | Action | Lines (approx) |
|------|--------|----------------|
| `go.mod` | modify | +2 deps |
| `go.sum` | modify | auto-generated |
| `pkg/cli/frame_model.go` | **create** | ~60 lines |
| `pkg/cli/frame.go` | modify | Replace runLive/renderFrame, add Init/Update/View, add focus helpers (~150 lines changed) |
| `pkg/cli/frame_test.go` | modify | Add ~200 lines of new tests |

## Commit Sequence

1. `deps: add bubbletea and lipgloss for Frame upgrade`
2. `feat(frame): add FrameModel interface and modelAdapter`
3. `feat(frame): add KeyMap with default bindings`
4. `feat(frame): add focus management fields, Focused(), Focus(), WithKeyMap()`
5. `feat(frame): implement tea.Model (Init, Update, View) with lipgloss layout`
6. `feat(frame): replace raw ANSI runLive with tea.Program`
7. `refactor(frame): remove unused renderFrame method`
8. `refactor(frame): unify String() with View() via viewLocked()`
9. `test(frame): add spatial focus navigation tests`
10. `test(frame): add Navigate/Back tests with FrameModel`
11. `test(frame): add message routing edge case tests`
12. `chore(frame): final cleanup after bubbletea upgrade`
