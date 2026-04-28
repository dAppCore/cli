package frame

import (
	"bytes"
	"strings"
	"time"

	"dappco.re/go"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCli_Frame_Good(t *core.T) {
	t.Run("static render HCF", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("header"))
		f.Content(StaticModel("content"))
		f.Footer(StaticModel("footer"))

		out := f.String()
		core.AssertContains(t, out, "header")
		core.AssertContains(t, out, "content")
		core.AssertContains(t, out, "footer")
	})

	t.Run("region order preserved", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("AAA"))
		f.Content(StaticModel("BBB"))
		f.Footer(StaticModel("CCC"))

		out := f.String()
		posA := strings.Index(out, "AAA")
		posB := strings.Index(out, "BBB")
		posC := strings.Index(out, "CCC")
		core.AssertLess(t, posA, posB, "header before content")
		core.AssertLess(t, posB, posC, "content before footer")
	})

	t.Run("navigate and back", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Header(StaticModel("nav"))
		f.Content(StaticModel("page-1"))
		f.Footer(StaticModel("hints"))
		core.AssertContains(t, f.String(), "page-1")

		// Navigate to page 2
		f.Navigate(StaticModel("page-2"))
		core.AssertContains(t, f.String(), "page-2")
		core.AssertNotContains(t, f.String(), "page-1")

		// Navigate to page 3
		f.Navigate(StaticModel("page-3"))
		core.AssertContains(t, f.String(), "page-3")

		// Back to page 2
		ok := f.Back()
		core.RequireTrue(t, ok)
		core.AssertContains(t, f.String(), "page-2")

		// Back to page 1
		ok = f.Back()
		core.RequireTrue(t, ok)
		core.AssertContains(t, f.String(), "page-1")

		// No more history
		ok = f.Back()
		core.AssertFalse(t, ok)
	})

	t.Run("empty regions skipped", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("only content"))

		out := f.String()
		core.AssertEqual(t, "only content\n", out)
	})

	t.Run("non-TTY run renders once", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		var buf bytes.Buffer
		f := NewFrame("HCF")
		f.out = &buf
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))

		f.Run()
		core.AssertContains( // non-TTY, should return immediately
			t, buf.String(), "h")
		core.AssertContains(t, buf.String(), "c")
		core.AssertContains(t, buf.String(), "f")
	})

	t.Run("ModelFunc adapter", func(t *core.T) {
		called := false
		m := ModelFunc(func(w, h int) string {
			called = true
			return "dynamic"
		})

		out := m.View(80, 24)
		core.AssertTrue(t, called)
		core.AssertEqual(t, "dynamic", out)
	})

	t.Run("RunFor exits after duration", func(t *core.T) {
		var buf bytes.Buffer
		f := NewFrame("C")
		f.out = &buf // non-TTY → RunFor renders once and returns
		f.Content(StaticModel("timed"))

		start := time.Now()
		f.RunFor(50 * time.Millisecond)
		elapsed := time.Since(start)
		core.AssertLess(t, elapsed, 200*time.Millisecond)
		core.AssertContains(t, buf.String(), "timed")
	})
}

func TestCli_Frame_Bad(t *core.T) {
	t.Run("empty frame", func(t *core.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		core.AssertEqual(t, "", f.String())
	})

	t.Run("back on empty history", func(t *core.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))
		core.AssertFalse(t, f.Back())
	})

	t.Run("invalid variant degrades gracefully", func(t *core.T) {
		f := NewFrame("XYZ")
		f.out = &bytes.Buffer{}
		core.AssertEqual( // No valid regions, so nothing renders
			t, "", f.String())
	})
}

func TestStatusLine_Good(t *core.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := StatusLine("core dev", "18 repos", "main")
	out := m.View(80, 1)
	core.AssertContains(t, out, "core dev")
	core.AssertContains(t, out, "18 repos")
	core.AssertContains(t, out, "main")
}

func TestKeyHints_Good(t *core.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := KeyHints("↑/↓ navigate", "q quit")
	out := m.View(80, 1)
	core.AssertContains(t, out, "navigate")
	core.AssertContains(t, out, "quit")
}

