package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestErrors_Good(t *testing.T) {
	// Err creates a formatted error.
	err := Err("key not found: %s", "theme")
	if err == nil {
		t.Fatal("Err: expected non-nil error")
	}
	if !strings.Contains(err.Error(), "theme") {
		t.Errorf("Err: expected 'theme' in message, got %q", err.Error())
	}

	// Wrap prepends a message.
	base := errors.New("connection refused")
	wrapped := Wrap(base, "connect to database")
	if !strings.Contains(wrapped.Error(), "connect to database") {
		t.Errorf("Wrap: expected prefix in message, got %q", wrapped.Error())
	}
	if !Is(wrapped, base) {
		t.Error("Wrap: errors.Is should unwrap to original")
	}
}

func TestErrors_Bad(t *testing.T) {
	// Wrap with nil error returns nil.
	if Wrap(nil, "should be nil") != nil {
		t.Error("Wrap(nil): expected nil return")
	}

	// WrapVerb with nil error returns nil.
	if WrapVerb(nil, "load", "config") != nil {
		t.Error("WrapVerb(nil): expected nil return")
	}

	// WrapAction with nil error returns nil.
	if WrapAction(nil, "connect") != nil {
		t.Error("WrapAction(nil): expected nil return")
	}
}

func TestErrors_Ugly(t *testing.T) {
	// Join with multiple errors.
	err1 := Err("first error")
	err2 := Err("second error")
	joined := Join(err1, err2)
	if joined == nil {
		t.Fatal("Join: expected non-nil error")
	}
	if !Is(joined, err1) {
		t.Error("Join: errors.Is should find first error")
	}

	// Exit creates ExitError with correct code.
	exitErr := Exit(2, Err("exit with code 2"))
	if exitErr == nil {
		t.Fatal("Exit: expected non-nil error")
	}
	var exitErrorValue *ExitError
	if !As(exitErr, &exitErrorValue) {
		t.Fatal("Exit: expected *ExitError type")
	}
	if exitErrorValue.Code != 2 {
		t.Errorf("Exit: expected code 2, got %d", exitErrorValue.Code)
	}

	// Exit with nil returns nil.
	if Exit(1, nil) != nil {
		t.Error("Exit(nil): expected nil return")
	}
}
