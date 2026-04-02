package cli

import (
	"strings"
	"testing"
)

func TestCheckBuilder_Good(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII() // Deterministic output

	checkResult := Check("database").Pass()
	got := checkResult.String()
	if got == "" {
		t.Error("Pass: expected non-empty output")
	}
	if !strings.Contains(got, "database") {
		t.Errorf("Pass: expected name in output, got %q", got)
	}
}

func TestCheckBuilder_Bad(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII()

	checkResult := Check("lint").Fail()
	got := checkResult.String()
	if got == "" {
		t.Error("Fail: expected non-empty output")
	}

	checkResult = Check("build").Skip()
	got = checkResult.String()
	if got == "" {
		t.Error("Skip: expected non-empty output")
	}

	checkResult = Check("tests").Warn()
	got = checkResult.String()
	if got == "" {
		t.Error("Warn: expected non-empty output")
	}
}

func TestCheckBuilder_Ugly(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII()

	// Zero-value builder should not panic.
	checkResult := &CheckBuilder{}
	got := checkResult.String()
	if got == "" {
		t.Error("Ugly: empty builder should still produce output")
	}

	// Duration and Message chaining.
	checkResult = Check("audit").Pass().Duration("2.3s").Message("all clear")
	got = checkResult.String()
	if !strings.Contains(got, "2.3s") {
		t.Errorf("Ugly: expected duration in output, got %q", got)
	}
}
