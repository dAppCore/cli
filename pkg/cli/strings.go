package cli

import "fmt"

// Sprintf formats a string using a format template.
//
//	msg := cli.Sprintf("Hello, %s! You have %d messages.", name, count)
func Sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// Sprint formats using default formats without a format string.
//
//	label := cli.Sprint("count:", count)
func Sprint(args ...any) string {
	return fmt.Sprint(args...)
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
		return compileGlyphs(fmt.Sprintf(format, args...))
	}
	return style.Render(compileGlyphs(fmt.Sprintf(format, args...)))
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
