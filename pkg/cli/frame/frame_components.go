package frame

import (
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// Built-in Region Components
// ─────────────────────────────────────────────────────────────────────────────

// statusLineModel renders a "title  key:value  key:value" bar.
type statusLineModel struct {
	title string
	pairs []string
}

// StatusLine creates a header/footer bar with a title and key:value pairs.
//
//	frame.Header(frame.StatusLine("core dev", "18 repos", "main"))
func StatusLine(title string, pairs ...string) Model {
	return &statusLineModel{title: title, pairs: pairs}
}

func (s *statusLineModel) View(width, _ int) string {
	parts := []string{BoldStyle.Render(compileGlyphs(s.title))}
	for _, p := range s.pairs {
		parts = append(parts, DimStyle.Render(compileGlyphs(p)))
	}
	line := strings.Join(parts, "  ")
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// keyHintsModel renders keyboard shortcut hints.
type keyHintsModel struct {
	hints []string
}

// KeyHints creates a footer showing keyboard shortcuts.
//
//	frame.Footer(frame.KeyHints("↑/↓ navigate", "enter select", "q quit"))
func KeyHints(hints ...string) Model {
	return &keyHintsModel{hints: hints}
}

func (k *keyHintsModel) View(width, _ int) string {
	parts := make([]string, len(k.hints))
	for i, h := range k.hints {
		parts[i] = DimStyle.Render(compileGlyphs(h))
	}
	line := strings.Join(parts, "  ")
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// breadcrumbModel renders a navigation path.
type breadcrumbModel struct {
	parts []string
}

// Breadcrumb creates a navigation breadcrumb bar.
//
//	frame.Header(frame.Breadcrumb("core", "dev", "health"))
func Breadcrumb(parts ...string) Model {
	return &breadcrumbModel{parts: parts}
}

func (b *breadcrumbModel) View(width, _ int) string {
	styled := make([]string, len(b.parts))
	for i, p := range b.parts {
		part := compileGlyphs(p)
		if i == len(b.parts)-1 {
			styled[i] = BoldStyle.Render(part)
		} else {
			styled[i] = DimStyle.Render(part)
		}
	}
	line := strings.Join(styled, DimStyle.Render(" > "))
	if width > 0 {
		line = Truncate(line, width)
	}
	return line
}

// staticModel wraps a plain string as a Model.
type staticModel struct {
	text string
}

// StaticModel wraps a static string as a Model, for use in Frame regions.
func StaticModel(text string) Model {
	return &staticModel{text: text}
}

func (s *staticModel) View(_, _ int) string {
	return compileGlyphs(s.text)
}
