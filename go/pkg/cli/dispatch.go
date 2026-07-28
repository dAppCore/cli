// dispatch.go is the CLI's own argument resolution: command lookup and option
// parsing.
//
// core/go's Cli.Run is deliberately surface-agnostic, and its parser recognises
// only the joined form `--key=value`. Given `--key value` it sets key to the
// boolean true and drops the value into "_arg", so every command reading
// opts.String("key") saw an empty string and silently fell back to its default.
// `core build installers --repo Acme/widget` generated scripts pointing at the
// default repository, reporting success the whole way.
//
// Argument parsing is a CLI-surface concern in exactly the way the catalog
// renderer in help.go is, so it lives here rather than in the core primitive.
// core/go keeps its own parser for callers that want it.
package cli

import (
	core "dappco.re/go"
)

// ArgsKey holds every positional argument. "_arg" carries the first for the
// commands that predate the slice.
const ArgsKey = "_args"

// resolveCommand returns the runnable command matching the longest leading run
// of args, along with the arguments that follow it.
//
// Longest-first matters: "build installers" must win over "build", or the
// subcommand's name would be parsed as a positional argument to its parent.
func resolveCommand(c *core.Core, args []string) (*core.Command, []string) {
	if c == nil {
		return nil, args
	}
	for i := len(args); i > 0; i-- {
		result := c.Command(core.JoinPath(args[:i]...))
		if !result.OK {
			continue
		}
		cmd, ok := result.Value.(*core.Command)
		if !ok || cmd.Action == nil {
			continue
		}
		return cmd, args[i:]
	}
	return nil, args
}

// takesValue reports whether a flag declared on cmd expects a value.
//
// A command may declare its flags in core.Command.Flags, in which case the
// declared default's type decides: a bool default is a switch and never
// consumes the following token, so `--verbose ./path` keeps ./path positional.
// Undeclared flags fall back to consuming the next non-flag token, which is
// what a user typing `--repo Acme/widget` means.
func takesValue(cmd *core.Command, key string) bool {
	if cmd == nil {
		return true
	}
	declared := cmd.Flags.Get(key)
	if !declared.OK {
		return true
	}
	_, isBool := declared.Value.(bool)
	return !isBool
}

// ParseOptions converts the arguments following a command into core.Options.
//
// Recognised forms:
//
//	--key=value  -k=value   value joined to the flag
//	--key value  -k value   value as the next token
//	--key        -k         boolean true (at the end, or before another flag)
//	--                      every argument after this is positional
//
// Positionals are collected under ArgsKey, with the first also under "_arg".
//
//	opts := cli.ParseOptions([]string{"--repo", "Acme/widget", "--force"})
//	opts.String("repo") // "Acme/widget"
//	opts.Bool("force")  // true
func ParseOptions(args []string) core.Options {
	return parseOptionsFor(nil, args)
}

func parseOptionsFor(cmd *core.Command, args []string) core.Options {
	opts := core.NewOptions()
	positional := make([]string, 0, len(args))
	terminated := false

	for index := 0; index < len(args); index++ {
		arg := args[index]

		if terminated {
			positional = append(positional, arg)
			continue
		}

		// Checked before ParseFlag, which reads "--" as a malformed flag.
		if arg == "--" {
			terminated = true
			continue
		}

		key, value, valid := core.ParseFlag(arg)
		if !valid {
			if core.IsFlag(arg) {
				// Malformed ("-verbose", "--v"). Dropping it silently is what
				// core/go did; a negative number is the reason not to treat it
				// as positional.
				continue
			}
			positional = append(positional, arg)
			continue
		}

		if core.Contains(arg, "=") {
			opts.Set(key, value)
			continue
		}

		next := index + 1
		if next < len(args) && args[next] != "--" && !core.IsFlag(args[next]) && takesValue(cmd, key) {
			opts.Set(key, args[next])
			index = next
			continue
		}

		opts.Set(key, true)
	}

	if len(positional) > 0 {
		opts.Set(ArgsKey, positional)
		opts.Set("_arg", positional[0])
	}

	return opts
}

// Args returns every positional argument passed to a command.
//
// core/go's parser kept only one, so a command wanting two values had to ask
// for them as flags. Commands reading opts.String("_arg") still see the first.
//
//	func action(opts core.Options) core.Result {
//	    args := cli.Args(opts) // ["dev.editor", "vim"]
//	}
func Args(opts core.Options) []string {
	result := opts.Get(ArgsKey)
	if !result.OK {
		return nil
	}
	args, ok := result.Value.([]string)
	if !ok {
		return nil
	}
	return args
}

// Arg returns the positional argument at index, or "" when there is none.
//
//	key := cli.Arg(opts, 0)
//	value := cli.Arg(opts, 1)
func Arg(opts core.Options, index int) string {
	args := Args(opts)
	if index < 0 || index >= len(args) {
		return ""
	}
	return args[index]
}
