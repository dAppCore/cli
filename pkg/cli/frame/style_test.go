package frame

import (
	core "dappco.re/go"
)

func TestStyle_ColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_ColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestStyle_ColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(!original)
	defer SetColorEnabled(original)

	core.AssertEqual(t, !original, ColorEnabled())
}

func TestStyle_SetColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_SetColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestStyle_SetColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_NewStyle_Good(t *core.T) {
	s := NewStyle()

	core.AssertNotNil(t, s)
	core.AssertEqual(t, "plain", s.Render("plain"))
}

func TestStyle_NewStyle_Bad(t *core.T) {
	s := NewStyle()

	core.AssertFalse(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestStyle_NewStyle_Ugly(t *core.T) {
	s := NewStyle().Bold().Dim().Italic().Underline()

	core.AssertTrue(t, s.bold)
	core.AssertTrue(t, s.underline)
}

func TestStyle_AnsiStyle_Bold_Good(t *core.T) {
	s := NewStyle().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertEqual(t, s, s.Bold())
}

func TestStyle_AnsiStyle_Bold_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Bold() })
	core.AssertNil(t, s)
}

func TestStyle_AnsiStyle_Bold_Ugly(t *core.T) {
	s := NewStyle().Bold().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestStyle_AnsiStyle_Dim_Good(t *core.T) {
	s := NewStyle().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertEqual(t, s, s.Dim())
}

func TestStyle_AnsiStyle_Dim_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Dim() })
	core.AssertNil(t, s)
}

func TestStyle_AnsiStyle_Dim_Ugly(t *core.T) {
	s := NewStyle().Dim().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertFalse(t, s.bold)
}

func TestStyle_AnsiStyle_Foreground_Good(t *core.T) {
	s := NewStyle().Foreground("#ff0000")

	core.AssertContains(t, s.fg, "38;2;255;0;0")
	core.AssertEqual(t, s, s.Foreground("#00ff00"))
}

func TestStyle_AnsiStyle_Foreground_Bad(t *core.T) {
	s := NewStyle().Foreground("bad")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestStyle_AnsiStyle_Foreground_Ugly(t *core.T) {
	s := NewStyle().Foreground("")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestStyle_AnsiStyle_Render_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertContains(t, got, "text")
	core.AssertContains(t, got, "\033[1m")
}

func TestStyle_AnsiStyle_Render_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestStyle_AnsiStyle_Render_Ugly(t *core.T) {
	var s *AnsiStyle
	got := s.Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestStyle_Truncate_Good(t *core.T) {
	got := Truncate("abcdef", 4)

	core.AssertEqual(t, "a...", got)
	core.AssertLen(t, got, 4)
}

func TestStyle_Truncate_Bad(t *core.T) {
	got := Truncate("abcdef", 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStyle_Truncate_Ugly(t *core.T) {
	got := Truncate("abcdef", 2)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestStyle_Glyph_Good(t *core.T) {
	UseUnicode()
	got := Glyph(":check:")

	core.AssertEqual(t, "✓", got)
	core.AssertNotEqual(t, ":check:", got)
}

func TestStyle_Glyph_Bad(t *core.T) {
	got := Glyph(":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestStyle_Glyph_Ugly(t *core.T) {
	got := Glyph("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStyle_ColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_ColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestStyle_ColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(!original)
	defer SetColorEnabled(original)

	core.AssertEqual(t, !original, ColorEnabled())
}

func TestStyle_SetColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_SetColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestStyle_SetColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestStyle_NewStyle_Good(t *core.T) {
	s := NewStyle()

	core.AssertNotNil(t, s)
	core.AssertEqual(t, "text", s.Render("text"))
}

func TestStyle_NewStyle_Bad(t *core.T) {
	s := NewStyle()

	core.AssertFalse(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestStyle_NewStyle_Ugly(t *core.T) {
	s := NewStyle().Bold().Dim()

	core.AssertTrue(t, s.bold)
	core.AssertTrue(t, s.dim)
}

func TestStyle_AnsiStyle_Bold_Good(t *core.T) {
	s := NewStyle().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertEqual(t, s, s.Bold())
}

func TestStyle_AnsiStyle_Bold_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Bold() })
	core.AssertNil(t, s)
}

func TestStyle_AnsiStyle_Bold_Ugly(t *core.T) {
	s := NewStyle().Bold().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestStyle_AnsiStyle_Dim_Good(t *core.T) {
	s := NewStyle().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertEqual(t, s, s.Dim())
}

func TestStyle_AnsiStyle_Dim_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Dim() })
	core.AssertNil(t, s)
}

func TestStyle_AnsiStyle_Dim_Ugly(t *core.T) {
	s := NewStyle().Dim().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertFalse(t, s.bold)
}

func TestStyle_AnsiStyle_Foreground_Good(t *core.T) {
	s := NewStyle().Foreground("#ff0000")

	core.AssertContains(t, s.fg, "38;2;255;0;0")
	core.AssertEqual(t, s, s.Foreground("#00ff00"))
}

func TestStyle_AnsiStyle_Foreground_Bad(t *core.T) {
	s := NewStyle().Foreground("bad")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestStyle_AnsiStyle_Foreground_Ugly(t *core.T) {
	s := NewStyle().Foreground("")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestStyle_AnsiStyle_Render_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertContains(t, got, "text")
	core.AssertContains(t, got, "\033[1m")
}

func TestStyle_AnsiStyle_Render_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestStyle_AnsiStyle_Render_Ugly(t *core.T) {
	var s *AnsiStyle
	got := s.Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestStyle_Truncate_Good(t *core.T) {
	got := Truncate("abcdef", 4)

	core.AssertEqual(t, "a...", got)
	core.AssertLen(t, got, 4)
}

func TestStyle_Truncate_Bad(t *core.T) {
	got := Truncate("abcdef", 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestStyle_Truncate_Ugly(t *core.T) {
	got := Truncate("abcdef", 2)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestStyle_Glyph_Good(t *core.T) {
	got := Glyph(":check:")

	core.AssertEqual(t, "✓", got)
	core.AssertNotEqual(t, ":check:", got)
}

func TestStyle_Glyph_Bad(t *core.T) {
	got := Glyph(":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestStyle_Glyph_Ugly(t *core.T) {
	got := Glyph("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}
