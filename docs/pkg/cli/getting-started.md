---
title: Getting Started
description: How to use cli.Main(), WithCommands(), and build CLI binaries with the framework.
---

# Getting Started

## The `core` Binary

The `core` binary is built from `main.go` at the repo root. It composes commands from both local packages and ecosystem modules:

```go
package main

import (
    "forge.lthn.ai/core/cli/cmd/config"
    "forge.lthn.ai/core/cli/cmd/gocmd"
    "forge.lthn.ai/core/cli/pkg/cli"

    // Ecosystem packages self-register via init()
    _ "forge.lthn.ai/core/go-devops/cmd/dev"
    _ "forge.lthn.ai/core/go-build/cmd/build"
)

func main() {
    cli.Main(
        cli.WithCommands("config", config.AddConfigCommands),
        cli.WithCommands("go", gocmd.AddGoCommands),
    )
}
```

## `cli.Main()`

`Main()` is the primary entry point. It:

1. Registers core services (i18n, log, crypt, workspace)
2. Appends your command services
3. Creates the cobra root command and signal handler
4. Starts all services via the Core DI framework
5. Adds the `completion` command
6. Executes the matched command
7. Shuts down all services in reverse order
8. Exits with the appropriate code

```go
cli.Main(
    cli.WithCommands("score", score.AddScoreCommands),
    cli.WithCommands("gen", gen.AddGenCommands),
)
```

If a command returns an `*ExitError`, the process exits with that code. All other errors exit with code 1.

## `cli.WithCommands()`

This is the preferred way to register commands. It wraps your registration function in a Core service that participates in the lifecycle:

```go
func WithCommands(name string, register func(root *Command), localeFS ...fs.FS) CommandSetup
```

During `Main()`, the CLI calls your function with the Core instance. Internally it retrieves the root cobra command and passes it to your register function:

```go
func AddScoreCommands(root *cli.Command) {
    scoreCmd := cli.NewGroup("score", "Scoring commands", "")

    grammarCmd := cli.NewCommand("grammar", "Grammar analysis", "", runGrammar)
    cli.StringFlag(grammarCmd, &inputPath, "input", "i", "", "Input file")

    scoreCmd.AddCommand(grammarCmd)
    root.AddCommand(scoreCmd)
}
```

**Startup order:**
1. Core services start (i18n, log, crypt, workspace, signal)
2. Command services start (your `WithCommands` functions run)
3. `Execute()` runs the matched command

## Building a Variant Binary

To create a standalone binary (not the `core` binary), set the app name and compose your commands:

```go
// cmd/lem/main.go
package main

import (
    "forge.lthn.ai/core/cli/pkg/cli"
    "forge.lthn.ai/lthn/lem/cmd/lemcmd"
)

func main() {
    cli.WithAppName("lem")
    cli.Main(lemcmd.Commands()...)
}
```

Where `Commands()` returns a slice of `CommandSetup` functions:

```go
package lemcmd

import (
    "forge.lthn.ai/core/cli/pkg/cli"
)

func Commands() []cli.CommandSetup {
    return []cli.CommandSetup{
        cli.WithCommands("score", addScoreCommands),
        cli.WithCommands("gen", addGenCommands),
        cli.WithCommands("data", addDataCommands),
    }
}
```

## `cli.RegisterCommands()` (Legacy)

For ecosystem packages that need to self-register via `init()`:

```go
func init() {
    cli.RegisterCommands(func(root *cobra.Command) {
        root.AddCommand(myCmd)
    })
}
```

The `core` binary imports these packages with blank imports (`_ "forge.lthn.ai/core/go-build/cmd/build"`), triggering their `init()` functions.

**Prefer `WithCommands`** -- it is explicit and does not rely on import side effects.

## Manual Initialisation (Advanced)

If you need more control over the lifecycle:

```go
cli.Init(cli.Options{
    AppName:  "myapp",
    Version:  "1.0.0",
    Services: []core.Service{...},
    OnReload: func() error { return reloadConfig() },
})
defer cli.Shutdown()

// Add commands manually
cli.RootCmd().AddCommand(myCmd)

if err := cli.Execute(); err != nil {
    os.Exit(1)
}
```

## Version Info

Version fields are set via ldflags at build time:

```go
cli.AppVersion      // "1.2.0"
cli.BuildCommit     // "df94c24"
cli.BuildDate       // "2026-02-06"
cli.BuildPreRelease // "dev.8"
cli.SemVer()        // "1.2.0-dev.8+df94c24.20260206"
```

Build command:

```bash
go build -ldflags="-X forge.lthn.ai/core/cli/pkg/cli.AppVersion=1.2.0 \
  -X forge.lthn.ai/core/cli/pkg/cli.BuildCommit=$(git rev-parse --short HEAD) \
  -X forge.lthn.ai/core/cli/pkg/cli.BuildDate=$(date +%Y-%m-%d)"
```

## Accessing Core Services

Inside a command handler, you can access the Core DI container and retrieve services:

```go
func runMyCommand(cmd *cli.Command, args []string) error {
    ctx := cli.Context()     // Root context (cancelled on signal)
    core := cli.Core()       // Framework Core instance
    root := cli.RootCmd()    // Root cobra command

    // Type-safe service retrieval
    ws, err := framework.ServiceFor[*workspace.Service](core)
    if err != nil {
        return cli.WrapVerb(err, "get", "workspace service")
    }

    return nil
}
```

## Signal Handling

Signal handling is automatic. SIGINT and SIGTERM cancel `cli.Context()`. Use this context in your commands for graceful cancellation:

```go
func runServer(cmd *cli.Command, args []string) error {
    ctx := cli.Context()
    // ctx is cancelled when the user presses Ctrl+C
    <-ctx.Done()
    return nil
}
```

Optional SIGHUP handling for configuration reload:

```go
cli.Init(cli.Options{
    AppName: "daemon",
    OnReload: func() error {
        return reloadConfig()
    },
})
```
