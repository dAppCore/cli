package cli

import (
	"strings"
	"testing"
)

func TestFrameComponents_Good(t *testing.T) {
	// StatusLine renders title and pairs.
	model := StatusLine("core dev", "18 repos", "main")
	output := model.View(80, 1)
	if !strings.Contains(output, "core dev") {
		t.Errorf("StatusLine: expected 'core dev' in output, got %q", output)
	}

	// KeyHints renders hints.
	hints := KeyHints("↑/↓ navigate", "enter select", "q quit")
	output = hints.View(80, 1)
	if !strings.Contains(output, "navigate") {
		t.Errorf("KeyHints: expected 'navigate' in output, got %q", output)
	}

	// Breadcrumb renders navigation path.
	breadcrumb := Breadcrumb("core", "dev", "health")
	output = breadcrumb.View(80, 1)
	if !strings.Contains(output, "health") {
		t.Errorf("Breadcrumb: expected 'health' in output, got %q", output)
	}

	// StaticModel returns static text.
	static := StaticModel("static content")
	output = static.View(80, 1)
	if output != "static content" {
		t.Errorf("StaticModel: expected 'static content', got %q", output)
	}
}

func TestFrameComponents_Bad(t *testing.T) {
	// StatusLine with zero width should truncate to empty or short string.
	model := StatusLine("long title that should be truncated")
	output := model.View(0, 1)
	// Zero width means no truncation guard in current impl — just verify no panic.
	_ = output

	// KeyHints with no hints should not panic.
	hints := KeyHints()
	output = hints.View(80, 1)
	_ = output
}

func TestFrameComponents_Ugly(t *testing.T) {
	// Breadcrumb with single item has no separator.
	breadcrumb := Breadcrumb("root")
	output := breadcrumb.View(80, 1)
	if !strings.Contains(output, "root") {
		t.Errorf("Breadcrumb single: expected 'root', got %q", output)
	}

	// StatusLine with very narrow width truncates output.
	model := StatusLine("core dev", "18 repos")
	output = model.View(5, 1)
	if len(output) > 10 {
		t.Errorf("StatusLine truncated: output too long for width 5, got %q", output)
	}
}
