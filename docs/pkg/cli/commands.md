---
title: Command Builders
description: Creating commands, flag helpers, args validation, and the config struct pattern.
---

# Command Builders

The framework provides three command constructors and a full set of flag helpers. All wrap cobra but remove the need to import it directly.

## Command Types

### `NewCommand` -- Standard command (returns error)

The most common form. The handler returns an error which `Main()` handles:

```go
cmd := cli.NewCommand("build", "Build the project", "", func(cmd *cli.Command, args []string) error {
    if err := compile(); err != nil {
        return cli.WrapVerb(err, "compile", "project")
    }
    cli.Success("Build complete")
    return nil
})
```

The third parameter is the long description (shown in `--help`). Pass `""` to omit it.

### `NewGroup` -- Parent command (subcommands only)

Creates a command with no handler, used to group subcommands:

```go
scoreCmd := cli.NewGroup("score", "Scoring commands", "")
scoreCmd.AddCommand(grammarCmd, attentionCmd, tierCmd)
root.AddCommand(scoreCmd)
```

### `NewRun` -- Simple command (no error return)

For commands that cannot fail:

```go
cmd := cli.NewRun("version", "Show version", "", func(cmd *cli.Command, args []string) {
    cli.Println("v1.0.0")
})
```

## Re-exports

The framework re-exports cobra types so you never need to import cobra directly:

```go
cli.Command        // = cobra.Command
cli.PositionalArgs // = cobra.PositionalArgs
```

## Flag Helpers

All flag helpers follow the same signature: `(cmd, ptr, name, short, default, usage)`. Pass `""` for the short name to omit the short flag.

```go
var cfg struct {
    Model    string
    Verbose  bool
    Count    int
    Score    float64
    Seed     int64
    Timeout  time.Duration
    Tags     []string
}

cli.StringFlag(cmd, &cfg.Model, "model", "m", "", "Model path")
cli.BoolFlag(cmd, &cfg.Verbose, "verbose", "v", false, "Verbose output")
cli.IntFlag(cmd, &cfg.Count, "count", "n", 10, "Item count")
cli.Float64Flag(cmd, &cfg.Score, "score", "s", 0.0, "Min score")
cli.Int64Flag(cmd, &cfg.Seed, "seed", "", 0, "Random seed")
cli.DurationFlag(cmd, &cfg.Timeout, "timeout", "t", 30*time.Second, "Timeout")
cli.StringSliceFlag(cmd, &cfg.Tags, "tag", "", nil, "Tags")
```

### Persistent Flags

Persistent flags are inherited by all subcommands:

```go
cli.PersistentStringFlag(parentCmd, &dbPath, "db", "d", "", "Database path")
cli.PersistentBoolFlag(parentCmd, &debug, "debug", "", false, "Debug mode")
```

## Args Validation

Constrain the number of positional arguments:

```go
cmd := cli.NewCommand("deploy", "Deploy to env", "", deployFn)
cli.WithArgs(cmd, cli.ExactArgs(1))    // Exactly 1 arg
cli.WithArgs(cmd, cli.MinimumNArgs(1)) // At least 1
cli.WithArgs(cmd, cli.MaximumNArgs(3)) // At most 3
cli.WithArgs(cmd, cli.RangeArgs(1, 3)) // Between 1 and 3
cli.WithArgs(cmd, cli.NoArgs())        // No args allowed
cli.WithArgs(cmd, cli.ArbitraryArgs()) // Any number of args
```

## Command Configuration

Add examples to help text:

```go
cli.WithExample(cmd, `  core build --targets linux/amd64
  core build --ci`)
```

## Pattern: Config Struct + Flags

The idiomatic pattern for commands with many flags is to define a config struct, bind flags to its fields, then pass the struct to the business logic:

```go
type DistillOpts struct {
    Model    string
    Probes   string
    Runs     int
    DryRun   bool
}

func addDistillCommand(parent *cli.Command) {
    var cfg DistillOpts

    cmd := cli.NewCommand("distill", "Run distillation", "", func(cmd *cli.Command, args []string) error {
        return RunDistill(cfg)
    })

    cli.StringFlag(cmd, &cfg.Model, "model", "m", "", "Model config path")
    cli.StringFlag(cmd, &cfg.Probes, "probes", "p", "", "Probe set name")
    cli.IntFlag(cmd, &cfg.Runs, "runs", "r", 3, "Runs per probe")
    cli.BoolFlag(cmd, &cfg.DryRun, "dry-run", "", false, "Preview without executing")

    parent.AddCommand(cmd)
}
```

## Registration Function Pattern

Commands are organised in packages under `cmd/`. Each package exports an `Add*Commands` function:

```go
// cmd/score/commands.go
package score

import "forge.lthn.ai/core/cli/pkg/cli"

func AddScoreCommands(root *cli.Command) {
    scoreCmd := cli.NewGroup("score", "Scoring commands", "")

    grammarCmd := cli.NewCommand("grammar", "Grammar analysis", "", runGrammar)
    cli.StringFlag(grammarCmd, &inputPath, "input", "i", "", "Input file")
    scoreCmd.AddCommand(grammarCmd)

    root.AddCommand(scoreCmd)
}
```

Then in `main.go`:

```go
cli.Main(
    cli.WithCommands("score", score.AddScoreCommands),
)
```
