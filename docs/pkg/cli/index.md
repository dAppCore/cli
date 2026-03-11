---
title: CLI Framework Overview
description: Go CLI framework built on cobra, with styled output, streaming, daemon mode, and TUI components.
---

# CLI Framework (`pkg/cli`)

`pkg/cli` is the CLI framework that powers the `core` binary and all derivative binaries in the ecosystem. It wraps cobra with the Core DI framework, adding styled output, streaming, interactive prompts, daemon management, and TUI components.

**Import:** `forge.lthn.ai/core/cli/pkg/cli`

## Quick Start

```go
package main

import "forge.lthn.ai/core/cli/pkg/cli"

func main() {
    cli.WithAppName("myapp")
    cli.Main(
        cli.WithCommands("greet", addGreetCommands),
    )
}

func addGreetCommands(root *cli.Command) {
    cmd := cli.NewCommand("greet", "Say hello", "", func(cmd *cli.Command, args []string) error {
        cli.Success("Hello, world!")
        return nil
    })
    root.AddCommand(cmd)
}
```

## Architecture

The framework has three layers:

1. **Runtime** (`runtime.go`, `app.go`) -- Initialises the Core DI container, cobra root command, and signal handling. Provides global accessors (`cli.Core()`, `cli.Context()`, `cli.RootCmd()`).

2. **Command registration** (`commands.go`, `command.go`) -- Two mechanisms for adding commands: `WithCommands` (lifecycle-aware, preferred) and `RegisterCommands` (init-time, for ecosystem packages).

3. **Output & interaction** (`output.go`, `stream.go`, `prompt.go`, `utils.go`) -- Styled output functions, streaming text renderer, interactive prompts, tables, trees, and task trackers.

## Key Types

| Type | Description |
|------|-------------|
| `Command` | Re-export of `cobra.Command` |
| `Stream` | Token-by-token text renderer with optional word-wrap |
| `Table` | Aligned tabular output with optional box-drawing borders |
| `TreeNode` | Tree structure with box-drawing connectors |
| `TaskTracker` | Concurrent task display with live spinners |
| `CheckBuilder` | Fluent API for pass/fail/skip result lines |
| `AnsiStyle` | Terminal text styling (bold, dim, colour) |

## Built-in Services

When you call `cli.Main()`, these services are registered automatically:

| Service | Name | Purpose |
|---------|------|---------|
| I18nService | `i18n` | Internationalisation and grammar composition |
| LogService | `log` | Structured logging with CLI-styled output |
| openpgp | `crypt` | OpenPGP encryption |
| workspace | `workspace` | Project root detection |
| signalService | `signal` | SIGINT/SIGTERM/SIGHUP handling |

## Documentation

- [Getting Started](getting-started.md) -- `Main()`, `WithCommands()`, building binaries
- [Commands](commands.md) -- Command builders, flag helpers, args validation
- [Output](output.md) -- Styled output, tables, trees, task trackers
- [Prompts](prompts.md) -- Interactive prompts, confirmations, selections
- [Streaming](streaming.md) -- Real-time token streaming with word-wrap
- [Daemon](daemon.md) -- Daemon mode, PID files, health checks
- [Errors](errors.md) -- Error creation, wrapping, exit codes

## Colour & Theme Control

Colours are enabled by default and respect the `NO_COLOR` environment variable and `TERM=dumb`. You can also control them programmatically:

```go
cli.SetColorEnabled(false)  // Disable ANSI colours
cli.UseASCII()              // ASCII glyphs + disable colours
cli.UseEmoji()              // Emoji glyph theme
cli.UseUnicode()            // Default Unicode glyph theme
```

## Execution Modes

The framework auto-detects the execution environment:

```go
mode := cli.DetectMode()
// cli.ModeInteractive -- TTY attached, colours enabled
// cli.ModePipe         -- stdout piped, colours disabled
// cli.ModeDaemon       -- CORE_DAEMON=1, log-only output

cli.IsTTY()       // stdout is a terminal?
cli.IsStdinTTY()  // stdin is a terminal?
cli.IsStderrTTY() // stderr is a terminal?
```
