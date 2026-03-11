---
title: Error Handling
description: Error creation, wrapping with i18n grammar, exit codes, and inspection helpers.
---

# Error Handling

The framework provides error creation and wrapping functions that integrate with i18n grammar composition, plus re-exports of the standard `errors` package for convenience.

## Error Creation

```go
// Simple error (replaces fmt.Errorf)
return cli.Err("invalid model: %s", name)

// Wrap with context (nil-safe -- returns nil if err is nil)
return cli.Wrap(err, "load config")  // "load config: <original>"

// Wrap with i18n grammar
return cli.WrapVerb(err, "load", "config")    // "Failed to load config: <original>"
return cli.WrapAction(err, "connect")          // "Failed to connect: <original>"
```

`WrapVerb` and `WrapAction` use the i18n `ActionFailed` function, which produces grammatically correct error messages across languages.

All wrapping functions are nil-safe: they return `nil` if the input error is `nil`.

## Error Inspection

Re-exports of the `errors` package for convenience, so you do not need a separate import:

```go
if cli.Is(err, os.ErrNotExist) { ... }

var exitErr *cli.ExitError
if cli.As(err, &exitErr) {
    os.Exit(exitErr.Code)
}

combined := cli.Join(err1, err2, err3)
```

## Exit Codes

Return a specific exit code from a command:

```go
return cli.Exit(2, fmt.Errorf("validation failed"))
```

The `ExitError` type wraps an error with an exit code. `Main()` checks for `*ExitError` and uses its code when exiting the process. All other errors exit with code 1.

```go
type ExitError struct {
    Code int
    Err  error
}
```

`ExitError` implements `error` and `Unwrap()`, so it works with `errors.Is` and `errors.As`.

## The Pattern: Commands Return Errors

Commands should return errors rather than calling `os.Exit` directly. `Main()` handles the exit:

```go
// Correct: return error, let Main() handle exit
func runBuild(cmd *cli.Command, args []string) error {
    if err := compile(); err != nil {
        return cli.WrapVerb(err, "compile", "project")
    }
    cli.Success("Build complete")
    return nil
}

// Wrong: calling os.Exit from command code
func runBuild(cmd *cli.Command, args []string) error {
    if err := compile(); err != nil {
        cli.Fatal(err) // Do not do this
    }
    return nil
}
```

## Fatal Functions (Deprecated)

These exist for legacy code but should not be used in new commands:

```go
cli.Fatal(err)                        // prints + os.Exit(1)
cli.Fatalf("bad: %v", err)           // prints + os.Exit(1)
cli.FatalWrap(err, "load config")    // prints + os.Exit(1)
cli.FatalWrapVerb(err, "load", "x")  // prints + os.Exit(1)
```

All `Fatal*` functions log the error, print it to stderr with the error style, and call `os.Exit(1)`. They bypass the normal shutdown sequence.

## Error Output Functions

For displaying errors without returning them (see [Output](output.md) for full details):

```go
cli.Error("message")                         // ✗ message (stderr)
cli.Errorf("port %d in use", port)          // ✗ port 8080 in use (stderr)
cli.ErrorWrap(err, "context")               // ✗ context: <error> (stderr)
cli.ErrorWrapVerb(err, "load", "config")    // ✗ Failed to load config: <error> (stderr)
cli.ErrorWrapAction(err, "connect")         // ✗ Failed to connect: <error> (stderr)
```

These print styled error messages but do not terminate the process or return an error value. Use them when you want to report a problem and continue.
