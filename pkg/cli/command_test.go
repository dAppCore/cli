package cli

import (
	"testing"

	"dappco.re/go"
)

func TestCommand_Good(t *testing.T) {
	// RegisterCommand registers a command on Core.
	c := core.New()
	RegisterCommand(c, "build", core.Command{
		Description: "Build the project",
		Action: func(_ core.Options) core.Result {
			return core.Ok(nil)
		},
	})

	r := c.Command("build")
	if !r.OK {
		t.Fatal("RegisterCommand: command not found after registration")
	}

	cmd := r.Value.(*core.Command)
	if cmd.Name != "build" {
		t.Errorf("RegisterCommand: Name=%q, expected 'build'", cmd.Name)
	}
}

func TestCommand_Bad(t *testing.T) {
	// RequireArgs with no args should return error message.
	opts := core.NewOptions()
	msg := RequireArgs(opts, 1)
	if msg == "" {
		t.Error("RequireArgs: should return error message when no args present")
	}

	// RequireArgs with args should return empty.
	opts.Set("_arg", "value")
	msg = RequireArgs(opts, 1)
	if msg != "" {
		t.Errorf("RequireArgs: should return empty string when args present, got %q", msg)
	}
}

func TestCommand_Ugly(t *testing.T) {
	// RequireExactArgs with 0 and no arg should pass.
	opts := core.NewOptions()
	msg := RequireExactArgs(opts, 0)
	if msg != "" {
		t.Errorf("RequireExactArgs(0): expected empty, got %q", msg)
	}

	// RequireExactArgs with 0 but arg present should fail.
	opts.Set("_arg", "unexpected")
	msg = RequireExactArgs(opts, 0)
	if msg == "" {
		t.Error("RequireExactArgs(0): should fail when args present")
	}

	// Path-based nested commands work.
	c := core.New()
	RegisterCommand(c, "deploy/to/homelab", core.Command{
		Description: "Deploy to homelab",
		Action: func(_ core.Options) core.Result {
			return core.Ok(nil)
		},
	})
	r := c.Command("deploy/to/homelab")
	if !r.OK {
		t.Error("RegisterCommand: nested path command not found")
	}
}
