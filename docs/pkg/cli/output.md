---
title: Output Functions
description: Styled output, tables, trees, task trackers, glyphs, and formatting utilities.
---

# Output Functions

All output functions support glyph shortcodes (`:check:`, `:cross:`, `:warn:`, `:info:`) which are auto-converted to Unicode symbols.

## Styled Messages

```go
cli.Success("All tests passed")     // ✓ All tests passed (green, bold)
cli.Successf("Built %d files", n)   // ✓ Built 5 files (green, bold)
cli.Error("Connection refused")     // ✗ Connection refused (red, bold, stderr)
cli.Errorf("Port %d in use", port)  // ✗ Port 8080 in use (red, bold, stderr)
cli.Warn("Deprecated flag used")    // ⚠ Deprecated flag used (amber, bold, stderr)
cli.Warnf("Skipping %s", name)     // ⚠ Skipping foo (amber, bold, stderr)
cli.Info("Using default config")    // ℹ Using default config (blue)
cli.Infof("Found %d items", n)     // ℹ Found 42 items (blue)
cli.Dim("Optional detail")         // Optional detail (grey, dimmed)
```

`Error` and `Warn` write to stderr and also log the message. `Success` and `Info` write to stdout.

## Plain Output

```go
cli.Println("Hello %s", name)      // fmt.Sprintf + glyph conversion + newline
cli.Print("Loading...")             // No newline
cli.Text("raw", "text")            // Like fmt.Println but with glyphs
cli.Blank()                         // Empty line
cli.Echo("key.label", args...)     // i18n.T translation + newline
```

## Structured Output

```go
cli.Label("version", "1.2.0")     // Version: 1.2.0 (styled key)
cli.Task("php", "Running tests")  // [php] Running tests
cli.Section("audit")               // ── AUDIT ──
cli.Hint("fix", "go mod tidy")    //   fix: go mod tidy
cli.Result(passed, "Tests passed") // ✓ or ✗ based on bool
```

## Error Wrapping for Output

These display errors without returning them -- useful when you need to show an error but continue:

```go
cli.ErrorWrap(err, "load config")            // ✗ load config: <error>
cli.ErrorWrapVerb(err, "load", "config")     // ✗ Failed to load config: <error>
cli.ErrorWrapAction(err, "connect")          // ✗ Failed to connect: <error>
```

All three are nil-safe -- they do nothing if `err` is nil.

## Progress Indicator

Overwrites the current terminal line to show progress:

```go
for i, item := range items {
    cli.Progress("check", i+1, len(items), item.Name) // Overwrites line
}
cli.ProgressDone() // Clears progress line
```

The verb is passed through `i18n.Progress()` for gerund form ("Checking...").

## Severity Levels

```go
cli.Severity("critical", "SQL injection found")  // [critical] red, bold
cli.Severity("high", "XSS vulnerability")        // [high] orange, bold
cli.Severity("medium", "Missing CSRF token")     // [medium] amber
cli.Severity("low", "Debug mode enabled")         // [low] grey
```

## Check Results

Fluent API for building pass/fail/skip/warn result lines:

```go
cli.Check("audit").Pass().Print()                  //   ✓ audit passed
cli.Check("fmt").Fail().Duration("2.3s").Print()   //   ✗ fmt    failed    2.3s
cli.Check("test").Skip().Print()                   //   - test skipped
cli.Check("lint").Warn().Message("3 warnings").Print()
```

## Tables

Aligned tabular output with optional box-drawing borders:

```go
t := cli.NewTable("REPO", "STATUS", "BRANCH")
t.AddRow("core-php", "clean", "main")
t.AddRow("core-tenant", "dirty", "feature/x")
t.Render()
```

Output:

```
REPO          STATUS  BRANCH
core-php      clean   main
core-tenant   dirty   feature/x
```

### Bordered Tables

```go
t := cli.NewTable("REPO", "STATUS").
    WithBorders(cli.BorderRounded)
t.AddRow("core-php", "clean")
t.Render()
```

Output:

```
╭──────────┬────────╮
│ REPO     │ STATUS │
├──────────┼────────┤
│ core-php │ clean  │
╰──────────┴────────╯
```

Border styles: `BorderNone` (default), `BorderNormal`, `BorderRounded`, `BorderHeavy`, `BorderDouble`.

### Per-Column Cell Styling

```go
t := cli.NewTable("REPO", "STATUS", "BRANCH").
    WithCellStyle(1, func(val string) *cli.AnsiStyle {
        if val == "clean" {
            return cli.SuccessStyle
        }
        return cli.WarningStyle
    }).
    WithMaxWidth(80)
```

