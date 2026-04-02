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
	reader, writer, _ := os.Pipe()
	os.Stdout = writer
	os.Stderr = writer

	f()

	_ = writer.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, reader)
	return buf.String()
}

func TestSemanticOutput_Good(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII()
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	cases := []struct {
		name string
		fn   func()
	}{
		{"Success", func() { Success("done") }},
		{"Info", func() { Info("info") }},
		{"Task", func() { Task("task", "msg") }},
		{"Section", func() { Section("section") }},
		{"Hint", func() { Hint("hint", "msg") }},
		{"Result_pass", func() { Result(true, "pass") }},
	}

	for _, testCase := range cases {
		output := captureOutput(testCase.fn)
		if output == "" {
			t.Errorf("%s: output was empty", testCase.name)
		}
	}
}

func TestSemanticOutput_Bad(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII()
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	// Error and Warn go to stderr — both captured here.
	errorOutput := captureOutput(func() { Error("fail") })
	if errorOutput == "" {
		t.Error("Error: output was empty")
	}

	warnOutput := captureOutput(func() { Warn("warn") })
	if warnOutput == "" {
		t.Error("Warn: output was empty")
	}

	failureOutput := captureOutput(func() { Result(false, "fail") })
	if failureOutput == "" {
		t.Error("Result(false): output was empty")
	}
}

func TestSemanticOutput_Ugly(t *testing.T) {
	restoreThemeAndColors(t)
	UseASCII()

	// Severity with various levels should not panic.
	levels := []string{"critical", "high", "medium", "low", "unknown", ""}
	for _, level := range levels {
		output := captureOutput(func() { Severity(level, "test message") })
		if output == "" {
			t.Errorf("Severity(%q): output was empty", level)
		}
	}

	// Section uppercases the name.
	output := captureOutput(func() { Section("audit") })
	if !strings.Contains(output, "AUDIT") {
		t.Errorf("Section: expected AUDIT in output, got %q", output)
	}
}
