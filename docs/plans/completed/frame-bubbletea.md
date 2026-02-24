# Frame Bubbletea Upgrade — Completion Summary

**Completed:** 22 February 2026
**Module:** `forge.lthn.ai/core/cli`
**Status:** Complete — Frame implements tea.Model with full bubbletea lifecycle

## What Was Built

Upgraded the Frame layout system from a static HLCRF renderer to a full
bubbletea `tea.Model` with lifecycle management, keyboard handling, and
panel navigation.

### Key changes

- **Frame implements `tea.Model`** — `Init()`, `Update()`, `View()` lifecycle
- **`KeyMap`** — configurable keybindings with default set (quit, navigate,
  help, focus cycling)
- **`Navigate(name)` / `Back()`** — panel switching with history stack
- **Focus management** — Tab/Shift-Tab cycles focus between visible models
- **lipgloss layout** — HLCRF regions (Header, Left, Content, Right, Footer)
  rendered with lipgloss instead of raw ANSI
- **`FrameModel` interface** — models register with `Frame.Header()`,
  `.Content()`, `.Footer()` etc., receiving focus/blur/resize messages

### Tests

Navigate/Back stack tests, focus cycling, key dispatch, resize propagation.
All passing with `-race`.

### Dependencies

- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`

### Consumer

`go-blockchain/cmd/chain/` is the first consumer — TUI dashboard uses
Frame with StatusModel (header), ExplorerModel (content), KeyHintsModel
(footer).
