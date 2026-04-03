package cli

import "testing"

func TestRuntime_Good(t *testing.T) {
	// Init with valid options should succeed.
	err := Init(Options{
		AppName: "test-cli",
		Version: "0.0.1",
	})
	if err != nil {
		t.Fatalf("Init: unexpected error: %v", err)
	}
	defer Shutdown()

	// Core() returns non-nil after Init.
	coreInstance := Core()
	if coreInstance == nil {
		t.Error("Core(): returned nil after Init")
	}

	// RootCmd() returns non-nil after Init.
	rootCommand := RootCmd()
	if rootCommand == nil {
		t.Error("RootCmd(): returned nil after Init")
	}

	// Context() returns non-nil after Init.
	ctx := Context()
	if ctx == nil {
		t.Error("Context(): returned nil after Init")
	}
}

func TestRuntime_Bad(t *testing.T) {
	// Shutdown when not initialised should not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Shutdown() panicked when not initialised: %v", r)
		}
	}()
	// Reset singleton so this test can run standalone.
	// We use a fresh Shutdown here — it should be a no-op.
	Shutdown()
}

func TestRuntime_Ugly(t *testing.T) {
	// Once is idempotent: calling Init twice should succeed.
	err := Init(Options{AppName: "test-ugly"})
	if err != nil {
		t.Fatalf("Init (second call): unexpected error: %v", err)
	}
	defer Shutdown()
}
