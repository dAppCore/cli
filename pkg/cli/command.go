package cli

import (
	"time"

	"github.com/spf13/cobra"
)

// ─────────────────────────────────────────────────────────────────────────────
// Cobra Re-exports
// ─────────────────────────────────────────────────────────────────────────────

// PositionalArgs is the cobra positional args type.
type PositionalArgs = cobra.PositionalArgs

// ─────────────────────────────────────────────────────────────────────────────
// Command Type Re-export
// ─────────────────────────────────────────────────────────────────────────────

// Command is the cobra command type.
// Re-exported for convenience so packages don't need to import cobra directly.
type Command = cobra.Command

// ─────────────────────────────────────────────────────────────────────────────
// Command Builders
// ─────────────────────────────────────────────────────────────────────────────

// NewCommand creates a new command with a RunE handler.
// This is the standard way to create commands that may return errors.
//
//	cmd := cli.NewCommand("build", "Build the project", "", func(cmd *cli.Command, args []string) error {
//	    // Build logic
//	    return nil
//	})
func NewCommand(use, short, long string, run func(cmd *Command, args []string) error) *Command {
	cmd := &Command{
		Use:   use,
		Short: short,
		RunE:  run,
	}
	if long != "" {
		cmd.Long = long
	}
	return cmd
}

// NewGroup creates a new command group (no RunE).
// Use this for parent commands that only contain subcommands.
//
//	devCmd := cli.NewGroup("dev", "Development commands", "")
//	devCmd.AddCommand(buildCmd, testCmd)
func NewGroup(use, short, long string) *Command {
	cmd := &Command{
		Use:   use,
		Short: short,
	}
	if long != "" {
		cmd.Long = long
	}
	return cmd
}

// NewRun creates a new command with a simple Run handler (no error return).
// Use when the command cannot fail.
//
//	cmd := cli.NewRun("version", "Show version", "", func(cmd *cli.Command, args []string) {
//	    cli.Println("v1.0.0")
//	})
func NewRun(use, short, long string, run func(cmd *Command, args []string)) *Command {
	cmd := &Command{
		Use:   use,
		Short: short,
		Run:   run,
	}
	if long != "" {
		cmd.Long = long
	}
	return cmd
}

// ─────────────────────────────────────────────────────────────────────────────
// Flag Helpers
// ─────────────────────────────────────────────────────────────────────────────

// StringFlag adds a string flag to a command.
// The value will be stored in the provided pointer.
//
//	var output string
//	cli.StringFlag(cmd, &output, "output", "o", "", "Output file path")
func StringFlag(cmd *Command, ptr *string, name, short, def, usage string) {
	if short != "" {
		cmd.Flags().StringVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().StringVar(ptr, name, def, usage)
	}
}

// BoolFlag adds a boolean flag to a command.
// The value will be stored in the provided pointer.
//
//	var verbose bool
//	cli.BoolFlag(cmd, &verbose, "verbose", "v", false, "Enable verbose output")
func BoolFlag(cmd *Command, ptr *bool, name, short string, def bool, usage string) {
	if short != "" {
		cmd.Flags().BoolVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().BoolVar(ptr, name, def, usage)
	}
}

// IntFlag adds an integer flag to a command.
// The value will be stored in the provided pointer.
//
//	var count int
//	cli.IntFlag(cmd, &count, "count", "n", 10, "Number of items")
func IntFlag(cmd *Command, ptr *int, name, short string, def int, usage string) {
	if short != "" {
		cmd.Flags().IntVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().IntVar(ptr, name, def, usage)
	}
}

// Float64Flag adds a float64 flag to a command.
// The value will be stored in the provided pointer.
//
//	var threshold float64
//	cli.Float64Flag(cmd, &threshold, "threshold", "t", 0.0, "Score threshold")
func Float64Flag(cmd *Command, ptr *float64, name, short string, def float64, usage string) {
	if short != "" {
		cmd.Flags().Float64VarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().Float64Var(ptr, name, def, usage)
	}
}

// Int64Flag adds an int64 flag to a command.
// The value will be stored in the provided pointer.
//
//	var seed int64
//	cli.Int64Flag(cmd, &seed, "seed", "s", 0, "Random seed")
func Int64Flag(cmd *Command, ptr *int64, name, short string, def int64, usage string) {
	if short != "" {
		cmd.Flags().Int64VarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().Int64Var(ptr, name, def, usage)
	}
}

// DurationFlag adds a time.Duration flag to a command.
// The value will be stored in the provided pointer.
//
//	var timeout time.Duration
//	cli.DurationFlag(cmd, &timeout, "timeout", "t", 30*time.Second, "Request timeout")
func DurationFlag(cmd *Command, ptr *time.Duration, name, short string, def time.Duration, usage string) {
	if short != "" {
		cmd.Flags().DurationVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().DurationVar(ptr, name, def, usage)
	}
}

// StringSliceFlag adds a string slice flag to a command.
// The value will be stored in the provided pointer.
//
//	var tags []string
//	cli.StringSliceFlag(cmd, &tags, "tag", "t", nil, "Tags to apply")
func StringSliceFlag(cmd *Command, ptr *[]string, name, short string, def []string, usage string) {
	if short != "" {
		cmd.Flags().StringSliceVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().StringSliceVar(ptr, name, def, usage)
	}
}

