package cli

import (
	core "dappco.re/go"
)

func TestAnsi_AnsiStyle_Italic_Good(t *core.T) {
	s := NewStyle().Italic()

	core.AssertTrue(t, s.italic)
	core.AssertEqual(t, s, s.Italic())
}

func TestAnsi_AnsiStyle_Italic_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Italic() })
	core.AssertNil(t, s)
}

func TestAnsi_AnsiStyle_Italic_Ugly(t *core.T) {
	s := NewStyle().Italic().Italic()

	core.AssertTrue(t, s.italic)
	core.AssertFalse(t, s.underline)
}

func TestAnsi_AnsiStyle_Underline_Good(t *core.T) {
	s := NewStyle().Underline()

	core.AssertTrue(t, s.underline)
	core.AssertEqual(t, s, s.Underline())
}

func TestAnsi_AnsiStyle_Underline_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Underline() })
	core.AssertNil(t, s)
}

func TestAnsi_AnsiStyle_Underline_Ugly(t *core.T) {
	s := NewStyle().Underline().Underline()

	core.AssertTrue(t, s.underline)
	core.AssertFalse(t, s.bold)
}

func TestAnsi_AnsiStyle_Background_Good(t *core.T) {
	s := NewStyle().Background("#0000ff")

	core.AssertContains(t, s.bg, "48;2;0;0;255")
	core.AssertEqual(t, s, s.Background("#00ff00"))
}

func TestAnsi_AnsiStyle_Background_Bad(t *core.T) {
	s := NewStyle().Background("bad")

	core.AssertContains(t, s.bg, "255;255;255")
	core.AssertNotEmpty(t, s.bg)
}

func TestAnsi_AnsiStyle_Background_Ugly(t *core.T) {
	s := NewStyle().Background("")

	core.AssertContains(t, s.bg, "255;255;255")
	core.AssertNotEmpty(t, s.bg)
}

func TestAnsi_ColorEnabled_Good(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(true)

	core.AssertTrue(t, ColorEnabled())
	core.AssertContains(t, NewStyle().Bold().Render("x"), "\033[")
}

func TestAnsi_ColorEnabled_Bad(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(false)

	core.AssertFalse(t, ColorEnabled())
	core.AssertEqual(t, "x", NewStyle().Bold().Render("x"))
}

func TestAnsi_ColorEnabled_Ugly(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(false)
	SetColorEnabled(true)

	core.AssertTrue(t, ColorEnabled())
	core.AssertNotEqual(t, "x", NewStyle().Bold().Render("x"))
}

func TestAnsi_SetColorEnabled_Good(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(false)

	core.AssertFalse(t, ColorEnabled())
	core.AssertEqual(t, "plain", NewStyle().Foreground("#00ff00").Render("plain"))
}

func TestAnsi_SetColorEnabled_Bad(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(true)

	core.AssertTrue(t, ColorEnabled())
	core.AssertContains(t, NewStyle().Foreground("#00ff00").Render("plain"), "38;2;0;255;0")
}

func TestAnsi_SetColorEnabled_Ugly(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(false)
	SetColorEnabled(false)

	core.AssertFalse(t, ColorEnabled())
	core.AssertEqual(t, "plain", NewStyle().Render("plain"))
}

func TestAnsi_NewStyle_Good(t *core.T) {
	s := NewStyle()

	core.AssertNotNil(t, s)
	core.AssertEqual(t, "plain", s.Render("plain"))
}

func TestAnsi_NewStyle_Bad(t *core.T) {
	s := NewStyle()

	core.AssertFalse(t, s.bold)
	core.AssertEmpty(t, s.fg)
}

func TestAnsi_NewStyle_Ugly(t *core.T) {
	s := NewStyle().Bold().Dim()

	core.AssertTrue(t, s.bold)
	core.AssertTrue(t, s.dim)
}

func TestAnsi_AnsiStyle_Bold_Good(t *core.T) {
	s := NewStyle().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertEqual(t, s, s.Bold())
}

func TestAnsi_AnsiStyle_Bold_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Bold() })
	core.AssertNil(t, s)
}

func TestAnsi_AnsiStyle_Bold_Ugly(t *core.T) {
	s := NewStyle().Bold().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestAnsi_AnsiStyle_Dim_Good(t *core.T) {
	s := NewStyle().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertEqual(t, s, s.Dim())
}

func TestAnsi_AnsiStyle_Dim_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Dim() })
	core.AssertNil(t, s)
}

func TestAnsi_AnsiStyle_Dim_Ugly(t *core.T) {
	s := NewStyle().Dim().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertFalse(t, s.bold)
}

func TestAnsi_AnsiStyle_Foreground_Good(t *core.T) {
	s := NewStyle().Foreground("#0000ff")

	core.AssertContains(t, s.fg, "38;2;0;0;255")
	core.AssertEqual(t, s, s.Foreground("#00ff00"))
}

func TestAnsi_AnsiStyle_Foreground_Bad(t *core.T) {
	s := NewStyle().Foreground("bad")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAnsi_AnsiStyle_Foreground_Ugly(t *core.T) {
	s := NewStyle().Foreground("")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAnsi_AnsiStyle_Render_Good(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(true)

	got := NewStyle().Bold().Render("ready")
	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, ansiReset)
}

func TestAnsi_AnsiStyle_Render_Bad(t *core.T) {
	old := ColorEnabled()
	t.Cleanup(func() { SetColorEnabled(old) })
	SetColorEnabled(false)

	got := NewStyle().Bold().Render("ready")
	core.AssertEqual(t, "ready", got)
	core.AssertNotContains(t, got, ansiReset)
}

func TestAnsi_AnsiStyle_Render_Ugly(t *core.T) {
	got := NewStyle().Render("ready")

	core.AssertEqual(t, "ready", got)
	core.AssertNotContains(t, got, ansiReset)
}
