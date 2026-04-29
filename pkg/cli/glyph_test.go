package cli

import (
	core "dappco.re/go"
)

func TestGlyph_UseUnicode_Good(t *core.T) {
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertEqual(t, "✓", Glyph(":check:"))
}

func TestGlyph_UseUnicode_Bad(t *core.T) {
	UseASCII()
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertTrue(t, ColorEnabled())
}

func TestGlyph_UseUnicode_Ugly(t *core.T) {
	UseUnicode()
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertEqual(t, "✓", Glyph(":check:"))
}

func TestGlyph_UseEmoji_Good(t *core.T) {
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertNotEqual(t, "✓", Glyph(":check:"))
}

func TestGlyph_UseEmoji_Bad(t *core.T) {
	UseASCII()
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertTrue(t, ColorEnabled())
}

func TestGlyph_UseEmoji_Ugly(t *core.T) {
	UseEmoji()
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertNotEmpty(t, Glyph(":check:"))
}

func TestGlyph_UseASCII_Good(t *core.T) {
	UseASCII()
	defer UseUnicode()

	core.AssertEqual(t, ThemeASCII, currentTheme)
	core.AssertEqual(t, "[OK]", Glyph(":check:"))
}

func TestGlyph_UseASCII_Bad(t *core.T) {
	UseASCII()
	defer UseUnicode()

	core.AssertFalse(t, ColorEnabled())
	core.AssertEqual(t, "[FAIL]", Glyph(":cross:"))
}

func TestGlyph_UseASCII_Ugly(t *core.T) {
	UseASCII()
	UseASCII()
	defer UseUnicode()

	core.AssertEqual(t, ThemeASCII, currentTheme)
	core.AssertFalse(t, ColorEnabled())
}

func TestGlyph_Glyph_Good(t *core.T) {
	UseUnicode()
	got := Glyph(":check:")

	core.AssertEqual(t, "✓", got)
	core.AssertNotEqual(t, ":check:", got)
}

func TestGlyph_Glyph_Bad(t *core.T) {
	UseUnicode()
	got := Glyph(":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestGlyph_Glyph_Ugly(t *core.T) {
	UseASCII()
	defer UseUnicode()
	got := Glyph(":cross:")

	core.AssertEqual(t, "[FAIL]", got)
	core.AssertFalse(t, ColorEnabled())
}
