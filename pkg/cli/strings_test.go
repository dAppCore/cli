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

func TestStrings_Repeat_Good(t *testing.T) {
	if got := Repeat("-", 5); got != "-----" {
		t.Errorf("Repeat: got %q, want %q", got, "-----")
	}
	if got := Repeat("ab", 3); got != "ababab" {
		t.Errorf("Repeat: got %q, want %q", got, "ababab")
	}
}

func TestStrings_Repeat_Bad(t *testing.T) {
	// Zero count returns empty.
	if got := Repeat("x", 0); got != "" {
		t.Errorf("Repeat zero: got %q, want \"\"", got)
	}
	// Empty string repeated returns empty.
	if got := Repeat("", 5); got != "" {
		t.Errorf("Repeat empty: got %q, want \"\"", got)
	}
}

func TestStrings_Repeat_Ugly(t *testing.T) {
	// Negative count returns empty (instead of panicking like strings.Repeat).
	if got := Repeat("x", -3); got != "" {
		t.Errorf("Repeat negative: got %q, want \"\"", got)
	}
}

func TestStrings_LastIndex_Good(t *testing.T) {
	if got := LastIndex("hello world", "o"); got != 7 {
		t.Errorf("LastIndex: got %d, want 7", got)
	}
}

func TestStrings_LastIndex_Bad(t *testing.T) {
	if got := LastIndex("hello", "z"); got != -1 {
		t.Errorf("LastIndex absent: got %d, want -1", got)
	}
}

func TestStrings_LastIndex_Ugly(t *testing.T) {
	if got := LastIndex("", ""); got != 0 {
		t.Errorf("LastIndex empty: got %d, want 0", got)
	}
}

func TestStrings_Atoi_Good(t *testing.T) {
	n, err := Atoi("42")
	if err != nil {
		t.Fatalf("Atoi: unexpected error %v", err)
	}
	if n != 42 {
		t.Errorf("Atoi: got %d, want 42", n)
	}
}

func TestStrings_Atoi_Bad(t *testing.T) {
	if _, err := Atoi("not a number"); err == nil {
		t.Error("Atoi: expected error for non-numeric input")
	}
}

func TestStrings_Atoi_Ugly(t *testing.T) {
	if _, err := Atoi(""); err == nil {
		t.Error("Atoi empty: expected error")
	}
}

func TestStrings_ParseHexByte_Good(t *testing.T) {
	v, err := ParseHexByte("ff")
	if err != nil {
		t.Fatalf("ParseHexByte: unexpected error %v", err)
	}
	if v != 255 {
		t.Errorf("ParseHexByte: got %d, want 255", v)
	}
}

func TestStrings_ParseHexByte_Bad(t *testing.T) {
	if _, err := ParseHexByte("zz"); err == nil {
		t.Error("ParseHexByte: expected error for invalid hex")
	}
}

func TestStrings_ParseHexByte_Ugly(t *testing.T) {
	// Out-of-range hex (3 chars) should fail with overflow.
	if _, err := ParseHexByte("fff"); err == nil {
		t.Error("ParseHexByte overflow: expected error for out-of-range value")
	}
}
