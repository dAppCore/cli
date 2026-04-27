package cli

import (
	"bytes"   // Note: AX-6 — bytes.Buffer provides rune-by-rune parser state for glyph shortcode rendering; core has no NewBuffer wrapper in this release.
	"unicode" // Note: AX-6 — unicode is structural for rune classification in glyph rendering.
)

// GlyphTheme defines which symbols to use.
type GlyphTheme int

const (
	// ThemeUnicode uses standard Unicode symbols.
	ThemeUnicode GlyphTheme = iota
	// ThemeEmoji uses Emoji symbols.
	ThemeEmoji
	// ThemeASCII uses ASCII fallback symbols.
	ThemeASCII
)

var currentTheme = ThemeUnicode

// UseUnicode switches the glyph theme to Unicode.
func UseUnicode() {
	currentTheme = ThemeUnicode
	restoreColorIfASCII()
}

// UseEmoji switches the glyph theme to Emoji.
func UseEmoji() {
	currentTheme = ThemeEmoji
	restoreColorIfASCII()
}

// UseASCII switches the glyph theme to ASCII and disables colors.
func UseASCII() {
	currentTheme = ThemeASCII
	SetColorEnabled(false)
	colorEnabledMu.Lock()
	asciiDisabledColors = true
	colorEnabledMu.Unlock()
}

func glyphMap() map[string]string {
	switch currentTheme {
	case ThemeEmoji:
		return glyphMapEmoji
	case ThemeASCII:
		return glyphMapASCII
	default:
		return glyphMapUnicode
	}
}

// Glyph converts a shortcode (e.g. ":check:") to its symbol based on the current theme.
func Glyph(code string) string {
	if sym, ok := glyphMap()[code]; ok {
		return sym
	}
	return code
}

func compileGlyphs(x string) string {
	if x == "" {
		return ""
	}
	input := bytes.NewBufferString(x)
	output := bytes.NewBufferString("")

	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}
		if r == ':' {
			output.WriteString(replaceGlyph(input))
		} else {
			output.WriteRune(r)
		}
	}
	return output.String()
}

func replaceGlyph(input *bytes.Buffer) string {
	code := bytes.NewBufferString(":")
	for {
		r, _, err := input.ReadRune()
		if err != nil {
			return code.String()
		}
		if r == ':' && code.Len() == 1 {
			return code.String() + replaceGlyph(input)
		}
		code.WriteRune(r)
		if unicode.IsSpace(r) {
			return code.String()
		}
		if r == ':' {
			return Glyph(code.String())
		}
	}
}
