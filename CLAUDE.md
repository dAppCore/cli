# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Core CLI (`forge.lthn.ai/core/cli`) is a Go CLI tool for managing development workflows across Go, PHP, and Wails projects. It wraps common tooling (testing, linting, building, releasing, multi-repo management) behind a unified `core` command. Built on Cobra with the Core framework's dependency injection for service lifecycle management.

## Build & Development Commands

```bash
# Run all tests
go test ./...

# Run a single test
go test -run TestName ./...

# Build CLI binary
go build -o dist/core .

# Install from source
go install forge.lthn.ai/core/cli@latest

# Verify environment
go run . doctor

# Format and lint
gofmt -w .
go vet ./...
golangci-lint run ./...
```

**Go version:** 1.26+

## Architecture

### Entry Point & Command Registration

`main.go` wires everything together. Commands register in two ways:

1. **Explicit registration** via `cli.WithCommands()` — local command packages in `cmd/` pass an `AddXCommands(root *cli.Command)` function that receives the root cobra command during service startup.

2. **Self-registration** via `cli.RegisterCommands()` — ecosystem packages (imported as blank `_` imports in `main.go`) call `cli.RegisterCommands()` in their `init()` functions. This is how external modules like `go-build`, `go-devops`, `go-scm`, `agent`, etc. contribute commands without coupling to `main.go`.

### `pkg/cli` — The CLI Runtime

This is the core package. Everything commands need is re-exported here to avoid direct Cobra/lipgloss/bubbletea imports:

- **`runtime.go`**: Singleton `Init()`/`Shutdown()`/`Execute()` lifecycle. Creates a Core framework instance, attaches services, handles signals (SIGINT/SIGTERM/SIGHUP).
- **`command.go`**: `NewCommand()`, `NewGroup()`, `NewRun()` builders, flag helpers (`StringFlag`, `BoolFlag`, etc.), arg validators — all wrapping Cobra types so command packages don't import cobra directly. `Command` is a type alias for `cobra.Command`.
- **`output.go`**: Styled output functions — `Success()`, `Error()`, `Warn()`, `Info()`, `Dim()`, `Progress()`, `Label()`, `Section()`, `Hint()`, `Severity()`. Use these instead of raw `fmt.Print`.
- **`errors.go`**: `Err()`, `Wrap()`, `WrapVerb()`, `Exit()` for error creation. Re-exports `errors.Is`/`As`/`Join`. Commands should return errors (via `RunE`), not call `Fatal()` (deprecated).
- **`frame.go`**: Bubbletea-based TUI framework. `NewFrame("HCF")` with HLCRF region layout (Header, Left, Content, Right, Footer), focus management, content navigation stack.
- **`tracker.go`**: `TaskTracker` for concurrent task display with spinners. TTY-aware (live updates vs static output).
- **`daemon.go`**: Execution mode detection (`ModeInteractive`/`ModePipe`/`ModeDaemon`).
- **`styles.go`**: Shared lipgloss styles and colour constants.
- **`glyph.go`**: Shortcode system for emoji/symbols (`:check:`, `:cross:`, `:warn:`, `:info:`).

### Command Package Pattern

Every command package in `cmd/` follows this structure:

```go
package mycommand

import "forge.lthn.ai/core/cli/pkg/cli"

func AddMyCommands(root *cli.Command) {
    myCmd := cli.NewGroup("my", "My commands", "")
    root.AddCommand(myCmd)
    // Add subcommands...
}
```

Commands use `cli.NewCommand()` (returns error) or `cli.NewRun()` (no error) and the `cli.StringFlag()`/`cli.BoolFlag()` helpers for flags.

### External Module Integration

The CLI imports ecosystem modules as blank imports that self-register via `init()`:
- `forge.lthn.ai/core/go-build` — build, CI, SDK commands
- `forge.lthn.ai/core/go-devops` — dev, deploy, docs, git, setup commands
- `forge.lthn.ai/core/go-scm` — forge, gitea, collect commands
- `forge.lthn.ai/core/agent` — agent, dispatch, task commands
- `forge.lthn.ai/core/lint` — QA commands
- Others: go-ansible, go-api, go-container, go-crypt, go-infra

### Daemon/Service Management

`cmd/service/` implements `start`/`stop`/`list`/`restart` for manifest-driven daemons. Reads `.core/manifest.yaml` from the project directory (walks up). Daemons run detached with `CORE_DAEMON=1` env var and are tracked in `~/.core/daemons/`.

## Test Naming Convention

Tests use `_Good`, `_Bad`, `_Ugly` suffix pattern:
- `_Good`: Happy path tests
- `_Bad`: Expected error conditions
- `_Ugly`: Panic/edge cases

## Commit Message Convention

[Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`

## Configuration

- Global config: `~/.core/config.yaml` (YAML, dot-notation keys)
- Project config: `.core/build.yaml`, `.core/release.yaml`, `.core/ci.yaml`
- Environment override: `CORE_CONFIG_<KEY>` (underscores become dots, lowercased)
- Multi-repo registry: `repos.yaml` (searched cwd upward, then `~/.config/core/repos.yaml`)

## i18n

Commands use `i18n.T("key")` for translatable strings and the grammar system (`i18n.ActionFailed()`, `i18n.Progress()`) for consistent error/progress messages. The CLI wraps these in `cli.Echo()`, `cli.WrapVerb()`, `cli.ErrorWrapVerb()`.
