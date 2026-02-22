package cli

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrame_Good(t *testing.T) {
	t.Run("static render HCF", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("header"))
		f.Content(StaticModel("content"))
		f.Footer(StaticModel("footer"))

		out := f.String()
		assert.Contains(t, out, "header")
		assert.Contains(t, out, "content")
		assert.Contains(t, out, "footer")
	})

	t.Run("region order preserved", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("AAA"))
		f.Content(StaticModel("BBB"))
		f.Footer(StaticModel("CCC"))

		out := f.String()
		posA := indexOf(out, "AAA")
		posB := indexOf(out, "BBB")
		posC := indexOf(out, "CCC")
		assert.Less(t, posA, posB, "header before content")
		assert.Less(t, posB, posC, "content before footer")
	})

	t.Run("navigate and back", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("nav"))
		f.Content(StaticModel("page-1"))
		f.Footer(StaticModel("hints"))

		assert.Contains(t, f.String(), "page-1")

		// Navigate to page 2
		f.Navigate(StaticModel("page-2"))
		assert.Contains(t, f.String(), "page-2")
		assert.NotContains(t, f.String(), "page-1")

		// Navigate to page 3
		f.Navigate(StaticModel("page-3"))
		assert.Contains(t, f.String(), "page-3")

		// Back to page 2
		ok := f.Back()
		require.True(t, ok)
		assert.Contains(t, f.String(), "page-2")

		// Back to page 1
		ok = f.Back()
		require.True(t, ok)
		assert.Contains(t, f.String(), "page-1")

		// No more history
		ok = f.Back()
		assert.False(t, ok)
	})

	t.Run("empty regions skipped", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("only content"))

		out := f.String()
		assert.Equal(t, "only content\n", out)
	})

	t.Run("non-TTY run renders once", func(t *testing.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		var buf bytes.Buffer
		f := NewFrame("HCF")
		f.out = &buf
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		f.Run() // non-TTY, should return immediately
		assert.Contains(t, buf.String(), "h")
		assert.Contains(t, buf.String(), "c")
		assert.Contains(t, buf.String(), "f")
	})

	t.Run("ModelFunc adapter", func(t *testing.T) {
		called := false
		m := ModelFunc(func(w, h int) string {
			called = true
			return "dynamic"
		})

		out := m.View(80, 24)
		assert.True(t, called)
		assert.Equal(t, "dynamic", out)
	})

	t.Run("RunFor exits after duration", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewFrame("C")
		f.out = &buf // non-TTY → RunFor renders once and returns
		f.Content(StaticModel("timed"))

		start := time.Now()
		f.RunFor(50 * time.Millisecond)
		elapsed := time.Since(start)

		assert.Less(t, elapsed, 200*time.Millisecond)
		assert.Contains(t, buf.String(), "timed")
	})
}

func TestFrame_Bad(t *testing.T) {
	t.Run("empty frame", func(t *testing.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		assert.Equal(t, "", f.String())
	})

	t.Run("back on empty history", func(t *testing.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))
		assert.False(t, f.Back())
	})

	t.Run("invalid variant degrades gracefully", func(t *testing.T) {
		f := NewFrame("XYZ")
		f.out = &bytes.Buffer{}
		// No valid regions, so nothing renders
		assert.Equal(t, "", f.String())
	})
}

func TestStatusLine_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := StatusLine("core dev", "18 repos", "main")
	out := m.View(80, 1)
	assert.Contains(t, out, "core dev")
	assert.Contains(t, out, "18 repos")
	assert.Contains(t, out, "main")
}

func TestKeyHints_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := KeyHints("↑/↓ navigate", "q quit")
	out := m.View(80, 1)
	assert.Contains(t, out, "navigate")
	assert.Contains(t, out, "quit")
}

func TestBreadcrumb_Good(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := Breadcrumb("core", "dev", "health")
	out := m.View(80, 1)
	assert.Contains(t, out, "core")
	assert.Contains(t, out, "dev")
	assert.Contains(t, out, "health")
	assert.Contains(t, out, ">")
}

func TestStaticModel_Good(t *testing.T) {
	m := StaticModel("hello")
	assert.Equal(t, "hello", m.View(80, 24))
}

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

// indexOf returns the position of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := range len(s) - len(substr) + 1 {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
