package frame

import (
	"dappco.re/go"
	tea "github.com/charmbracelet/bubbletea"
)

func TestFrameModel_AdaptModel_Good(t *core.T) {
	// A plain Model is wrapped in a modelAdapter.
	m := StaticModel("wrapped")
	adapted := adaptModel(m)

	if adapted == nil {
		t.Fatal("adaptModel returned nil for a valid Model")
	}
	core.AssertEqual(t, "wrapped", adapted.View(80, 24))
}

func TestFrameModel_AdaptModel_Bad(t *core.T) {
	// A FrameModel passes through unchanged — no double-wrapping.
	fm := &testFrameModelForAdapter{viewText: "native"}
	adapted := adaptModel(fm)

	same, ok := adapted.(*testFrameModelForAdapter)
	if !ok {
		t.Fatalf("FrameModel should not be wrapped, got %T", adapted)
	}
	if same != fm {
		t.Error("FrameModel was wrapped instead of passed through")
	}
}

func TestFrameModel_AdaptModel_Ugly(t *core.T) {
	// modelAdapter's Init, Update, and View must not panic on edge inputs.
	m := StaticModel("edge")
	adapted := adaptModel(m)
	core.AssertNotPanics(t, func() {
		_ = adapted.Init()
	})
	core.AssertNotPanics(t, func() {
		_, _ = adapted.Update(nil)
	})
	core.AssertNotPanics(t, func() {
		_ = adapted.View(-1, -1)
	})
	core.AssertNotPanics(t, func() {
		_ = adapted.View(0, 0)
	})
}

func TestFrameModel_DefaultKeyMap_Good(t *core.T) {
	// DefaultKeyMap must return the expected standard bindings.
	km := DefaultKeyMap()
	core.AssertEqual(t, tea.KeyTab, km.FocusNext)
	core.AssertEqual(t, tea.KeyShiftTab, km.FocusPrev)
	core.AssertEqual(t, tea.KeyUp, km.FocusUp)
	core.AssertEqual(t, tea.KeyDown, km.FocusDown)
	core.AssertEqual(t, tea.KeyLeft, km.FocusLeft)
	core.AssertEqual(t, tea.KeyRight, km.FocusRight)
	core.AssertEqual(t, tea.KeyEsc, km.Back)
	core.AssertEqual(t, tea.KeyCtrlC, km.Quit)
}

func TestFrameModel_DefaultKeyMap_Bad(t *core.T) {
	// The four spatial focus keys must all be distinct from each other.
	km := DefaultKeyMap()
	spatial := []tea.KeyType{km.FocusUp, km.FocusDown, km.FocusLeft, km.FocusRight}
	seen := make(map[tea.KeyType]bool)
	for _, k := range spatial {
		if seen[k] {
			t.Errorf("duplicate spatial binding: %v", k)
		}
		seen[k] = true
	}
}

func TestFrameModel_DefaultKeyMap_Ugly(t *core.T) {
	// Multiple calls must return identical, independent copies — no shared state.
	a := DefaultKeyMap()
	b := DefaultKeyMap()
	core.AssertEqual( // Same values.
		t, a, b)

	// Mutating one does not affect the other (value semantics).
	b.Quit = tea.KeyCtrlD
	if a.Quit == b.Quit {
		t.Error("DefaultKeyMap should return independent copies, mutation leaked")
	}
}

// testFrameModelForAdapter is a minimal FrameModel used only by
// frame_model_test.go to verify adaptModel's pass-through path.
// Named distinctly from testFrameModel in frame_test.go to avoid collisions.
type testFrameModelForAdapter struct {
	viewText string
}

func (m *testFrameModelForAdapter) View(_, _ int) string { return m.viewText }
func (m *testFrameModelForAdapter) Init() tea.Cmd        { return nil }
func (m *testFrameModelForAdapter) Update(_ tea.Msg) (FrameModel, tea.Cmd) {
	return m, nil
}
