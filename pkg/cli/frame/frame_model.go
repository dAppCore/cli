package frame

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

func (a *modelAdapter) View(w, h int) string                 { return a.m.View(w, h) }
func (a *modelAdapter) Init() tea.Cmd                        { return nil }
func (a *modelAdapter) Update(tea.Msg) (FrameModel, tea.Cmd) { return a, nil }

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
