---
title: Interactive Prompts
description: Text prompts, confirmations, single and multi-select menus, and type-safe generic selectors.
---

# Interactive Prompts

The framework provides several prompt functions, ranging from simple text input to type-safe generic selectors.

## Text Prompt

Basic text input with an optional default:

```go
name, err := cli.Prompt("Project name", "my-app")
// Project name [my-app]: _
// Returns default if user presses Enter
```

## Question (Enhanced Prompt)

`Question` extends `Prompt` with validation and required-input support:

```go
name := cli.Question("Enter your name:")
name := cli.Question("Enter your name:", cli.WithDefault("Anonymous"))
name := cli.Question("Enter your name:", cli.RequiredInput())
```

With validation:

```go
port := cli.Question("Port:", cli.WithValidator(func(s string) error {
    n, err := strconv.Atoi(s)
    if err != nil || n < 1 || n > 65535 {
        return fmt.Errorf("must be 1-65535")
    }
    return nil
}))
```

Grammar-composed question:

```go
name := cli.QuestionAction("rename", "old.txt")
// Rename old.txt? _
```

## Confirmation

Yes/no confirmation with sensible defaults:

```go
if cli.Confirm("Delete file?") { ... }             // Default: no  [y/N]
if cli.Confirm("Save?", cli.DefaultYes()) { ... }   // Default: yes [Y/n]
if cli.Confirm("Destroy?", cli.Required()) { ... }  // Must type y or n [y/n]
```

With auto-timeout:

```go
if cli.Confirm("Continue?", cli.Timeout(30*time.Second)) { ... }
// Continue? [y/N] (auto in 30s)
// Auto-selects default after timeout
```

Combine options:

```go
if cli.Confirm("Deploy?", cli.DefaultYes(), cli.Timeout(10*time.Second)) { ... }
// Deploy? [Y/n] (auto in 10s)
```

### Grammar-Composed Confirmation

```go
if cli.ConfirmAction("delete", "config.yaml") { ... }
// Delete config.yaml? [y/N]

if cli.ConfirmAction("save", "changes", cli.DefaultYes()) { ... }
// Save changes? [Y/n]
```

### Dangerous Action (Double Confirmation)

```go
if cli.ConfirmDangerousAction("delete", "production database") { ... }
// Delete production database? [y/n]    (must type y/n)
// Really delete production database? [y/n]
```

## Single Select

Numbered menu, returns the selected string:

```go
choice, err := cli.Select("Choose backend:", []string{"metal", "rocm", "cpu"})
// Choose backend:
//   1. metal
//   2. rocm
//   3. cpu
// Choose [1-3]: _
```

## Multi Select

Space-separated number input, returns selected strings:

```go
tags, err := cli.MultiSelect("Enable features:", []string{"auth", "api", "admin", "mcp"})
// Enable features:
//   1. auth
//   2. api
//   3. admin
//   4. mcp
// Choose (space-separated) [1-4]: _
```

## Type-Safe Generic Select (`Choose`)

For selecting from typed slices with custom display:

```go
type File struct {
    Name string
    Size int64
}

files := []File{{Name: "a.go", Size: 1024}, {Name: "b.go", Size: 2048}}

choice := cli.Choose("Select a file:", files,
    cli.Display(func(f File) string {
        return fmt.Sprintf("%s (%d bytes)", f.Name, f.Size)
    }),
)
```

Enable `cli.Filter()` to let users type a substring and narrow the visible choices before selecting a number:

```go
choice := cli.Choose("Select:", items, cli.Filter[Item]())
```

With a default selection:

```go
choice := cli.Choose("Select:", items, cli.WithDefaultIndex[Item](0))
// Items marked with * are the default when Enter is pressed
```

Grammar-composed:

```go
file := cli.ChooseAction("select", "file", files)
// Select file:
//   1. ...
```

## Type-Safe Generic Multi-Select (`ChooseMulti`)

Select multiple items with ranges:

```go
selected := cli.ChooseMulti("Select files:", files,
    cli.Display(func(f File) string { return f.Name }),
)
// Select files:
//   1. a.go
//   2. b.go
//   3. c.go
// Enter numbers (e.g., 1 3 5 or 1-3) or empty for none: _
```

Input formats:
- `1 3 5` -- select items 1, 3, and 5
- `1-3` -- select items 1, 2, and 3
- `1 3-5` -- select items 1, 3, 4, and 5
- (empty) -- select none

Grammar-composed:

```go
files := cli.ChooseMultiAction("select", "files", allFiles)
```

## Testing Prompts

Override stdin for testing:

```go
cli.SetStdin(strings.NewReader("test input\n"))
defer cli.SetStdin(os.Stdin)
```
