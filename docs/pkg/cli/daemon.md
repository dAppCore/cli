---
title: Daemon Mode
description: Daemon process management, PID files, health checks, and execution modes.
---

# Daemon Mode

The framework provides execution mode detection and signal handling for daemon processes.

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

## Simple Daemon

Use `cli.Context()` for cancellation-aware daemon loops:

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
