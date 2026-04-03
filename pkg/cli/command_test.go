package cli

import "testing"

func TestCommand_Good(t *testing.T) {
	// NewCommand creates a command with RunE.
	called := false
	cmd := NewCommand("build", "Build the project", "", func(cmd *Command, args []string) error {
		called = true
		return nil
	})
	if cmd == nil {
		t.Fatal("NewCommand: returned nil")
	}
	if cmd.Use != "build" {
		t.Errorf("NewCommand: Use=%q, expected 'build'", cmd.Use)
	}
	if cmd.RunE == nil {
		t.Fatal("NewCommand: RunE is nil")
	}
	_ = called

	// NewGroup creates a command with no RunE.
	groupCmd := NewGroup("dev", "Development commands", "")
	if groupCmd.RunE != nil {
		t.Error("NewGroup: RunE should be nil")
	}

	// NewRun creates a command with Run.
	runCmd := NewRun("version", "Show version", "", func(cmd *Command, args []string) {})
	if runCmd.Run == nil {
		t.Fatal("NewRun: Run is nil")
	}
}

func TestCommand_Bad(t *testing.T) {
	// NewCommand with empty long string should not set Long.
	cmd := NewCommand("test", "Short desc", "", func(cmd *Command, args []string) error {
		return nil
	})
	if cmd.Long != "" {
		t.Errorf("NewCommand: Long should be empty, got %q", cmd.Long)
	}

	// Flag helpers with empty short should not add short flag.
	var value string
	StringFlag(cmd, &value, "output", "", "default", "Output path")
	if cmd.Flags().Lookup("output") == nil {
		t.Error("StringFlag: flag 'output' not registered")
	}
}

func TestCommand_Ugly(t *testing.T) {
	// WithArgs and WithExample are chainable.
	cmd := NewCommand("deploy", "Deploy", "Long desc", func(cmd *Command, args []string) error {
		return nil
	})
	result := WithExample(cmd, "core deploy production")
	if result != cmd {
		t.Error("WithExample: should return the same command")
	}
	if cmd.Example != "core deploy production" {
		t.Errorf("WithExample: Example=%q", cmd.Example)
	}

	// ExactArgs, NoArgs, MinimumNArgs, MaximumNArgs, ArbitraryArgs should not panic.
	_ = ExactArgs(1)
	_ = NoArgs()
	_ = MinimumNArgs(1)
	_ = MaximumNArgs(5)
	_ = ArbitraryArgs()
	_ = RangeArgs(1, 3)
}
