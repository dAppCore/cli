package cli

import (
	"bytes"
	"strings"
	"testing"

	"forge.lthn.ai/core/go-log"
)

func TestLogSecurity_Good(t *testing.T) {
	var buf bytes.Buffer
	original := log.Default()
	t.Cleanup(func() {
		log.SetDefault(original)
	})

	logger := log.New(log.Options{Level: log.LevelDebug, Output: &buf})
	log.SetDefault(logger)

	LogSecurity("login attempt", "user", "admin")

	out := buf.String()
	if !strings.Contains(out, "login attempt") {
		t.Fatalf("expected security log message, got %q", out)
	}
	if !strings.Contains(out, "user") {
		t.Fatalf("expected structured key/value output, got %q", out)
	}
}
