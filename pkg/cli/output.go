package cli

import (
	"io"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/i18n"
)

// Blank prints an empty line.
func Blank() {
	core.Print(stdoutWriter(), "")
}

// Echo translates a key via i18n.T and prints with newline.
// No automatic styling - use Success/Error/Warn/Info for styled output.
func Echo(key string, args ...any) {
	core.Print(stdoutWriter(), "%s", compileGlyphs(i18n.T(key, args...)))
}

// Print outputs formatted text (no newline).
// Glyph shortcodes like :check: are converted.
func Print(format string, args ...any) {
	io.WriteString(stdoutWriter(), compileGlyphs(core.Sprintf(format, args...)))
}

// Println outputs formatted text with newline.
// Glyph shortcodes like :check: are converted.
func Println(format string, args ...any) {
	core.Print(stdoutWriter(), "%s", compileGlyphs(core.Sprintf(format, args...)))
}

// Text prints arguments space-separated with a trailing newline, handling glyphs.
//
//	cli.Text("count:", count)
func Text(args ...any) {
	core.Print(stdoutWriter(), "%s", compileGlyphs(core.Sprint(args...)))
}

// Success prints a success message with checkmark (green).
func Success(msg string) {
	core.Print(stdoutWriter(), "%s", SuccessStyle.Render(Glyph(":check:")+" "+compileGlyphs(msg)))
}

// Successf prints a formatted success message.
func Successf(format string, args ...any) {
	Success(core.Sprintf(format, args...))
}

// Error prints an error message with cross (red) to stderr and logs it.
func Error(msg string) {
	LogError(msg)
	core.Print(stderrWriter(), "%s", ErrorStyle.Render(Glyph(":cross:")+" "+compileGlyphs(msg)))
}

// Errorf prints a formatted error message to stderr and logs it.
func Errorf(format string, args ...any) {
	Error(core.Sprintf(format, args...))
}

// ErrorWrap prints a wrapped error message to stderr and logs it.
func ErrorWrap(err error, msg string) {
	if err == nil {
		return
	}
	Error(core.Sprintf("%s: %v", msg, err))
}

// ErrorWrapVerb prints a wrapped error using i18n grammar to stderr and logs it.
func ErrorWrapVerb(err error, verb, subject string) {
	if err == nil {
		return
	}
	msg := i18n.ActionFailed(verb, subject)
	Error(core.Sprintf("%s: %v", msg, err))
}

// ErrorWrapAction prints a wrapped error using i18n grammar to stderr and logs it.
func ErrorWrapAction(err error, verb string) {
	if err == nil {
		return
	}
	msg := i18n.ActionFailed(verb, "")
	Error(core.Sprintf("%s: %v", msg, err))
}

// Warn prints a warning message with warning symbol (amber) to stderr and logs it.
func Warn(msg string) {
	LogWarn(msg)
	core.Print(stderrWriter(), "%s", WarningStyle.Render(Glyph(":warn:")+" "+compileGlyphs(msg)))
}

// Warnf prints a formatted warning message to stderr and logs it.
func Warnf(format string, args ...any) {
	Warn(core.Sprintf(format, args...))
}

// Info prints an info message with info symbol (blue).
func Info(msg string) {
	core.Print(stdoutWriter(), "%s", InfoStyle.Render(Glyph(":info:")+" "+compileGlyphs(msg)))
}

// Infof prints a formatted info message.
func Infof(format string, args ...any) {
	Info(core.Sprintf(format, args...))
}

// Dim prints dimmed text.
func Dim(msg string) {
	core.Print(stdoutWriter(), "%s", DimStyle.Render(compileGlyphs(msg)))
}

// Progress prints a progress indicator that overwrites the current line.
// Uses i18n.Progress for gerund form ("Checking...").
func Progress(verb string, current, total int, item ...string) {
	msg := compileGlyphs(i18n.Progress(verb))
	if len(item) > 0 && item[0] != "" {
		io.WriteString(stderrWriter(), core.Sprintf("\033[2K\r%s %d/%d %s", DimStyle.Render(msg), current, total, compileGlyphs(item[0])))
	} else {
		io.WriteString(stderrWriter(), core.Sprintf("\033[2K\r%s %d/%d", DimStyle.Render(msg), current, total))
	}
}

// ProgressDone clears the progress line.
func ProgressDone() {
	io.WriteString(stderrWriter(), "\033[2K\r")
}

// Label prints a "Label: value" line.
func Label(word, value string) {
	core.Print(stdoutWriter(), "%s %s", KeyStyle.Render(compileGlyphs(i18n.Label(word))), compileGlyphs(value))
}

// Task prints a task header: "[label] message"
//
//	cli.Task("php", "Running tests...")  // [php] Running tests...
//	cli.Task("go", i18n.Progress("build"))  // [go] Building...
func Task(label, message string) {
	core.Print(stdoutWriter(), "%s %s\n", DimStyle.Render("["+compileGlyphs(label)+"]"), compileGlyphs(message))
}

// Section prints a section header: "── SECTION ──"
//
//	cli.Section("audit")  // ── AUDIT ──
func Section(name string) {
	dash := Glyph(":dash:")
	header := dash + dash + " " + core.Upper(compileGlyphs(name)) + " " + dash + dash
	core.Print(stdoutWriter(), "%s", AccentStyle.Render(header))
}

// Hint prints a labelled hint: "label: message"
//
//	cli.Hint("install", "composer require vimeo/psalm")
//	cli.Hint("fix", "core php fmt --fix")
func Hint(label, message string) {
	core.Print(stdoutWriter(), "  %s %s", DimStyle.Render(compileGlyphs(label)+":"), compileGlyphs(message))
}

// Severity prints a severity-styled message.
//
//	cli.Severity("critical", "SQL injection")  // red, bold
//	cli.Severity("high", "XSS vulnerability")  // orange
//	cli.Severity("medium", "Missing CSRF")     // amber
//	cli.Severity("low", "Debug enabled")       // gray
func Severity(level, message string) {
	var style *AnsiStyle
	switch core.Lower(level) {
	case "critical":
		style = NewStyle().Bold().Foreground(ColourRed500)
	case "high":
		style = NewStyle().Bold().Foreground(ColourOrange500)
	case "medium":
		style = NewStyle().Foreground(ColourAmber500)
	case "low":
		style = NewStyle().Foreground(ColourGray500)
	default:
		style = DimStyle
	}
	core.Print(stdoutWriter(), "  %s %s", style.Render("["+compileGlyphs(level)+"]"), compileGlyphs(message))
}

// Result prints a result line: "✓ message" or "✗ message"
//
//	cli.Result(passed, "All tests passed")
//	cli.Result(false, "3 tests failed")
func Result(passed bool, message string) {
	if passed {
		Success(message)
	} else {
		Error(message)
	}
}