func TestBreadcrumb_Good(t *core.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	m := Breadcrumb("core", "dev", "health")
	out := m.View(80, 1)
	core.AssertContains(t, out, "core")
	core.AssertContains(t, out, "dev")
	core.AssertContains(t, out, "health")
	core.AssertContains(t, out, ">")
}

func TestStaticModel_Good(t *core.T) {
	m := StaticModel("hello")
	out := m.View(80, 24)
	core.AssertEqual(t, "hello", out)
	core.AssertContains(t, out, "hello")
}

func TestFrameModel_Good(t *core.T) {
	t.Run("modelAdapter wraps plain Model", func(t *core.T) {
		m := StaticModel("hello")
		adapted := adaptModel(m)

		// Should return nil cmd from Init
		cmd := adapted.Init()
		core.AssertNil(t, cmd)

		// Should return itself from Update
		updated, cmd := adapted.Update(nil)
		core.AssertEqual(t, adapted, updated)
		core.AssertNil(t, cmd)
		core.AssertEqual( // Should delegate View to wrapped model
			t, "hello", adapted.View(80, 24))
	})

	t.Run("FrameModel passes through without wrapping", func(t *core.T) {
		fm := &testFrameModel{viewText: "interactive"}
		adapted := adaptModel(fm)

		// Should be the same object, not wrapped
		_, ok := adapted.(*testFrameModel)
		core.AssertTrue(t, ok, "FrameModel should not be wrapped")
		core.AssertEqual(t, "interactive", adapted.View(80, 24))
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

func TestKeyMap_Good(t *core.T) {
	t.Run("default keymap has expected bindings", func(t *core.T) {
		km := DefaultKeyMap()
		core.AssertEqual(t, tea.KeyTab, km.FocusNext)
		core.AssertEqual(t, tea.KeyShiftTab, km.FocusPrev)
		core.AssertEqual(t, tea.KeyUp, km.FocusUp)
		core.AssertEqual(t, tea.KeyDown, km.FocusDown)
		core.AssertEqual(t, tea.KeyLeft, km.FocusLeft)
		core.AssertEqual(t, tea.KeyRight, km.FocusRight)
		core.AssertEqual(t, tea.KeyEsc, km.Back)
		core.AssertEqual(t, tea.KeyCtrlC, km.Quit)
	})
}

func TestFrameFocus_Good(t *core.T) {
	t.Run("default focus is Content", func(t *core.T) {
		f := NewFrame("HCF")
		core.AssertEqual(t, RegionContent, f.Focused())
	})

	t.Run("Focus sets focused region", func(t *core.T) {
		f := NewFrame("HCF")
		f.Focus(RegionHeader)
		core.AssertEqual(t, RegionHeader, f.Focused())
	})

	t.Run("Focus ignores invalid region", func(t *core.T) {
		f := NewFrame("HCF")
		f.Focus(RegionLeft)
		core.AssertEqual( // Left not in "HCF"
			t, RegionContent, f.Focused()) // unchanged
	})

	t.Run("WithKeyMap returns frame for chaining", func(t *core.T) {
		km := DefaultKeyMap()
		km.Quit = tea.KeyCtrlQ
		f := NewFrame("HCF").WithKeyMap(km)
		core.AssertEqual(t, tea.KeyCtrlQ, f.keyMap.Quit)
	})

	t.Run("focusRing builds from variant", func(t *core.T) {
		f := NewFrame("HLCRF")
		ring := f.buildFocusRing()
		core.AssertEqual(t, []Region{RegionHeader, RegionLeft, RegionContent, RegionRight, RegionFooter}, ring)
	})

	t.Run("focusRing respects variant order", func(t *core.T) {
		f := NewFrame("HCF")
		ring := f.buildFocusRing()
		core.AssertEqual(t, []Region{RegionHeader, RegionContent, RegionFooter}, ring)
	})
}

func TestFrameTeaModel_Good(t *core.T) {
	t.Run("Init collects FrameModel inits", func(t *core.T) {
		f := NewFrame("HCF")
		fm := &testFrameModel{viewText: "x"}
		f.Content(fm)

		cmd := f.Init()
		// Should produce a batch command (non-nil if any FrameModel has Init)
		// fm.Init returns nil, so batch of nils = nil
		_ = cmd
		core.AssertTrue( // no panic = success
			t, fm.initCalled)
	})

	t.Run("Update routes key to focused region", func(t *core.T) {
		f := NewFrame("HCF")
		header := &testFrameModel{viewText: "h"}
		content := &testFrameModel{viewText: "c"}
		f.Header(header)
		f.Content(content)

		// Focus is Content by default
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		f.Update(keyMsg)
		core.AssertTrue(t, content.updateCalled, "focused region should receive key")
		core.AssertFalse(t, header.updateCalled, "unfocused region should not receive key")
	})

	t.Run("Update broadcasts WindowSizeMsg to all", func(t *core.T) {
		f := NewFrame("HCF")
		header := &testFrameModel{viewText: "h"}
		content := &testFrameModel{viewText: "c"}
		footer := &testFrameModel{viewText: "f"}
		f.Header(header)
		f.Content(content)
		f.Footer(footer)

		sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
		f.Update(sizeMsg)
		core.AssertTrue(t, header.updateCalled, "header should get resize")
		core.AssertTrue(t, content.updateCalled, "content should get resize")
		core.AssertTrue(t, footer.updateCalled, "footer should get resize")
		core.AssertEqual(t, 120, f.width)
		core.AssertEqual(t, 40, f.height)
	})

	t.Run("Update handles quit key", func(t *core.T) {
		f := NewFrame("HCF")
		f.Content(StaticModel("c"))

		quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
		_, cmd := f.Update(quitMsg)
		core.AssertNotNil( // cmd should be tea.Quit
			t, cmd)
	})

	t.Run("Update handles back key", func(t *core.T) {
		f := NewFrame("HCF")
		f.Content(StaticModel("page-1"))
		f.Navigate(StaticModel("page-2"))

		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		f.Update(escMsg)
		core.AssertContains(t, f.String(), "page-1")
	})

	t.Run("Update cycles focus with Tab", func(t *core.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))
		core.AssertEqual(t, RegionContent, f.Focused())

		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		f.Update(tabMsg)
		core.AssertEqual(t, RegionFooter, f.Focused())

		f.Update(tabMsg)
		core.AssertEqual(t, RegionHeader, f.Focused()) // wraps around

		shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
		f.Update(shiftTabMsg)
		core.AssertEqual(t, RegionFooter, f.Focused()) // back
	})

	t.Run("View produces non-empty output", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.Header(StaticModel("HEAD"))
		f.Content(StaticModel("BODY"))
		f.Footer(StaticModel("FOOT"))

		view := f.View()
		core.AssertContains(t, view, "HEAD")
		core.AssertContains(t, view, "BODY")
		core.AssertContains(t, view, "FOOT")
	})

	t.Run("View lipgloss layout: header before content before footer", func(t *core.T) {
		SetColorEnabled(false)
		defer SetColorEnabled(true)

		f := NewFrame("HCF")
		f.Header(StaticModel("AAA"))
		f.Content(StaticModel("BBB"))
		f.Footer(StaticModel("CCC"))
		f.width = 80
		f.height = 24

		view := f.View()
		posA := strings.Index(view, "AAA")
		posB := strings.Index(view, "BBB")
		posC := strings.Index(view, "CCC")
		core.AssertGreater(t, posA, -1, "header should be present")
		core.AssertGreater(t, posB, -1, "content should be present")
		core.AssertGreater(t, posC, -1, "footer should be present")
		core.AssertLess(t, posA, posB, "header before content")
		core.AssertLess(t, posB, posC, "content before footer")
	})
}