## Trees

Hierarchical output with box-drawing connectors:

```go
tree := cli.NewTree("core-php")
tree.Add("core-tenant").Add("core-bio")
tree.Add("core-admin")
tree.Add("core-api")
tree.Render()
```

Output:

```
core-php
├── core-tenant
│   └── core-bio
├── core-admin
└── core-api
```

Styled nodes:

```go
tree.AddStyled("core-tenant", cli.RepoStyle)
```

## Task Tracker

Displays multiple concurrent tasks with live spinners. Uses ANSI cursor manipulation when connected to a TTY; falls back to line-by-line output otherwise.

```go
tracker := cli.NewTaskTracker()
for _, repo := range repos {
    t := tracker.Add(repo.Name)
    go func(t *cli.TrackedTask) {
        t.Update("pulling...")
        if err := pull(repo); err != nil {
            t.Fail(err.Error())
            return
        }
        t.Done("up to date")
    }(t)
}
tracker.Wait()
cli.Println(tracker.Summary()) // "5/5 passed"
```

`TrackedTask` methods (`Update`, `Done`, `Fail`) are safe for concurrent use.

## String Builders (No Print)

These return styled strings without printing, for use in composing output:

```go
cli.SuccessStr("done")    // Returns "✓ done" styled green
cli.ErrorStr("failed")    // Returns "✗ failed" styled red
cli.WarnStr("warning")    // Returns "⚠ warning" styled amber
cli.InfoStr("note")       // Returns "ℹ note" styled blue
cli.DimStr("detail")      // Returns dimmed text
cli.Styled(style, "text") // Apply any AnsiStyle
cli.Styledf(style, "formatted %s", arg)
```

## Glyph Shortcodes

All output functions auto-convert shortcodes to symbols:

| Shortcode | Unicode | Emoji | ASCII |
|-----------|---------|-------|-------|
| `:check:` | ✓ | ✅ | [OK] |
| `:cross:` | ✗ | ❌ | [FAIL] |
| `:warn:` | ⚠ | ⚠️ | [WARN] |
| `:info:` | ℹ | ℹ️ | [INFO] |
| `:arrow_right:` | → | ➡️ | -> |
| `:bullet:` | • | • | * |
| `:dash:` | ─ | ─ | - |
| `:pipe:` | │ | │ | \| |
| `:corner:` | └ | └ | \` |
| `:tee:` | ├ | ├ | + |
| `:pending:` | … | ⏳ | ... |

Switch themes:

```go
cli.UseUnicode() // Default
cli.UseEmoji()   // Emoji symbols
cli.UseASCII()   // ASCII fallback (also disables colours)
```

## ANSI Styles

Build custom styles with the fluent `AnsiStyle` API:

```go
style := cli.NewStyle().Bold().Foreground(cli.ColourBlue500)
fmt.Println(style.Render("Important text"))
```

Available methods: `Bold()`, `Dim()`, `Italic()`, `Underline()`, `Foreground(hex)`, `Background(hex)`.

The framework provides pre-defined styles using the Tailwind colour palette:

| Style | Description |
|-------|-------------|
| `SuccessStyle` | Green, bold |
| `ErrorStyle` | Red, bold |
| `WarningStyle` | Amber, bold |
| `InfoStyle` | Blue |
| `SecurityStyle` | Purple, bold |
| `DimStyle` | Grey, dimmed |
| `HeaderStyle` | Grey-200, bold |
| `AccentStyle` | Cyan |
| `LinkStyle` | Blue, underlined |
| `CodeStyle` | Grey-300 |
| `NumberStyle` | Blue-300 |
| `RepoStyle` | Blue, bold |

Colours respect `NO_COLOR` and `TERM=dumb`. Use `cli.ColorEnabled()` to check and `cli.SetColorEnabled(false)` to disable programmatically.

## Formatting Utilities

```go
cli.Truncate("long string", 10)  // "long st..."
cli.Pad("short", 20)             // "short               "
cli.FormatAge(time.Now().Add(-2*time.Hour))  // "2h ago"
```

## Logging

The framework provides package-level log functions that delegate to the Core log service when available:

```go
cli.LogDebug("cache miss", "key", "foo")
cli.LogInfo("server started", "port", 8080)
cli.LogWarn("slow query", "duration", "3.2s")
cli.LogError("connection failed", "err", err)
cli.LogSecurity("login attempt", "user", "admin")
```
