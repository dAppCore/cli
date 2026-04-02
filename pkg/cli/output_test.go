package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	oldOut := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	f()

	_ = w.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestSemanticOutput(t *testing.T) {
	UseASCII()

	// Test Success
	out := captureOutput(func() {
		Success("done")
	})
	if out == "" {
		t.Error("Success output empty")
	}

	// Test Error
	out = captureOutput(func() {
		Error("fail")
	})
	if out == "" {
		t.Error("Error output empty")
	}

	// Test Warn
	out = captureOutput(func() {
		Warn("warn")
	})
	if out == "" {
		t.Error("Warn output empty")
	}

	// Test Info
	out = captureOutput(func() {
		Info("info")
	})
	if out == "" {
		t.Error("Info output empty")
	}

	// Test Task
	out = captureOutput(func() {
		Task("task", "msg")
	})
	if out == "" {
		t.Error("Task output empty")
	}

	// Test Section
	out = captureOutput(func() {
		Section("section")
	})
	if out == "" {
		t.Error("Section output empty")
	}

	// Test Hint
	out = captureOutput(func() {
		Hint("hint", "msg")
	})
	if out == "" {
		t.Error("Hint output empty")
	}

	// Test Result
	out = captureOutput(func() {
		Result(true, "pass")
	})
	if out == "" {
		t.Error("Result(true) output empty")
	}

	out = captureOutput(func() {
		Result(false, "fail")
	})
	if out == "" {
		t.Error("Result(false) output empty")
	}
}

func TestSemanticOutput_GlyphShortcodes(t *testing.T) {
	UseASCII()

	out := captureOutput(func() {
		Success("done :check:")
		Task(":cross:", "running :warn:")
		Section(":check: audit")
		Hint(":info:", "apply :check:")
		Label("status", "ready :warn:")
	})

	for _, want := range []string{"[OK]", "[FAIL]", "[WARN]"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got %q", want, out)
		}
	}
}