// StringArrayFlag adds a string array flag to a command.
// The value will be stored in the provided pointer.
//
//	var tags []string
//	cli.StringArrayFlag(cmd, &tags, "tag", "t", nil, "Tags to apply")
func StringArrayFlag(cmd *Command, ptr *[]string, name, short string, def []string, usage string) {
	if short != "" {
		cmd.Flags().StringArrayVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().StringArrayVar(ptr, name, def, usage)
	}
}

// StringToStringFlag adds a string-to-string map flag to a command.
// The value will be stored in the provided pointer.
//
//	var labels map[string]string
//	cli.StringToStringFlag(cmd, &labels, "label", "l", nil, "Labels to apply")
func StringToStringFlag(cmd *Command, ptr *map[string]string, name, short string, def map[string]string, usage string) {
	if short != "" {
		cmd.Flags().StringToStringVarP(ptr, name, short, def, usage)
	} else {
		cmd.Flags().StringToStringVar(ptr, name, def, usage)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Persistent Flag Helpers
// ─────────────────────────────────────────────────────────────────────────────

// PersistentStringFlag adds a persistent string flag (inherited by subcommands).
func PersistentStringFlag(cmd *Command, ptr *string, name, short, def, usage string) {
	if short != "" {
		cmd.PersistentFlags().StringVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().StringVar(ptr, name, def, usage)
	}
}

// PersistentBoolFlag adds a persistent boolean flag (inherited by subcommands).
func PersistentBoolFlag(cmd *Command, ptr *bool, name, short string, def bool, usage string) {
	if short != "" {
		cmd.PersistentFlags().BoolVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().BoolVar(ptr, name, def, usage)
	}
}

// PersistentIntFlag adds a persistent integer flag (inherited by subcommands).
func PersistentIntFlag(cmd *Command, ptr *int, name, short string, def int, usage string) {
	if short != "" {
		cmd.PersistentFlags().IntVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().IntVar(ptr, name, def, usage)
	}
}

// PersistentInt64Flag adds a persistent int64 flag (inherited by subcommands).
func PersistentInt64Flag(cmd *Command, ptr *int64, name, short string, def int64, usage string) {
	if short != "" {
		cmd.PersistentFlags().Int64VarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().Int64Var(ptr, name, def, usage)
	}
}

// PersistentFloat64Flag adds a persistent float64 flag (inherited by subcommands).
func PersistentFloat64Flag(cmd *Command, ptr *float64, name, short string, def float64, usage string) {
	if short != "" {
		cmd.PersistentFlags().Float64VarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().Float64Var(ptr, name, def, usage)
	}
}

// PersistentDurationFlag adds a persistent time.Duration flag (inherited by subcommands).
func PersistentDurationFlag(cmd *Command, ptr *time.Duration, name, short string, def time.Duration, usage string) {
	if short != "" {
		cmd.PersistentFlags().DurationVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().DurationVar(ptr, name, def, usage)
	}
}

// PersistentStringSliceFlag adds a persistent string slice flag (inherited by subcommands).
func PersistentStringSliceFlag(cmd *Command, ptr *[]string, name, short string, def []string, usage string) {
	if short != "" {
		cmd.PersistentFlags().StringSliceVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().StringSliceVar(ptr, name, def, usage)
	}
}

// PersistentStringArrayFlag adds a persistent string array flag (inherited by subcommands).
func PersistentStringArrayFlag(cmd *Command, ptr *[]string, name, short string, def []string, usage string) {
	if short != "" {
		cmd.PersistentFlags().StringArrayVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().StringArrayVar(ptr, name, def, usage)
	}
}

// PersistentStringToStringFlag adds a persistent string-to-string map flag (inherited by subcommands).
func PersistentStringToStringFlag(cmd *Command, ptr *map[string]string, name, short string, def map[string]string, usage string) {
	if short != "" {
		cmd.PersistentFlags().StringToStringVarP(ptr, name, short, def, usage)
	} else {
		cmd.PersistentFlags().StringToStringVar(ptr, name, def, usage)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Command Configuration
// ─────────────────────────────────────────────────────────────────────────────

// WithArgs sets the Args validation function for a command.
// Returns the command for chaining.
//
//	cmd := cli.NewCommand("build", "Build", "", run).WithArgs(cobra.ExactArgs(1))
func WithArgs(cmd *Command, args cobra.PositionalArgs) *Command {
	cmd.Args = args
	return cmd
}

// WithExample sets the Example field for a command.
// Returns the command for chaining.
func WithExample(cmd *Command, example string) *Command {
	cmd.Example = example
	return cmd
}

// ExactArgs returns a PositionalArgs that accepts exactly N arguments.
func ExactArgs(n int) cobra.PositionalArgs {
	return cobra.ExactArgs(n)
}

// MinimumNArgs returns a PositionalArgs that accepts minimum N arguments.
func MinimumNArgs(n int) cobra.PositionalArgs {
	return cobra.MinimumNArgs(n)
}

// MaximumNArgs returns a PositionalArgs that accepts maximum N arguments.
func MaximumNArgs(n int) cobra.PositionalArgs {
	return cobra.MaximumNArgs(n)
}

// RangeArgs returns a PositionalArgs that accepts between min and max arguments.
func RangeArgs(min int, max int) cobra.PositionalArgs {
	return cobra.RangeArgs(min, max)
}

// NoArgs returns a PositionalArgs that accepts no arguments.
func NoArgs() cobra.PositionalArgs {
	return cobra.NoArgs
}

// ArbitraryArgs returns a PositionalArgs that accepts any arguments.
func ArbitraryArgs() cobra.PositionalArgs {
	return cobra.ArbitraryArgs
}
