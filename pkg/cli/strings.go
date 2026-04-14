package cli

import (
	"strconv"
	"strings"

	"dappco.re/go/core"
)

// Sprintf formats a string using a format template.
//
//	msg := cli.Sprintf("Hello, %s! You have %d messages.", name, count)
func Sprintf(format string, args ...any) string {
	return core.Sprintf(format, args...)
}

// Sprint formats using default formats without a format string.
//
//	label := cli.Sprint("count:", count)
func Sprint(args ...any) string {
	return core.Sprint(args...)
}

// Repeat returns a new string consisting of count copies of s.
// Wraps strings.Repeat so consumer files do not import "strings" directly.
//
//	cli.Repeat("-", 10)  // "----------"
func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

// LastIndex returns the index of the last instance of substr in s,
// or -1 if substr is not present.
// Wraps strings.LastIndex so consumer files do not import "strings" directly.
//
//	cli.LastIndex("hello\nworld", "\n")  // 5
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

// Atoi parses a decimal integer from s.
// Wraps strconv.Atoi so consumer files do not import "strconv" directly.
//
//	n, err := cli.Atoi("42")  // 42, nil
func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseHexByte parses a 2-character hex string into a byte value (0-255).
// Wraps strconv.ParseUint so consumer files do not import "strconv" directly.
//
//	r, _ := cli.ParseHexByte("ff")  // 255, nil
func ParseHexByte(s string) (int, error) {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// Styled returns text with a style applied.
//
//	label := cli.Styled(cli.AccentStyle, "core dev")
func Styled(style *AnsiStyle, text string) string {
	if style == nil {
		return compileGlyphs(text)
	}
	return style.Render(compileGlyphs(text))
}

// Styledf returns formatted text with a style applied.
//
//	header := cli.Styledf(cli.HeaderStyle, "%s v%s", name, version)
func Styledf(style *AnsiStyle, format string, args ...any) string {
	if style == nil {
		return compileGlyphs(core.Sprintf(format, args...))
	}
	return style.Render(compileGlyphs(core.Sprintf(format, args...)))
}

// SuccessStr returns a success-styled string without printing it.
//
//	line := cli.SuccessStr("all tests passed")
func SuccessStr(msg string) string {
	return SuccessStyle.Render(Glyph(":check:") + " " + compileGlyphs(msg))
}

// ErrorStr returns an error-styled string without printing it.
//
//	line := cli.ErrorStr("connection refused")
func ErrorStr(msg string) string {
	return ErrorStyle.Render(Glyph(":cross:") + " " + compileGlyphs(msg))
}

// WarnStr returns a warning-styled string without printing it.
//
//	line := cli.WarnStr("deprecated flag")
func WarnStr(msg string) string {
	return WarningStyle.Render(Glyph(":warn:") + " " + compileGlyphs(msg))
}

// InfoStr returns an info-styled string without printing it.
//
//	line := cli.InfoStr("listening on :8080")
func InfoStr(msg string) string {
	return InfoStyle.Render(Glyph(":info:") + " " + compileGlyphs(msg))
}

// DimStr returns a dim-styled string without printing it.
//
//	line := cli.DimStr("optional: use --verbose for details")
func DimStr(msg string) string {
	return DimStyle.Render(compileGlyphs(msg))
}