func TestFrameSend_Good(t *core.T) {
	t.Run("Send is safe before Run", func(t *core.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))
		core.AssertNotPanics( // Should not panic when program is nil
			t, func() {
				f.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			})
	})
}

func TestFrameSpatialFocus_Good(t *core.T) {
	t.Run("arrow keys move to target region", func(t *core.T) {
		f := NewFrame("HLCRF")
		f.Header(StaticModel("h"))
		f.Left(StaticModel("l"))
		f.Content(StaticModel("c"))
		f.Right(StaticModel("r"))
		f.Footer(StaticModel("f"))
		core.AssertEqual( // Start at Content
			t, RegionContent, f.Focused())

		// Up → Header
		f.Update(tea.KeyMsg{Type: tea.KeyUp})
		core.AssertEqual(t, RegionHeader, f.Focused())

		// Down → Footer
		f.Update(tea.KeyMsg{Type: tea.KeyDown})
		core.AssertEqual(t, RegionFooter, f.Focused())

		// Left → Left sidebar
		f.Update(tea.KeyMsg{Type: tea.KeyLeft})
		core.AssertEqual(t, RegionLeft, f.Focused())

		// Right → Right sidebar
		f.Update(tea.KeyMsg{Type: tea.KeyRight})
		core.AssertEqual(t, RegionRight, f.Focused())
	})

	t.Run("spatial focus ignores missing regions", func(t *core.T) {
		f := NewFrame("HCF") // no Left or Right
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))
		core.AssertEqual(t, RegionContent, f.Focused())

		// Left arrow → no Left region, focus stays
		f.Update(tea.KeyMsg{Type: tea.KeyLeft})
		core.AssertEqual(t, RegionContent, f.Focused())
	})
}

