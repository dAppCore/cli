package cli

import "testing"

// expectedGlyphs lists shortcodes that must exist in every theme map.
// All three maps (unicode, emoji, ASCII) must provide a symbol for each.
var expectedGlyphs = []string{
	":check:", ":cross:", ":warn:", ":info:",
	":question:", ":skip:", ":dot:", ":circle:",
	":arrow_right:", ":arrow_left:", ":arrow_up:", ":arrow_down:",
	":pointer:", ":bullet:", ":dash:", ":pipe:",
	":corner:", ":tee:", ":pending:", ":spinner:",
}

func TestGlyphMaps_Good(t *testing.T) {
	// All three glyph maps must cover the same shortcode set.
	maps := map[string]map[string]string{
		"unicode": glyphMapUnicode,
		"emoji":   glyphMapEmoji,
		"ascii":   glyphMapASCII,
	}

	for name, m := range maps {
		if len(m) == 0 {
			t.Errorf("glyph map %q is empty", name)
			continue
		}
		for _, code := range expectedGlyphs {
			sym, ok := m[code]
			if !ok {
				t.Errorf("glyph map %q missing shortcode %q", name, code)
				continue
			}
			if sym == "" {
				t.Errorf("glyph map %q has empty symbol for %q", name, code)
			}
		}
	}
}

func TestGlyphMaps_Bad(t *testing.T) {
	// Unknown shortcode must NOT be present in any map.
	for name, m := range map[string]map[string]string{
		"unicode": glyphMapUnicode,
		"emoji":   glyphMapEmoji,
		"ascii":   glyphMapASCII,
	} {
		if _, ok := m[":nonexistent:"]; ok {
			t.Errorf("glyph map %q unexpectedly contains :nonexistent:", name)
		}
		if _, ok := m[""]; ok {
			t.Errorf("glyph map %q unexpectedly contains empty shortcode", name)
		}
	}
}

func TestGlyphMaps_Ugly(t *testing.T) {
	// Symbols for the same shortcode should differ across themes where expected.
	// ASCII :check: should not equal Unicode :check: — they are deliberately different.
	if glyphMapASCII[":check:"] == glyphMapUnicode[":check:"] {
		t.Error("ASCII and Unicode :check: glyphs should differ")
	}
	if glyphMapASCII[":cross:"] == glyphMapUnicode[":cross:"] {
		t.Error("ASCII and Unicode :cross: glyphs should differ")
	}

	// Every map must be the same size — themes are parallel alternatives,
	// not overlapping subsets.
	if len(glyphMapUnicode) != len(glyphMapEmoji) {
		t.Errorf("unicode map has %d entries, emoji map has %d — must be equal",
			len(glyphMapUnicode), len(glyphMapEmoji))
	}
	if len(glyphMapUnicode) != len(glyphMapASCII) {
		t.Errorf("unicode map has %d entries, ASCII map has %d — must be equal",
			len(glyphMapUnicode), len(glyphMapASCII))
	}
}
