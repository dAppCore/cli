package cli

import "testing"

func TestGlyph_Good(t *testing.T) {
	restoreThemeAndColors(t)
	UseUnicode()
	if Glyph(":check:") != "✓" {
		t.Errorf("Expected ✓, got %s", Glyph(":check:"))
	}

	UseASCII()
	if Glyph(":check:") != "[OK]" {
		t.Errorf("Expected [OK], got %s", Glyph(":check:"))
	}
}

func TestGlyph_Bad(t *testing.T) {
	restoreThemeAndColors(t)
	// Unknown shortcode returns the shortcode unchanged.
	UseUnicode()
	got := Glyph(":unknown:")
	if got != ":unknown:" {
		t.Errorf("Unknown shortcode should return unchanged, got %q", got)
	}
}

func TestGlyph_Ugly(t *testing.T) {
	restoreThemeAndColors(t)
	// Empty shortcode should not panic.
	got := Glyph("")
	if got != "" {
		t.Errorf("Empty shortcode should return empty string, got %q", got)
	}
}

func TestCompileGlyphs_Good(t *testing.T) {
	restoreThemeAndColors(t)
	UseUnicode()
	got := compileGlyphs("Status: :check:")
	if got != "Status: ✓" {
		t.Errorf("Expected 'Status: ✓', got %q", got)
	}
}

func TestCompileGlyphs_Bad(t *testing.T) {
	restoreThemeAndColors(t)
	UseUnicode()
	// Text with no shortcodes should be returned as-is.
	got := compileGlyphs("no glyphs here")
	if got != "no glyphs here" {
		t.Errorf("Expected unchanged text, got %q", got)
	}
}

func TestCompileGlyphs_Ugly(t *testing.T) {
	restoreThemeAndColors(t)
	// Empty string should not panic.
	got := compileGlyphs("")
	if got != "" {
		t.Errorf("Empty string should return empty, got %q", got)
	}
}