func TestFrameNavigateFrameModel_Good(t *core.T) {
	t.Run("Navigate preserves current focus", func(t *core.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(&testFrameModel{viewText: "page-1"})
		f.Footer(StaticModel("f"))

		// Focus something else
		f.Focus(RegionHeader)
		core.AssertEqual(t, RegionHeader, f.Focused())

		// Navigate replaces Content, focus should remain where it was
		f.Navigate(&testFrameModel{viewText: "page-2"})
		core.AssertEqual(t, RegionHeader, f.Focused())
		core.AssertContains(t, f.String(), "page-2")
	})

	t.Run("Back restores FrameModel", func(t *core.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		fm1 := &testFrameModel{viewText: "page-1"}
		fm2 := &testFrameModel{viewText: "page-2"}
		f.Header(StaticModel("h"))
		f.Content(fm1)
		f.Footer(StaticModel("f"))

		f.Navigate(fm2)
		core.AssertContains(t, f.String(), "page-2")

		ok := f.Back()
		core.AssertTrue(t, ok)
		core.AssertContains(t, f.String(), "page-1")
	})
}

func TestFrameMessageRouting_Good(t *core.T) {
	t.Run("custom message broadcasts to all FrameModels", func(t *core.T) {
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
		core.AssertTrue(t, header.updateCalled, "header should receive custom msg")
		core.AssertTrue(t, content.updateCalled, "content should receive custom msg")
		core.AssertTrue(t, footer.updateCalled, "footer should receive custom msg")
	})

	t.Run("plain Model regions ignore messages gracefully", func(t *core.T) {
		f := NewFrame("HCF")
		f.Header(StaticModel("h"))
		f.Content(StaticModel("c"))
		f.Footer(StaticModel("f"))
		core.AssertNotPanics( // Should not panic — modelAdapter ignores all messages
			t, func() {
				f.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
				f.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			})
	})
}

func TestCli_Frame_Ugly(t *core.T) {
	t.Run("navigate with nil model does not panic", func(t *core.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("base"))
		core.AssertNotPanics(t, func() {
			f.Navigate(nil)
		})
	})

	t.Run("deeply nested back stack does not panic", func(t *core.T) {
		f := NewFrame("C")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("p0"))
		for i := 1; i <= 20; i++ {
			f.Navigate(StaticModel("p" + string(rune('0'+i%10))))
		}
		for f.Back() {
			// drain the full history stack
		}
		core.AssertFalse(t, f.Back(), "no more history after full drain")
	})

	t.Run("zero-size window renders without panic", func(t *core.T) {
		f := NewFrame("HCF")
		f.out = &bytes.Buffer{}
		f.Content(StaticModel("x"))
		f.width = 0
		f.height = 0
		core.AssertNotPanics(t, func() {
			_ = f.View()
		})
	})
}
