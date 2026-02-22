# Frame Bubbletea Upgrade Design

**Issue:** core/go#15
**Date:** 2026-02-22
**Status:** Approved

**Goal:** Upgrade `cli.Frame` from raw ANSI + `golang.org/x/term` to bubbletea internally, adding keyboard navigation, focus management, and lipgloss layout composition while preserving the existing public API.

---

## Architecture

Single ownership model. Frame becomes the sole `tea.Model` wrapping a `tea.Program`. It owns the terminal (alt-screen, raw mode, resize events, input). Region models never touch the terminal directly.

Message routing:
- **Key messages** — routed to the focused region's `FrameModel.Update()` only
- **Tick/resize messages** — broadcast to all region `FrameModel.Update()` calls
- **Custom messages** — broadcast to all (enables cross-region communication)

Dual interface pattern:

```go
// Existing — view-only, no changes
type Model interface {
    View(width, height int) string
}

// New — interactive components
type FrameModel interface {
    Model
    Init() tea.Cmd
    Update(tea.Msg) (FrameModel, tea.Cmd)
}
```

Frame wraps plain `Model` in a no-op adapter internally, so existing code (StatusLine, KeyHints, Breadcrumb, StaticModel, ModelFunc) works without changes.

Layout composition replaces the manual ANSI cursor/clear dance in `runLive()` with lipgloss `JoinVertical` and `JoinHorizontal`. The existing HLCRF variant parser and region size calculations stay, but rendering uses lipgloss instead of raw escape codes.

---

## Focus Management

Focus ring. Frame maintains an ordered list of focusable regions (only regions with `FrameModel` components). Focus cycles through them.

Navigation:
- `Tab` / `Shift-Tab` — cycle focus forward/backward through the ring
- Arrow keys — spatial navigation (up to Header, down to Footer, left to Left sidebar, right to Right sidebar)
- Configurable via `KeyMap` struct with sensible defaults

```go
type KeyMap struct {
    FocusNext  key.Binding // Tab
    FocusPrev  key.Binding // Shift-Tab
    FocusUp    key.Binding // Up (to Header from Content)
    FocusDown  key.Binding // Down (to Footer from Content)
    FocusLeft  key.Binding // Left (to Left sidebar)
    FocusRight key.Binding // Right (to Right sidebar)
    Quit       key.Binding // q, Ctrl-C
    Back       key.Binding // Esc (triggers Navigate back)
}
```

Visual feedback: focused region gets a subtle border highlight (configurable via lipgloss border styling). Unfocused regions render normally.

Key filtering: focus keys are consumed by Frame and never forwarded to region models. All other keys go to the focused region's `Update()`.

---

## Public API

### Preserved (no changes)

- `NewFrame(variant string) *Frame`
- `Header(m Model)`, `Left(m Model)`, `Content(m Model)`, `Right(m Model)`, `Footer(m Model)`
- `Navigate(m Model)`, `Back() bool`
- `Run()`, `RunFor(d time.Duration)`, `Stop()`
- `String()` — static render for non-TTY
- `ModelFunc`, `StaticModel`, `StatusLine`, `KeyHints`, `Breadcrumb`

### New additions

```go
// WithKeyMap sets custom key bindings for Frame navigation.
func (f *Frame) WithKeyMap(km KeyMap) *Frame

// Focused returns the currently focused region.
func (f *Frame) Focused() Region

// Focus sets focus to a specific region.
func (f *Frame) Focus(r Region)

// Send injects a message into the Frame's tea.Program.
// Useful for triggering updates from external goroutines.
func (f *Frame) Send(msg tea.Msg)
```

### Behavioural changes

- `Run()` now starts a `tea.Program` in TTY mode (instead of raw ticker loop)
- Non-TTY path unchanged — still calls `String()` and returns
- `RunFor()` unchanged — uses `Stop()` after timer

### New dependencies

