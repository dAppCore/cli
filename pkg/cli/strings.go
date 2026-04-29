package cli

import "dappco.re/go"

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
//
//	cli.Repeat("-", 10)  // "----------"
func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	b := core.NewBuilder()
	for range count {
		b.WriteString(s)
	}
	return b.String()
}

// LastIndex returns the index of the last instance of substr in s,
// or -1 if substr is not present.
//
//	cli.LastIndex("hello\nworld", "\n")  // 5
func LastIndex(s, substr string) int {
	if substr == "" {
		return len(s)
	}
	last := -1
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			last = i
		}
	}
	return last
}

// Atoi parses a decimal integer from s.
//
//	r := cli.Atoi("42")
func Atoi(s string) core.Result {
	return core.Atoi(s)
}

// ParseHexByte parses a 2-character hex string into a byte value (0-255).
//
//	r := cli.ParseHexByte("ff")
func ParseHexByte(s string) core.Result {
	r := core.ParseInt(s, 16, 16)
	if !r.OK {
		return r
	}
	value := r.Value.(int64)
	if value > 255 {
		return core.Fail(core.NewError("hex byte out of range"))
	}
	return core.Ok(int(value))
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
