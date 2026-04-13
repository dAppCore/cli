package cli

import (
	"dappco.re/go/core"
)

// ─────────────────────────────────────────────────────────────────────────────
// Command Type
// ─────────────────────────────────────────────────────────────────────────────

// Command is the core command type.
// Re-exported for convenience so packages don't need to import core directly.
type Command = core.Command

// CommandAction is the function signature for command handlers.
type CommandAction = core.CommandAction

// ─────────────────────────────────────────────────────────────────────────────
// Command Registration Helpers
// ─────────────────────────────────────────────────────────────────────────────

// RegisterCommand registers a command on the Core instance using path-based routing.
// This is the primary way to register commands in the core/go Cli+Command pattern.
//
//	cli.RegisterCommand(c, "config/list", core.Command{
//	    Description: "List all configuration values",
//	    Action: func(opts core.Options) core.Result {
//	        cli.Println("listing...")
//	        return core.Result{OK: true}
//	    },
//	})
func RegisterCommand(c *core.Core, path string, cmd core.Command) {
	c.Command(path, cmd)
}

// ─────────────────────────────────────────────────────────────────────────────
// Arg Helpers
// ─────────────────────────────────────────────────────────────────────────────

// RequireArgs validates that at least n positional arguments are present in opts.
// Returns an error string if insufficient args, empty string if OK.
// Use inside a CommandAction to validate argument count.
//
//	func myAction(opts core.Options) core.Result {
//	    if msg := cli.RequireArgs(opts, 1); msg != "" {
//	        return core.Result{Value: cli.Err(msg), OK: false}
//	    }
//	    key := opts.String("_arg")
//	    // ...
//	}
func RequireArgs(opts core.Options, n int) string {
	arg := opts.String("_arg")
	if n > 0 && arg == "" {
		return Sprintf("requires at least %d argument(s)", n)
	}
	return ""
}

// RequireExactArgs validates that exactly n positional arguments are present.
// Core/go stores the first positional arg in "_arg". For commands needing
// multiple positional args, the remaining args are available from the raw
// args slice passed to Cli.Run().
//
//	func myAction(opts core.Options) core.Result {
//	    if msg := cli.RequireExactArgs(opts, 1); msg != "" {
//	        return core.Result{Value: cli.Err(msg), OK: false}
//	    }
//	}
func RequireExactArgs(opts core.Options, n int) string {
	if n == 0 {
		arg := opts.String("_arg")
		if arg != "" {
			return "accepts no arguments"
		}
		return ""
	}
	return RequireArgs(opts, n)
}
