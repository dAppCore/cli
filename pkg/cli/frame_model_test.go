package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestFrameModel_AdaptModel_Good(t *testing.T) {
	// A plain Model is wrapped in a modelAdapter.
	m := StaticModel("wrapped")
	adapted := adaptModel(m)

	if adapted == nil {
		t.Fatal("adaptModel returned nil for a valid Model")
	}
	assert.Equal(t, "wrapped", adapted.View(80, 24))
}

func TestFrameModel_AdaptModel_Bad(t *testing.T) {
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

func TestFrameModel_AdaptModel_Ugly(t *testing.T) {
	// modelAdapter's Init, Update, and View must not panic on edge inputs.
	m := StaticModel("edge")
	adapted := adaptModel(m)

	assert.NotPanics(t, func() {
		_ = adapted.Init()
	})
	assert.NotPanics(t, func() {
		_, _ = adapted.Update(nil)
	})
	assert.NotPanics(t, func() {
		_ = adapted.View(-1, -1)
	})
	assert.NotPanics(t, func() {
		_ = adapted.View(0, 0)
	})
}

func TestFrameModel_DefaultKeyMap_Good(t *testing.T) {
	// DefaultKeyMap must return the expected standard bindings.
	km := DefaultKeyMap()

	assert.Equal(t, tea.KeyTab, km.FocusNext)
	assert.Equal(t, tea.KeyShiftTab, km.FocusPrev)
	assert.Equal(t, tea.KeyUp, km.FocusUp)
	assert.Equal(t, tea.KeyDown, km.FocusDown)
	assert.Equal(t, tea.KeyLeft, km.FocusLeft)
	assert.Equal(t, tea.KeyRight, km.FocusRight)
	assert.Equal(t, tea.KeyEsc, km.Back)
	assert.Equal(t, tea.KeyCtrlC, km.Quit)
}

func TestFrameModel_DefaultKeyMap_Bad(t *testing.T) {
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

func TestFrameModel_DefaultKeyMap_Ugly(t *testing.T) {
	// Multiple calls must return identical, independent copies — no shared state.
	a := DefaultKeyMap()
	b := DefaultKeyMap()

	// Same values.
	assert.Equal(t, a, b)

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