- `github.com/charmbracelet/bubbletea` (already in core/go)
- `github.com/charmbracelet/lipgloss` (already in core/go)
- `github.com/charmbracelet/bubbles/key` (key bindings)

---

## Internal Implementation

Frame implements `tea.Model`:

```go
func (f *Frame) Init() tea.Cmd
func (f *Frame) Update(tea.Msg) (tea.Model, tea.Cmd)
func (f *Frame) View() string
```

`Init()` collects `Init()` from all `FrameModel` regions via `tea.Batch()`.

`Update()` handles:
1. `tea.WindowSizeMsg` — update dimensions, broadcast to all FrameModels
2. `tea.KeyMsg` matching focus keys — advance/retreat focus ring
3. `tea.KeyMsg` matching quit — return `tea.Quit`
4. `tea.KeyMsg` matching back — call `Back()`, return nil
5. All other `tea.KeyMsg` — forward to focused region's `Update()`
6. All other messages — broadcast to all FrameModels

`View()` uses lipgloss composition:

```
header  = renderRegion(H, width, 1)
footer  = renderRegion(F, width, 1)
middleH = height - headerH - footerH

left    = renderRegion(L, width/4, middleH)
right   = renderRegion(R, width/4, middleH)
content = renderRegion(C, contentW, middleH)

middle  = lipgloss.JoinHorizontal(Top, left, content, right)
output  = lipgloss.JoinVertical(Left, header, middle, footer)
```

`Run()` change:

```go
func (f *Frame) Run() {
    if !f.isTTY() {
        fmt.Fprint(f.out, f.String())
        return
    }
    p := tea.NewProgram(f, tea.WithAltScreen())
    f.program = p
    if _, err := p.Run(); err != nil {
        Fatal(err)
    }
}
```

Plain `Model` adapter:

```go
type modelAdapter struct{ m Model }
func (a *modelAdapter) Init() tea.Cmd                      { return nil }
func (a *modelAdapter) Update(tea.Msg) (FrameModel, tea.Cmd) { return a, nil }
func (a *modelAdapter) View(w, h int) string               { return a.m.View(w, h) }
```

---

## Testing Strategy

Existing 14 tests preserved. They use `bytes.Buffer` (non-TTY path), bypassing bubbletea.

New tests for interactive features:
- Focus cycling: Tab advances focus, Shift-Tab goes back
- Spatial navigation: arrow keys move focus to correct region
- Message routing: key events only reach focused model
- Tick broadcast: tick events reach all models
- Resize propagation: resize reaches all models
- FrameModel lifecycle: Init() called on Run(), Update() receives messages
- Adapter: plain Model wrapped correctly, receives no Update calls
- Navigate/Back with FrameModel: focus transfers correctly
- KeyMap customization: overridden bindings work
- Send(): external messages delivered to models

Testing approach: use bubbletea's `teatest` package for interactive tests. Non-TTY tests stay as-is with `bytes.Buffer`.

---

## Files Affected

| File | Action | Purpose |
|------|--------|---------|
| `pkg/cli/frame.go` | modify | Add bubbletea tea.Model implementation, lipgloss layout, focus management |
| `pkg/cli/frame_model.go` | new | FrameModel interface, modelAdapter, KeyMap |
| `pkg/cli/frame_test.go` | modify | Add interactive tests alongside existing ones |
| `go.mod` | modify | Add bubbletea, lipgloss, bubbles dependencies |

## Design Decisions

1. **Frame as tea.Model, not wrapping separate tea.Model** — Frame IS the model, simplest ownership
2. **Dual interface (Model + FrameModel)** — backward compatible, existing components unchanged
3. **Lipgloss for layout** — replaces manual ANSI, consistent with bubbletea ecosystem
4. **Focus ring with spatial override** — Tab for cycling, arrows for direct spatial jumps
5. **Non-TTY path untouched** — `String()` and non-TTY `Run()` stay exactly as-is
