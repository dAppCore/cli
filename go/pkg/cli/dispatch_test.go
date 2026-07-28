package cli

import (
	core "dappco.re/go"
)

// The regression this file exists for: core/go's parser recognised only
// `--key=value`, so `--repo Acme/widget` set repo to true and pushed the value
// into "_arg". Commands read opts.String("repo"), got "", and silently used
// their default.
func TestDispatch_ParseOptions_Good(t *core.T) {
	opts := ParseOptions([]string{"--repo", "Acme/widget", "--version", "v9.9.9"})

	core.AssertEqual(t, "Acme/widget", opts.String("repo"))
	core.AssertEqual(t, "v9.9.9", opts.String("version"))
}

func TestDispatch_ParseOptions_Joined_Good(t *core.T) {
	opts := ParseOptions([]string{"--repo=Acme/widget", "-n=widget"})

	core.AssertEqual(t, "Acme/widget", opts.String("repo"))
	core.AssertEqual(t, "widget", opts.String("n"))
}

// A flag with nothing after it, or another flag after it, is a switch.
func TestDispatch_ParseOptions_Boolean_Good(t *core.T) {
	opts := ParseOptions([]string{"--force", "--dry-run", "--verbose"})

	core.AssertEqual(t, true, opts.Bool("force"))
	core.AssertEqual(t, true, opts.Bool("dry-run"))
	core.AssertEqual(t, true, opts.Bool("verbose"))
}

func TestDispatch_ParseOptions_Positional_Good(t *core.T) {
	opts := ParseOptions([]string{"dev.editor", "vim"})

	// "_arg" is the first, not the last: `config get KEY` wants KEY.
	core.AssertEqual(t, "dev.editor", opts.String("_arg"))
	core.AssertEqual(t, "dev.editor", Arg(opts, 0))
	core.AssertEqual(t, "vim", Arg(opts, 1))
	core.AssertEqual(t, 2, len(Args(opts)))
}

func TestDispatch_ParseOptions_Terminator_Good(t *core.T) {
	opts := ParseOptions([]string{"--repo", "Acme/widget", "--", "--not-a-flag", "-x"})

	core.AssertEqual(t, "Acme/widget", opts.String("repo"))
	core.AssertEqual(t, "--not-a-flag", Arg(opts, 0))
	core.AssertEqual(t, "-x", Arg(opts, 1))
	// Everything after "--" is positional, so -x must not become a flag.
	core.AssertEqual(t, false, opts.Has("x"))
}

// A flag declared bool in Command.Flags must not swallow the token after it,
// so `--verbose ./path` keeps ./path positional.
func TestDispatch_ParseOptions_DeclaredBool_Good(t *core.T) {
	cmd := &core.Command{Flags: core.NewOptions(core.Option{Key: "verbose", Value: false})}
	opts := parseOptionsFor(cmd, []string{"--verbose", "./path"})

	core.AssertEqual(t, true, opts.Bool("verbose"))
	core.AssertEqual(t, "./path", opts.String("_arg"))
}

func TestDispatch_ParseOptions_Bad(t *core.T) {
	// Malformed flags: single dash with 2+ chars, double dash with 1.
	opts := ParseOptions([]string{"-verbose", "--v"})

	core.AssertEqual(t, false, opts.Has("verbose"))
	core.AssertEqual(t, false, opts.Has("v"))
	core.AssertEqual(t, 0, len(Args(opts)))
}

func TestDispatch_ParseOptions_Ugly(t *core.T) {
	core.AssertEqual(t, 0, ParseOptions(nil).Len())
	core.AssertEqual(t, 0, len(Args(ParseOptions(nil))))
	core.AssertEqual(t, "", Arg(ParseOptions(nil), 0))
	core.AssertEqual(t, "", Arg(ParseOptions([]string{"only"}), -1))
	core.AssertEqual(t, "", Arg(ParseOptions([]string{"only"}), 9))
}

// A trailing flag has no value to take, so it stays a switch rather than
// consuming past the end of the slice.
func TestDispatch_ParseOptions_TrailingFlag_Ugly(t *core.T) {
	opts := ParseOptions([]string{"--repo", "Acme/widget", "--force"})

	core.AssertEqual(t, "Acme/widget", opts.String("repo"))
	core.AssertEqual(t, true, opts.Bool("force"))
}

func TestDispatch_resolveCommand_Good(t *core.T) {
	c := core.New()
	c.Command("build", core.Command{Action: func(core.Options) core.Result { return core.Ok(nil) }})
	c.Command("build/installers", core.Command{Action: func(core.Options) core.Result { return core.Ok(nil) }})

	// Longest match wins, or "installers" would be a positional arg to "build".
	cmd, rest := resolveCommand(c, []string{"build", "installers", "--variant", "ci"})

	core.AssertEqual(t, "build/installers", cmd.Path)
	core.AssertEqual(t, 2, len(rest))
	core.AssertEqual(t, "--variant", rest[0])
}

func TestDispatch_resolveCommand_Shortest_Good(t *core.T) {
	c := core.New()
	c.Command("build", core.Command{Action: func(core.Options) core.Result { return core.Ok(nil) }})

	cmd, rest := resolveCommand(c, []string{"build", "somewhere"})

	core.AssertEqual(t, "build", cmd.Path)
	core.AssertEqual(t, 1, len(rest))
	core.AssertEqual(t, "somewhere", rest[0])
}

func TestDispatch_resolveCommand_Bad(t *core.T) {
	c := core.New()
	c.Command("build", core.Command{Action: func(core.Options) core.Result { return core.Ok(nil) }})

	cmd, _ := resolveCommand(c, []string{"floopadoop"})

	core.AssertEqual(t, true, cmd == nil)
}

func TestDispatch_resolveCommand_Ugly(t *core.T) {
	cmd, rest := resolveCommand(nil, []string{"build"})
	core.AssertEqual(t, true, cmd == nil)
	core.AssertEqual(t, 1, len(rest))

	c := core.New()
	// A group placeholder has no Action and must not resolve as runnable.
	c.Command("group/leaf", core.Command{Action: func(core.Options) core.Result { return core.Ok(nil) }})
	groupCmd, _ := resolveCommand(c, []string{"group"})
	core.AssertEqual(t, true, groupCmd == nil)
}
