package cli

import (
	"strings"
	"testing"
)

func TestCheckBuilder(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII() // Deterministic output

	// Pass
	c := Check("foo").Pass()
	got := c.String()
	if got == "" {
		t.Error("Empty output for Pass")
	}

	// Fail
	c = Check("foo").Fail()
	got = c.String()
	if got == "" {
		t.Error("Empty output for Fail")
	}

	// Skip
	c = Check("foo").Skip()
	got = c.String()
	if got == "" {
		t.Error("Empty output for Skip")
	}

	// Warn
	c = Check("foo").Warn()
	got = c.String()
	if got == "" {
		t.Error("Empty output for Warn")
	}

	// Duration
	c = Check("foo").Pass().Duration("1s")
	got = c.String()
	if got == "" {
		t.Error("Empty output for Duration")
	}

	// Message
	c = Check("foo").Message("status")
	got = c.String()
	if got == "" {
		t.Error("Empty output for Message")
	}

	// Glyph shortcodes
	c = Check(":check: foo").Warn().Message(":warn:")
	got = c.String()
	if got == "" {
		t.Error("Empty output for glyph shortcode rendering")
	}
	if !strings.Contains(got, "[OK] foo") {
		t.Error("Expected shortcode-rendered name")
	}
	if strings.Count(got, "[WARN]") < 2 {
		t.Error("Expected shortcode-rendered warning icon and message")
	}
}
