package cli

import (
	"strings"
	"testing"
)

func TestStrings_Good(t *testing.T) {
	// Sprintf formats correctly.
	result := Sprintf("Hello, %s! Count: %d", "world", 42)
	if result != "Hello, world! Count: 42" {
		t.Errorf("Sprintf: got %q", result)
	}

	// Sprint joins with spaces.
	result = Sprint("foo", "bar")
	if result == "" {
		t.Error("Sprint: got empty string")
	}

	// SuccessStr, ErrorStr, WarnStr, InfoStr, DimStr return non-empty strings.
	if SuccessStr("done") == "" {
		t.Error("SuccessStr: got empty string")
	}
	if ErrorStr("fail") == "" {
		t.Error("ErrorStr: got empty string")
	}
	if WarnStr("warn") == "" {
		t.Error("WarnStr: got empty string")
	}
	if InfoStr("info") == "" {
		t.Error("InfoStr: got empty string")
	}
	if DimStr("dim") == "" {
		t.Error("DimStr: got empty string")
	}
}

func TestStrings_Bad(t *testing.T) {
	// Sprintf with no args returns the format string unchanged.
	result := Sprintf("no args here")
	if result != "no args here" {
		t.Errorf("Sprintf no-args: got %q", result)
	}

	// Styled with nil style should not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Styled with nil style panicked: %v", r)
		}
	}()
	Styled(nil, "text")
}

func TestStrings_Ugly(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	// Without colour, styled strings contain the raw text.
	result := Styled(NewStyle().Bold(), "core")
	if !strings.Contains(result, "core") {
		t.Errorf("Styled: expected 'core' in result, got %q", result)
	}

	// Styledf with empty format.
	result = Styledf(DimStyle, "")
	_ = result // should not panic
}
