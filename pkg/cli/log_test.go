package cli

import "testing"

func TestLog_Good(t *testing.T) {
	// All log functions should not panic when called without a configured logger.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogInfo panicked: %v", r)
		}
	}()
	LogInfo("test info message", "key", "value")
}

func TestLog_Bad(t *testing.T) {
	// LogError should not panic with an empty message.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogError panicked: %v", r)
		}
	}()
	LogError("")
}

func TestLog_Ugly(t *testing.T) {
	// All log levels should not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("log function panicked: %v", r)
		}
	}()
	LogDebug("debug", "k", "v")
	LogInfo("info", "k", "v")
	LogWarn("warn", "k", "v")
	LogError("error", "k", "v")

	// Level constants should be accessible.
	_ = LogLevelQuiet
	_ = LogLevelError
	_ = LogLevelWarn
	_ = LogLevelInfo
	_ = LogLevelDebug
}
