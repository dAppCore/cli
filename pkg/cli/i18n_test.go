package cli

import "testing"

func TestT_Good(t *testing.T) {
	// T should return a non-empty string for any key
	// (falls back to the key itself when no translation is found).
	result := T("some.key")
	if result == "" {
		t.Error("T: returned empty string for unknown key")
	}
}

func TestT_Bad(t *testing.T) {
	// T with args map should not panic.
	result := T("cmd.doctor.issues", map[string]any{"Count": 0})
	if result == "" {
		t.Error("T with args: returned empty string")
	}
}

func TestT_Ugly(t *testing.T) {
	// T with empty key should not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("T(\"\") panicked: %v", r)
		}
	}()
	_ = T("")
}
