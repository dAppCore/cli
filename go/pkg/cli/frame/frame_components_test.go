package frame

import (
	core "dappco.re/go"
)

func TestFrameComponents_StatusLine_Good(t *core.T) {
	got := StatusLine("core", "ready").View(80, 1)

	core.AssertContains(t, got, "core")
	core.AssertContains(t, got, "ready")
}

func TestFrameComponents_StatusLine_Bad(t *core.T) {
	got := StatusLine("").View(80, 1)

	core.AssertNotContains(t, got, "ready")
	core.AssertEqual(t, "", got)
}

func TestFrameComponents_StatusLine_Ugly(t *core.T) {
	got := StatusLine("abcdef").View(3, 1)

	core.AssertEqual(t, "abc", got)
	core.AssertLen(t, got, 3)
}

func TestFrameComponents_LineModel_View_Good(t *core.T) {
	m := &statusLineModel{title: "core", pairs: []string{"ready"}}

	core.AssertContains(t, m.View(80, 1), "core")
	core.AssertContains(t, m.View(80, 1), "ready")
}

func TestFrameComponents_LineModel_View_Bad(t *core.T) {
	m := &statusLineModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 1))
}

func TestFrameComponents_LineModel_View_Ugly(t *core.T) {
	m := &statusLineModel{title: ":check:"}

	core.AssertEqual(t, "✓", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "✓")
}

func TestFrameComponents_KeyHints_Good(t *core.T) {
	got := KeyHints("q quit", "enter open").View(80, 1)

	core.AssertContains(t, got, "q quit")
	core.AssertContains(t, got, "enter open")
}

func TestFrameComponents_KeyHints_Bad(t *core.T) {
	got := KeyHints().View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestFrameComponents_KeyHints_Ugly(t *core.T) {
	got := KeyHints("abcdef").View(2, 1)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestFrameComponents_HintsModel_View_Good(t *core.T) {
	m := &keyHintsModel{hints: []string{"tab focus"}}

	core.AssertContains(t, m.View(80, 1), "tab")
	core.AssertContains(t, m.View(80, 1), "focus")
}

func TestFrameComponents_HintsModel_View_Bad(t *core.T) {
	m := &keyHintsModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 1))
}

func TestFrameComponents_HintsModel_View_Ugly(t *core.T) {
	m := &keyHintsModel{hints: []string{":cross:"}}

	core.AssertEqual(t, "✗", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "✗")
}

func TestFrameComponents_Breadcrumb_Good(t *core.T) {
	got := Breadcrumb("core", "dev").View(80, 1)

	core.AssertContains(t, got, "core")
	core.AssertContains(t, got, "dev")
}

func TestFrameComponents_Breadcrumb_Bad(t *core.T) {
	got := Breadcrumb().View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestFrameComponents_Breadcrumb_Ugly(t *core.T) {
	got := Breadcrumb("abcdef").View(3, 1)

	core.AssertEqual(t, "abc", got)
	core.AssertLen(t, got, 3)
}

func TestFrameComponents_StaticModel_Good(t *core.T) {
	m := StaticModel("content")
	got := m.View(80, 1)

	core.AssertEqual(t, "content", got)
	core.AssertContains(t, got, "content")
}

func TestFrameComponents_StaticModel_Bad(t *core.T) {
	m := StaticModel("")
	got := m.View(80, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestFrameComponents_StaticModel_Ugly(t *core.T) {
	m := StaticModel(":check:")
	got := m.View(80, 1)

	core.AssertEqual(t, "✓", got)
	core.AssertContains(t, got, "✓")
}

func TestFrameComponents_Model_View_Good(t *core.T) {
	m := &staticModel{text: "content"}

	core.AssertEqual(t, "content", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "content")
}

func TestFrameComponents_Model_View_Bad(t *core.T) {
	m := &staticModel{}

	core.AssertEqual(t, "", m.View(80, 1))
	core.AssertEmpty(t, m.View(0, 0))
}

func TestFrameComponents_Model_View_Ugly(t *core.T) {
	m := &staticModel{text: ":warn:"}

	core.AssertEqual(t, "⚠", m.View(80, 1))
	core.AssertContains(t, m.View(80, 1), "⚠")
}
