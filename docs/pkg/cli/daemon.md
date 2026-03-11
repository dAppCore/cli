---
title: Daemon Mode
description: Daemon process management, PID files, health checks, and execution modes.
---

# Daemon Mode

The framework provides both low-level daemon primitives and a high-level command group that adds `start`, `stop`, `status`, and `run` subcommands to your CLI.

## Execution Modes

The framework auto-detects the execution environment:

```go
mode := cli.DetectMode()
```

| Mode | Condition | Behaviour |
|------|-----------|-----------|
| `ModeInteractive` | TTY attached | Colours enabled, spinners active |
| `ModePipe` | stdout piped | Colours disabled, plain output |
| `ModeDaemon` | `CORE_DAEMON=1` env var | Log-only output |

Helper functions:

```go
cli.IsTTY()       // stdout is a terminal?
cli.IsStdinTTY()  // stdin is a terminal?
cli.IsStderrTTY() // stderr is a terminal?
```

## Adding Daemon Commands

`AddDaemonCommand` registers a command group with four subcommands:

```go
func AddMyCommands(root *cli.Command) {
    cli.AddDaemonCommand(root, cli.DaemonCommandConfig{
        Name:        "daemon",            // Command group name (default: "daemon")
        Description: "Manage the worker", // Short description
        PIDFile:     "/var/run/myapp.pid",
        HealthAddr:  ":9090",
        RunForeground: func(ctx context.Context, daemon *process.Daemon) error {
            // Your long-running service logic here.
            // ctx is cancelled on SIGINT/SIGTERM.
            return runWorker(ctx)
        },
    })
}
```

This creates:

- `myapp daemon start` -- Re-executes the binary as a background process with `CORE_DAEMON=1`
- `myapp daemon stop` -- Sends SIGTERM to the daemon, waits for shutdown (30s timeout, then SIGKILL)
- `myapp daemon status` -- Reports whether the daemon is running and queries health endpoints
- `myapp daemon run` -- Runs in the foreground (for development or process managers like systemd)

### Custom Persistent Flags

Add flags that apply to all daemon subcommands:

```go
cli.AddDaemonCommand(root, cli.DaemonCommandConfig{
    // ...
    Flags: func(cmd *cli.Command) {
        cli.PersistentStringFlag(cmd, &configPath, "config", "c", "", "Config file")
    },
    ExtraStartArgs: func() []string {
        return []string{"--config", configPath}
    },
})
```

`ExtraStartArgs` passes additional flags when re-executing the binary as a daemon.

### Health Endpoints

When `HealthAddr` is set, the daemon serves:

- `GET /health` -- Liveness check (200 if server is up, 503 if health checks fail)
- `GET /ready` -- Readiness check (200 if `daemon.SetReady(true)` has been called)

The `start` command waits up to 5 seconds for the health endpoint to become available before reporting success.

## Simple Daemon (Manual)

For cases where you do not need the full command group:

```go
func runDaemon(cmd *cli.Command, args []string) error {
    ctx := cli.Context() // Cancelled on SIGINT/SIGTERM
    // ... start your work ...
    <-ctx.Done()
    return nil
}
```

## Shutdown with Timeout

The daemon stop logic sends SIGTERM and waits up to 30 seconds. If the process has not exited by then, it sends SIGKILL and removes the PID file.

## Signal Handling

Signal handling is built into the CLI runtime:

- **SIGINT/SIGTERM** cancel `cli.Context()`
- **SIGHUP** calls the `OnReload` handler if configured:

```go
cli.Init(cli.Options{
    AppName: "daemon",
    OnReload: func() error {
        return reloadConfig()
    },
})
```

No manual signal handling is needed in commands. Use `cli.Context()` for cancellation-aware operations.

## DaemonCommandConfig Reference

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Command group name (default: `"daemon"`) |
| `Description` | `string` | Short description for help text |
| `PIDFile` | `string` | PID file path (default flag value) |
| `HealthAddr` | `string` | Health check listen address (default flag value) |
| `RunForeground` | `func(ctx, daemon) error` | Service logic for foreground/daemon mode |
| `Flags` | `func(cmd)` | Registers custom persistent flags |
| `ExtraStartArgs` | `func() []string` | Additional args for background re-exec |
