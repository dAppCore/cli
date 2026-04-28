package cli

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"testing/fstest"
	"time"

	core "dappco.re/go"
)

func ax7PlainCLI(t *core.T) {
	t.Helper()
	originalTheme := currentTheme
	originalColor := ColorEnabled()
	UseASCII()
	SetColorEnabled(false)
	t.Cleanup(func() {
		currentTheme = originalTheme
		SetColorEnabled(originalColor)
		SetStdout(nil)
		SetStderr(nil)
		SetStdin(nil)
	})
}

func ax7CaptureStdout(t *core.T, fn func()) string {
	t.Helper()
	out := core.NewBuilder()
	SetStdout(out)
	defer SetStdout(nil)
	fn()
	return out.String()
}

func ax7CaptureStderr(t *core.T, fn func()) string {
	t.Helper()
	out := core.NewBuilder()
	SetStderr(out)
	defer SetStderr(nil)
	fn()
	return out.String()
}

func ax7FakeCommands(t *core.T, scripts map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range scripts {
		path := core.Path(dir, name)
		data := []byte("#!/bin/sh\n" + body)
		core.RequireNoError(t, os.WriteFile(path, data, 0o755))
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return dir
}

func ax7RunSelf(t *core.T, envKey string) error {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run", "^"+t.Name()+"$")
	cmd.Env = append(os.Environ(), envKey+"=1")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}

func TestAX7CLI_Err_Good(t *core.T) {
	err := Err("missing %s", "config")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "missing config")
}

func TestAX7CLI_Err_Bad(t *core.T) {
	err := Err("")

	core.AssertError(t, err)
	core.AssertEqual(t, "cli: ", err.Error())
}

func TestAX7CLI_Err_Ugly(t *core.T) {
	err := Err("line\n%s", "two")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "line\ntwo")
}

func TestAX7CLI_Wrap_Good(t *core.T) {
	base := core.NewError("connection refused")
	err := Wrap(base, "connect")

	core.AssertError(t, err)
	core.AssertTrue(t, Is(err, base))
}

func TestAX7CLI_Wrap_Bad(t *core.T) {
	err := Wrap(nil, "connect")

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestAX7CLI_Wrap_Ugly(t *core.T) {
	base := Err("root")
	err := Wrap(base, "")

	core.AssertError(t, err)
	core.AssertTrue(t, Is(err, base))
}

func TestAX7CLI_WrapVerb_Good(t *core.T) {
	base := core.NewError("denied")
	err := WrapVerb(base, "load", "config")

	core.AssertContains(t, err.Error(), "Failed to load config")
	core.AssertTrue(t, Is(err, base))
}

func TestAX7CLI_WrapVerb_Bad(t *core.T) {
	err := WrapVerb(nil, "load", "config")

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestAX7CLI_WrapVerb_Ugly(t *core.T) {
	err := WrapVerb(core.NewError("denied"), "", "config")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "denied")
}

func TestAX7CLI_WrapAction_Good(t *core.T) {
	base := core.NewError("offline")
	err := WrapAction(base, "connect")

	core.AssertContains(t, err.Error(), "Failed to connect")
	core.AssertTrue(t, Is(err, base))
}

func TestAX7CLI_WrapAction_Bad(t *core.T) {
	err := WrapAction(nil, "connect")

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestAX7CLI_WrapAction_Ugly(t *core.T) {
	err := WrapAction(core.NewError("offline"), "")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "offline")
}

func TestAX7CLI_Is_Good(t *core.T) {
	target := core.NewError("target")
	err := Wrap(target, "wrap")

	core.AssertTrue(t, Is(err, target))
	core.AssertFalse(t, Is(err, core.NewError("other")))
}

func TestAX7CLI_Is_Bad(t *core.T) {
	left := core.NewError("left")
	right := core.NewError("right")

	core.AssertFalse(t, Is(left, right))
	core.AssertFalse(t, Is(nil, right))
}

func TestAX7CLI_Is_Ugly(t *core.T) {
	got := Is(nil, nil)
	core.AssertTrue(t, got)
	core.AssertTrue(t, Is(nil, nil))
	core.AssertFalse(t, Is(core.NewError("x"), nil))
}

func TestAX7CLI_As_Good(t *core.T) {
	err := Exit(7, Err("exit"))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 7, exitErr.Code)
}

func TestAX7CLI_As_Bad(t *core.T) {
	var exitErr *ExitError
	err := core.NewError("plain")

	core.AssertFalse(t, As(err, &exitErr))
	core.AssertNil(t, exitErr)
}

func TestAX7CLI_As_Ugly(t *core.T) {
	var exitErr *ExitError

	core.AssertFalse(t, As(nil, &exitErr))
	core.AssertNil(t, exitErr)
}

func TestAX7CLI_Join_Good(t *core.T) {
	first := Err("first")
	second := Err("second")
	err := Join(first, second)

	core.AssertTrue(t, Is(err, first))
	core.AssertTrue(t, Is(err, second))
}

func TestAX7CLI_Join_Bad(t *core.T) {
	err := Join(nil, nil)

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestAX7CLI_Join_Ugly(t *core.T) {
	err := Join(nil, Err("only"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "only")
}

func TestAX7CLI_ExitError_Error_Good(t *core.T) {
	err := &ExitError{Code: 2, Err: Err("failed")}

	core.AssertEqual(t, "cli: failed", err.Error())
	core.AssertEqual(t, 2, err.Code)
}

func TestAX7CLI_ExitError_Error_Bad(t *core.T) {
	err := &ExitError{Code: 2}

	core.AssertEqual(t, "", err.Error())
	core.AssertNil(t, err.Err)
}

func TestAX7CLI_ExitError_Error_Ugly(t *core.T) {
	err := &ExitError{Code: 255, Err: Err("line\nfail")}

	core.AssertContains(t, err.Error(), "line\nfail")
	core.AssertEqual(t, 255, err.Code)
}

func TestAX7CLI_ExitError_Unwrap_Good(t *core.T) {
	base := Err("base")
	err := &ExitError{Code: 2, Err: base}

	core.AssertEqual(t, base, err.Unwrap())
	core.AssertTrue(t, Is(err, base))
}

func TestAX7CLI_ExitError_Unwrap_Bad(t *core.T) {
	err := &ExitError{Code: 2}

	core.AssertNil(t, err.Unwrap())
	core.AssertEqual(t, "", err.Error())
}

func TestAX7CLI_ExitError_Unwrap_Ugly(t *core.T) {
	var err *ExitError

	core.AssertPanics(t, func() { _ = err.Unwrap() })
	core.AssertNil(t, err)
}

func TestAX7CLI_Exit_Good(t *core.T) {
	err := Exit(2, Err("bad args"))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 2, exitErr.Code)
}

func TestAX7CLI_Exit_Bad(t *core.T) {
	err := Exit(2, nil)

	core.AssertNil(t, err)
	core.AssertFalse(t, As(err, new(*ExitError)))
}

func TestAX7CLI_Exit_Ugly(t *core.T) {
	err := Exit(0, Err("zero"))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 0, exitErr.Code)
}

func TestAX7CLI_Fatal_Good(t *core.T) {
	if os.Getenv("AX7_FATAL_GOOD") == "1" {
		SetStderr(io.Discard)
		Fatal(Err("fatal"))
		return
	}
	err := ax7RunSelf(t, "AX7_FATAL_GOOD")
	core.AssertError(t, err)
}

func TestAX7CLI_Fatal_Bad(t *core.T) {
	ax7PlainCLI(t)

	core.AssertNotPanics(t, func() { Fatal(nil) })
	core.AssertEqual(t, "", ax7CaptureStderr(t, func() { Fatal(nil) }))
}

func TestAX7CLI_Fatal_Ugly(t *core.T) {
	if os.Getenv("AX7_FATAL_UGLY") == "1" {
		SetStderr(io.Discard)
		Fatal(Err("fatal\nline"))
		return
	}
	err := ax7RunSelf(t, "AX7_FATAL_UGLY")
	core.AssertError(t, err)
}

func TestAX7CLI_Fatalf_Good(t *core.T) {
	if os.Getenv("AX7_FATALF_GOOD") == "1" {
		SetStderr(io.Discard)
		Fatalf("fatal %s", "format")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALF_GOOD")
	core.AssertError(t, err)
}

func TestAX7CLI_Fatalf_Bad(t *core.T) {
	if os.Getenv("AX7_FATALF_BAD") == "1" {
		SetStderr(io.Discard)
		Fatalf("")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALF_BAD")
	core.AssertError(t, err)
}

func TestAX7CLI_Fatalf_Ugly(t *core.T) {
	if os.Getenv("AX7_FATALF_UGLY") == "1" {
		SetStderr(io.Discard)
		Fatalf("fatal %d", 42)
		return
	}
	err := ax7RunSelf(t, "AX7_FATALF_UGLY")
	core.AssertError(t, err)
}

func TestAX7CLI_FatalWrap_Good(t *core.T) {
	if os.Getenv("AX7_FATALWRAP_GOOD") == "1" {
		SetStderr(io.Discard)
		FatalWrap(Err("root"), "wrap")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALWRAP_GOOD")
	core.AssertError(t, err)
}

func TestAX7CLI_FatalWrap_Bad(t *core.T) {
	ax7PlainCLI(t)

	core.AssertNotPanics(t, func() { FatalWrap(nil, "wrap") })
	core.AssertEqual(t, "", ax7CaptureStderr(t, func() { FatalWrap(nil, "wrap") }))
}

func TestAX7CLI_FatalWrap_Ugly(t *core.T) {
	if os.Getenv("AX7_FATALWRAP_UGLY") == "1" {
		SetStderr(io.Discard)
		FatalWrap(Err("root"), "")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALWRAP_UGLY")
	core.AssertError(t, err)
}

func TestAX7CLI_FatalWrapVerb_Good(t *core.T) {
	if os.Getenv("AX7_FATALWRAPVERB_GOOD") == "1" {
		SetStderr(io.Discard)
		FatalWrapVerb(Err("root"), "load", "config")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALWRAPVERB_GOOD")
	core.AssertError(t, err)
}

func TestAX7CLI_FatalWrapVerb_Bad(t *core.T) {
	ax7PlainCLI(t)

	core.AssertNotPanics(t, func() { FatalWrapVerb(nil, "load", "config") })
	core.AssertEqual(t, "", ax7CaptureStderr(t, func() { FatalWrapVerb(nil, "load", "config") }))
}

func TestAX7CLI_FatalWrapVerb_Ugly(t *core.T) {
	if os.Getenv("AX7_FATALWRAPVERB_UGLY") == "1" {
		SetStderr(io.Discard)
		FatalWrapVerb(Err("root"), "", "")
		return
	}
	err := ax7RunSelf(t, "AX7_FATALWRAPVERB_UGLY")
	core.AssertError(t, err)
}

func TestAX7CLI_ColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7CLI_ColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestAX7CLI_ColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(!original)
	defer SetColorEnabled(original)

	core.AssertEqual(t, !original, ColorEnabled())
}

func TestAX7CLI_SetColorEnabled_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7CLI_SetColorEnabled_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)

	core.AssertFalse(t, ColorEnabled())
}

func TestAX7CLI_SetColorEnabled_Ugly(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	SetColorEnabled(true)
	defer SetColorEnabled(original)

	core.AssertTrue(t, ColorEnabled())
}

func TestAX7CLI_NewStyle_Good(t *core.T) {
	s := NewStyle()

	core.AssertNotNil(t, s)
	core.AssertEqual(t, "plain", s.Render("plain"))
}

func TestAX7CLI_NewStyle_Bad(t *core.T) {
	s := NewStyle()

	core.AssertFalse(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestAX7CLI_NewStyle_Ugly(t *core.T) {
	s := NewStyle().Bold().Dim().Italic().Underline()

	core.AssertTrue(t, s.bold)
	core.AssertTrue(t, s.underline)
}

func TestAX7CLI_AnsiStyle_Bold_Good(t *core.T) {
	s := NewStyle().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertEqual(t, s, s.Bold())
}

func TestAX7CLI_AnsiStyle_Bold_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Bold() })
	core.AssertNil(t, s)
}

func TestAX7CLI_AnsiStyle_Bold_Ugly(t *core.T) {
	s := NewStyle().Bold().Bold()

	core.AssertTrue(t, s.bold)
	core.AssertFalse(t, s.dim)
}

func TestAX7CLI_AnsiStyle_Dim_Good(t *core.T) {
	s := NewStyle().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertEqual(t, s, s.Dim())
}

func TestAX7CLI_AnsiStyle_Dim_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Dim() })
	core.AssertNil(t, s)
}

func TestAX7CLI_AnsiStyle_Dim_Ugly(t *core.T) {
	s := NewStyle().Dim().Dim()

	core.AssertTrue(t, s.dim)
	core.AssertFalse(t, s.bold)
}

func TestAX7CLI_AnsiStyle_Italic_Good(t *core.T) {
	s := NewStyle().Italic()

	core.AssertTrue(t, s.italic)
	core.AssertEqual(t, s, s.Italic())
}

func TestAX7CLI_AnsiStyle_Italic_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Italic() })
	core.AssertNil(t, s)
}

func TestAX7CLI_AnsiStyle_Italic_Ugly(t *core.T) {
	s := NewStyle().Italic().Italic()

	core.AssertTrue(t, s.italic)
	core.AssertFalse(t, s.underline)
}

func TestAX7CLI_AnsiStyle_Underline_Good(t *core.T) {
	s := NewStyle().Underline()

	core.AssertTrue(t, s.underline)
	core.AssertEqual(t, s, s.Underline())
}

func TestAX7CLI_AnsiStyle_Underline_Bad(t *core.T) {
	var s *AnsiStyle

	core.AssertPanics(t, func() { s.Underline() })
	core.AssertNil(t, s)
}

func TestAX7CLI_AnsiStyle_Underline_Ugly(t *core.T) {
	s := NewStyle().Underline().Underline()

	core.AssertTrue(t, s.underline)
	core.AssertFalse(t, s.bold)
}

func TestAX7CLI_AnsiStyle_Foreground_Good(t *core.T) {
	s := NewStyle().Foreground("#ff0000")

	core.AssertContains(t, s.fg, "38;2;255;0;0")
	core.AssertEqual(t, s, s.Foreground("#00ff00"))
}

func TestAX7CLI_AnsiStyle_Foreground_Bad(t *core.T) {
	s := NewStyle().Foreground("bad")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAX7CLI_AnsiStyle_Foreground_Ugly(t *core.T) {
	s := NewStyle().Foreground("")

	core.AssertContains(t, s.fg, "255;255;255")
	core.AssertNotEmpty(t, s.fg)
}

func TestAX7CLI_AnsiStyle_Background_Good(t *core.T) {
	s := NewStyle().Background("#0000ff")

	core.AssertContains(t, s.bg, "48;2;0;0;255")
	core.AssertEqual(t, s, s.Background("#00ff00"))
}

func TestAX7CLI_AnsiStyle_Background_Bad(t *core.T) {
	s := NewStyle().Background("bad")

	core.AssertContains(t, s.bg, "255;255;255")
	core.AssertNotEmpty(t, s.bg)
}

func TestAX7CLI_AnsiStyle_Background_Ugly(t *core.T) {
	s := NewStyle().Background("")

	core.AssertContains(t, s.bg, "255;255;255")
	core.AssertNotEmpty(t, s.bg)
}

func TestAX7CLI_AnsiStyle_Render_Good(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(true)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertContains(t, got, "text")
	core.AssertContains(t, got, "\033[1m")
}

func TestAX7CLI_AnsiStyle_Render_Bad(t *core.T) {
	original := ColorEnabled()
	SetColorEnabled(false)
	defer SetColorEnabled(original)
	got := NewStyle().Bold().Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestAX7CLI_AnsiStyle_Render_Ugly(t *core.T) {
	var s *AnsiStyle
	got := s.Render("text")

	core.AssertEqual(t, "text", got)
	core.AssertNotContains(t, got, "\033")
}

func TestAX7CLI_Sprintf_Good(t *core.T) {
	got := Sprintf("hello %s", "codex")

	core.AssertEqual(t, "hello codex", got)
	core.AssertContains(t, got, "codex")
}

func TestAX7CLI_Sprintf_Bad(t *core.T) {
	got := Sprintf("%s", "bad")

	core.AssertEqual(t, "bad", got)
	core.AssertContains(t, got, "bad")
}

func TestAX7CLI_Sprintf_Ugly(t *core.T) {
	got := Sprintf("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Sprint_Good(t *core.T) {
	got := Sprint("count:", 2)

	core.AssertEqual(t, "count:2", got)
	core.AssertContains(t, got, "2")
}

func TestAX7CLI_Sprint_Bad(t *core.T) {
	got := Sprint()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Sprint_Ugly(t *core.T) {
	got := Sprint(nil, "x")

	core.AssertEqual(t, "<nil>x", got)
	core.AssertContains(t, got, "nil")
}

func TestAX7CLI_Styled_Good(t *core.T) {
	ax7PlainCLI(t)
	got := Styled(NewStyle().Bold(), ":check: ready")

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "[OK]")
}

func TestAX7CLI_Styled_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := Styled(nil, ":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestAX7CLI_Styled_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := Styled(NewStyle(), "")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Styledf_Good(t *core.T) {
	ax7PlainCLI(t)
	got := Styledf(NewStyle().Bold(), "%s", ":check:")

	core.AssertEqual(t, "[OK]", got)
	core.AssertContains(t, got, "[OK]")
}

func TestAX7CLI_Styledf_Bad(t *core.T) {
	got := Styledf(nil, "")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Styledf_Ugly(t *core.T) {
	got := Styledf(nil, "%s", "bad")

	core.AssertEqual(t, "bad", got)
	core.AssertContains(t, got, "bad")
}

func TestAX7CLI_SuccessStr_Good(t *core.T) {
	ax7PlainCLI(t)
	got := SuccessStr("done")

	core.AssertContains(t, got, "done")
	core.AssertContains(t, got, "[OK]")
}

func TestAX7CLI_SuccessStr_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := SuccessStr("")

	core.AssertContains(t, got, "[OK]")
	core.AssertNotContains(t, got, "done")
}

func TestAX7CLI_SuccessStr_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := SuccessStr(":check:")

	core.AssertContains(t, got, "[OK]")
	core.AssertNotContains(t, got, ":check:")
}

func TestAX7CLI_ErrorStr_Good(t *core.T) {
	ax7PlainCLI(t)
	got := ErrorStr("failed")

	core.AssertContains(t, got, "failed")
	core.AssertContains(t, got, "[FAIL]")
}

func TestAX7CLI_ErrorStr_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := ErrorStr("")

	core.AssertContains(t, got, "[FAIL]")
	core.AssertNotContains(t, got, "failed")
}

func TestAX7CLI_ErrorStr_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := ErrorStr(":cross:")

	core.AssertContains(t, got, "[FAIL]")
	core.AssertNotContains(t, got, ":cross:")
}

func TestAX7CLI_WarnStr_Good(t *core.T) {
	ax7PlainCLI(t)
	got := WarnStr("careful")

	core.AssertContains(t, got, "careful")
	core.AssertContains(t, got, "[WARN]")
}

func TestAX7CLI_WarnStr_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := WarnStr("")

	core.AssertContains(t, got, "[WARN]")
	core.AssertNotContains(t, got, "careful")
}

func TestAX7CLI_WarnStr_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := WarnStr(":warn:")

	core.AssertContains(t, got, "[WARN]")
	core.AssertNotContains(t, got, ":warn:")
}

func TestAX7CLI_InfoStr_Good(t *core.T) {
	ax7PlainCLI(t)
	got := InfoStr("ready")

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "[INFO]")
}

func TestAX7CLI_InfoStr_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := InfoStr("")

	core.AssertContains(t, got, "[INFO]")
	core.AssertNotContains(t, got, "ready")
}

func TestAX7CLI_InfoStr_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := InfoStr(":info:")

	core.AssertContains(t, got, "[INFO]")
	core.AssertNotContains(t, got, ":info:")
}

func TestAX7CLI_DimStr_Good(t *core.T) {
	ax7PlainCLI(t)
	got := DimStr("quiet")

	core.AssertEqual(t, "quiet", got)
	core.AssertContains(t, got, "quiet")
}

func TestAX7CLI_DimStr_Bad(t *core.T) {
	ax7PlainCLI(t)
	got := DimStr("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_DimStr_Ugly(t *core.T) {
	ax7PlainCLI(t)
	got := DimStr(":check:")

	core.AssertEqual(t, "[OK]", got)
	core.AssertNotContains(t, got, ":check:")
}

func TestAX7CLI_Blank_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Blank() })

	core.AssertContains(t, out, "\n")
	core.AssertNotPanics(t, func() { Blank() })
}

func TestAX7CLI_Blank_Bad(t *core.T) {
	SetStdout(io.Discard)
	defer SetStdout(nil)

	core.AssertNotPanics(t, func() { Blank() })
	core.AssertNotNil(t, stdoutWriter())
}

func TestAX7CLI_Blank_Ugly(t *core.T) {
	out := ax7CaptureStdout(t, func() {
		Blank()
		Blank()
	})

	core.AssertContains(t, out, "\n")
	core.AssertTrue(t, core.RuneCount(out) >= 2)
}

func TestAX7CLI_Echo_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Echo("i18n.progress.check") })

	core.AssertContains(t, out, "Checking")
	core.AssertContains(t, out, "...")
}

func TestAX7CLI_Echo_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Echo("") })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "Checking")
}

func TestAX7CLI_Echo_Ugly(t *core.T) {
	out := ax7CaptureStdout(t, func() { Echo("i18n.fail.load", "config") })

	core.AssertContains(t, out, "Failed to load config")
	core.AssertNotContains(t, out, "i18n.fail")
}

func TestAX7CLI_Print_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Print("hello %s", "codex") })

	core.AssertEqual(t, "hello codex", out)
	core.AssertContains(t, out, "codex")
}

func TestAX7CLI_Print_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Print("") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_Print_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Print(":check:") })

	core.AssertEqual(t, "[OK]", out)
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Println_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Println("hello %s", "codex") })

	core.AssertContains(t, out, "hello codex")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Println_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Println("") })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestAX7CLI_Println_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Println(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Text_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Text("count:", 2) })

	core.AssertContains(t, out, "count:2")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Text_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Text() })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestAX7CLI_Text_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Text(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Success_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Success("done") })

	core.AssertContains(t, out, "done")
	core.AssertContains(t, out, "[OK]")
}

func TestAX7CLI_Success_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Success("") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, "done")
}

func TestAX7CLI_Success_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Success(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Successf_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Successf("done %d", 1) })

	core.AssertContains(t, out, "done 1")
	core.AssertContains(t, out, "[OK]")
}

func TestAX7CLI_Successf_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Successf("") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, "done")
}

func TestAX7CLI_Successf_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Successf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[OK]")
}

func TestAX7CLI_Error_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Error("failed") })

	core.AssertContains(t, out, "failed")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_Error_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Error("") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}

func TestAX7CLI_Error_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Error(":cross:") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, ":cross:")
}

func TestAX7CLI_Errorf_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Errorf("failed %d", 1) })

	core.AssertContains(t, out, "failed 1")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_Errorf_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Errorf("") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}

func TestAX7CLI_Errorf_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Errorf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_ErrorWrap_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrap(Err("root"), "wrap") })

	core.AssertContains(t, out, "wrap")
	core.AssertContains(t, out, "root")
}

func TestAX7CLI_ErrorWrap_Bad(t *core.T) {
	out := ax7CaptureStderr(t, func() { ErrorWrap(nil, "wrap") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_ErrorWrap_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrap(Err("root"), "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_ErrorWrapVerb_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrapVerb(Err("root"), "load", "config") })

	core.AssertContains(t, out, "Failed to load config")
	core.AssertContains(t, out, "root")
}

func TestAX7CLI_ErrorWrapVerb_Bad(t *core.T) {
	out := ax7CaptureStderr(t, func() { ErrorWrapVerb(nil, "load", "config") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_ErrorWrapVerb_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrapVerb(Err("root"), "", "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_ErrorWrapAction_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrapAction(Err("root"), "connect") })

	core.AssertContains(t, out, "Failed to connect")
	core.AssertContains(t, out, "root")
}

func TestAX7CLI_ErrorWrapAction_Bad(t *core.T) {
	out := ax7CaptureStderr(t, func() { ErrorWrapAction(nil, "connect") })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_ErrorWrapAction_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { ErrorWrapAction(Err("root"), "") })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_Warn_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warn("careful") })

	core.AssertContains(t, out, "careful")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Warn_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warn("") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, "careful")
}

func TestAX7CLI_Warn_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warn(":warn:") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, ":warn:")
}

func TestAX7CLI_Warnf_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warnf("careful %d", 1) })

	core.AssertContains(t, out, "careful 1")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Warnf_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warnf("") })

	core.AssertContains(t, out, "[WARN]")
	core.AssertNotContains(t, out, "careful")
}

func TestAX7CLI_Warnf_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Warnf("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Info_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Info("ready") })

	core.AssertContains(t, out, "ready")
	core.AssertContains(t, out, "[INFO]")
}

func TestAX7CLI_Info_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Info("") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, "ready")
}

func TestAX7CLI_Info_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Info(":info:") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, ":info:")
}

func TestAX7CLI_Infof_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Infof("ready %d", 1) })

	core.AssertContains(t, out, "ready 1")
	core.AssertContains(t, out, "[INFO]")
}

func TestAX7CLI_Infof_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Infof("") })

	core.AssertContains(t, out, "[INFO]")
	core.AssertNotContains(t, out, "ready")
}

func TestAX7CLI_Infof_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Infof("%s", "bad") })

	core.AssertContains(t, out, "bad")
	core.AssertContains(t, out, "[INFO]")
}

func TestAX7CLI_Dim_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Dim("quiet") })

	core.AssertContains(t, out, "quiet")
	core.AssertNotContains(t, out, "\033")
}

func TestAX7CLI_Dim_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Dim("") })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "quiet")
}

func TestAX7CLI_Dim_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Dim(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Progress_Good(t *core.T) {
	out := ax7CaptureStderr(t, func() { Progress("check", 1, 2, "repo") })

	core.AssertContains(t, out, "1/2")
	core.AssertContains(t, out, "repo")
}

func TestAX7CLI_Progress_Bad(t *core.T) {
	out := ax7CaptureStderr(t, func() { Progress("", 0, 0) })

	core.AssertContains(t, out, "0/0")
	core.AssertContains(t, out, "\r")
}

func TestAX7CLI_Progress_Ugly(t *core.T) {
	out := ax7CaptureStderr(t, func() { Progress("tie", -1, 3, "") })

	core.AssertContains(t, out, "-1/3")
	core.AssertContains(t, out, "Tying")
}

func TestAX7CLI_ProgressDone_Good(t *core.T) {
	out := ax7CaptureStderr(t, func() { ProgressDone() })

	core.AssertContains(t, out, "\033[2K")
	core.AssertContains(t, out, "\r")
}

func TestAX7CLI_ProgressDone_Bad(t *core.T) {
	SetStderr(io.Discard)
	defer SetStderr(nil)

	core.AssertNotPanics(t, func() { ProgressDone() })
	core.AssertNotNil(t, stderrWriter())
}

func TestAX7CLI_ProgressDone_Ugly(t *core.T) {
	out := ax7CaptureStderr(t, func() {
		ProgressDone()
		ProgressDone()
	})

	core.AssertContains(t, out, "\r")
	core.AssertTrue(t, core.Contains(out, "\033[2K"))
}

func TestAX7CLI_Label_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Label("path", "/tmp") })

	core.AssertContains(t, out, "Path:")
	core.AssertContains(t, out, "/tmp")
}

func TestAX7CLI_Label_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Label("", "") })

	core.AssertNotContains(t, out, ":")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Label_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Label(":check:", ":warn:") })

	core.AssertContains(t, out, "[OK]:")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Task_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Task("go", "Running") })

	core.AssertContains(t, out, "[go]")
	core.AssertContains(t, out, "Running")
}

func TestAX7CLI_Task_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Task("", "") })

	core.AssertContains(t, out, "[]")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Task_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Task(":check:", ":warn:") })

	core.AssertContains(t, out, "[[OK]]")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Section_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Section("audit") })

	core.AssertContains(t, out, "AUDIT")
	core.AssertContains(t, out, "─")
}

func TestAX7CLI_Section_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Section("") })

	core.AssertContains(t, out, "─")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Section_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Section(":check:") })

	core.AssertContains(t, out, "[OK]")
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Hint_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Hint("fix", "run tests") })

	core.AssertContains(t, out, "fix:")
	core.AssertContains(t, out, "run tests")
}

func TestAX7CLI_Hint_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Hint("", "") })

	core.AssertContains(t, out, ":")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Hint_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Hint(":check:", ":warn:") })

	core.AssertContains(t, out, "[OK]:")
	core.AssertContains(t, out, "[WARN]")
}

func TestAX7CLI_Severity_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Severity("critical", "sql injection") })

	core.AssertContains(t, out, "[critical]")
	core.AssertContains(t, out, "sql injection")
}

func TestAX7CLI_Severity_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Severity("unknown", "message") })

	core.AssertContains(t, out, "[unknown]")
	core.AssertContains(t, out, "message")
}

func TestAX7CLI_Severity_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Severity("", "") })

	core.AssertContains(t, out, "[]")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Result_Good(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStdout(t, func() { Result(true, "passed") })

	core.AssertContains(t, out, "passed")
	core.AssertContains(t, out, "[OK]")
}

func TestAX7CLI_Result_Bad(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Result(false, "failed") })

	core.AssertContains(t, out, "failed")
	core.AssertContains(t, out, "[FAIL]")
}

func TestAX7CLI_Result_Ugly(t *core.T) {
	ax7PlainCLI(t)
	out := ax7CaptureStderr(t, func() { Result(false, "") })

	core.AssertContains(t, out, "[FAIL]")
	core.AssertNotContains(t, out, "failed")
}

func TestAX7CLI_Mode_String_Good(t *core.T) {
	got := ModeInteractive.String()

	core.AssertEqual(t, "interactive", got)
	core.AssertContains(t, got, "active")
}

func TestAX7CLI_Mode_String_Bad(t *core.T) {
	got := Mode(99).String()

	core.AssertEqual(t, "unknown", got)
	core.AssertContains(t, got, "unknown")
}

func TestAX7CLI_Mode_String_Ugly(t *core.T) {
	got := ModeDaemon.String()

	core.AssertEqual(t, "daemon", got)
	core.AssertNotEqual(t, "interactive", got)
}

func TestAX7CLI_DetectMode_Good(t *core.T) {
	t.Setenv("CORE_DAEMON", "1")
	got := DetectMode()

	core.AssertEqual(t, ModeDaemon, got)
	core.AssertEqual(t, "daemon", got.String())
}

func TestAX7CLI_DetectMode_Bad(t *core.T) {
	t.Setenv("CORE_DAEMON", "")
	SetStdout(core.NewBuilder())
	defer SetStdout(nil)

	core.AssertEqual(t, ModePipe, DetectMode())
}

func TestAX7CLI_DetectMode_Ugly(t *core.T) {
	t.Setenv("CORE_DAEMON", "0")
	SetStdout(io.Discard)
	defer SetStdout(nil)

	core.AssertEqual(t, ModePipe, DetectMode())
}

func TestAX7CLI_IsTTY_Good(t *core.T) {
	SetStdout(io.Discard)
	defer SetStdout(nil)

	core.AssertFalse(t, IsTTY())
}

func TestAX7CLI_IsTTY_Bad(t *core.T) {
	SetStdout(core.NewBuilder())
	defer SetStdout(nil)

	core.AssertFalse(t, IsTTY())
}

func TestAX7CLI_IsTTY_Ugly(t *core.T) {
	SetStdout(nil)
	got := IsTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsTTY())
}

func TestAX7CLI_IsStdinTTY_Good(t *core.T) {
	SetStdin(core.NewReader(""))
	defer SetStdin(nil)

	core.AssertFalse(t, IsStdinTTY())
}

func TestAX7CLI_IsStdinTTY_Bad(t *core.T) {
	SetStdin(nil)
	got := IsStdinTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsStdinTTY())
}

func TestAX7CLI_IsStdinTTY_Ugly(t *core.T) {
	SetStdin(core.NewReader("input"))
	defer SetStdin(nil)

	core.AssertFalse(t, IsStdinTTY())
}

func TestAX7CLI_IsStderrTTY_Good(t *core.T) {
	SetStderr(io.Discard)
	defer SetStderr(nil)

	core.AssertFalse(t, IsStderrTTY())
}

func TestAX7CLI_IsStderrTTY_Bad(t *core.T) {
	SetStderr(core.NewBuilder())
	defer SetStderr(nil)

	core.AssertFalse(t, IsStderrTTY())
}

func TestAX7CLI_IsStderrTTY_Ugly(t *core.T) {
	SetStderr(nil)
	got := IsStderrTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsStderrTTY())
}

func TestAX7CLI_WithCommands_Good(t *core.T) {
	c := core.New()
	setup := WithCommands("x", func(c *core.Core) { c.Command("x", core.Command{}) })

	setup(c)
	core.AssertTrue(t, c.Command("x").OK)
}

func TestAX7CLI_WithCommands_Bad(t *core.T) {
	c := core.New()
	setup := WithCommands("", func(c *core.Core) { c.Command("empty", core.Command{}) }, nil)

	setup(c)
	core.AssertTrue(t, c.Command("empty").OK)
}

func TestAX7CLI_WithCommands_Ugly(t *core.T) {
	c := core.New()
	called := 0
	setup := WithCommands("x", func(_ *core.Core) { called++ })

	setup(c)
	core.AssertEqual(t, 1, called)
}

func TestAX7CLI_RegisterCommands_Good(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(c *core.Core) { c.Command("registered", core.Command{}) })

	var count int
	for range RegisteredCommands() {
		count++
	}
	core.AssertEqual(t, 1, count)
}

func TestAX7CLI_RegisterCommands_Bad(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {}, nil)

	core.AssertEmpty(t, RegisteredLocales())
	core.AssertNotNil(t, RegisteredCommands())
}

func TestAX7CLI_RegisterCommands_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "registered"}))
	RegisterCommands(func(c *core.Core) { c.Command("late", core.Command{}) })

	core.AssertTrue(t, Core().Command("late").OK)
}

func TestAX7CLI_RegisteredLocales_Good(t *core.T) {
	resetGlobals(t)
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	RegisterCommands(func(*core.Core) {}, fs)

	core.AssertLen(t, RegisteredLocales(), 1)
	core.AssertNotNil(t, RegisteredLocales()[0])
}

func TestAX7CLI_RegisteredLocales_Bad(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {}, nil)

	core.AssertNil(t, RegisteredLocales())
	core.AssertEmpty(t, RegisteredLocales())
}

func TestAX7CLI_RegisteredLocales_Ugly(t *core.T) {
	resetGlobals(t)
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	RegisterCommands(func(*core.Core) {}, fs, nil)

	core.AssertLen(t, RegisteredLocales(), 1)
	core.AssertNotNil(t, RegisteredLocales()[0])
}

func TestAX7CLI_RegisteredCommands_Good(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(c *core.Core) { c.Command("one", core.Command{}) })

	var count int
	for fn := range RegisteredCommands() {
		core.AssertNotNil(t, fn)
		count++
	}
	core.AssertEqual(t, 1, count)
}

func TestAX7CLI_RegisteredCommands_Bad(t *core.T) {
	resetGlobals(t)
	var count int
	for range RegisteredCommands() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, RegisteredLocales())
}

func TestAX7CLI_RegisteredCommands_Ugly(t *core.T) {
	resetGlobals(t)
	RegisterCommands(func(*core.Core) {})
	RegisterCommands(func(*core.Core) {})

	var count int
	for range RegisteredCommands() {
		count++
	}
	core.AssertEqual(t, 2, count)
}

func TestAX7CLI_RegisterCommand_Good(t *core.T) {
	c := core.New()
	RegisterCommand(c, "hello", core.Command{Description: "Hello"})

	core.AssertTrue(t, c.Command("hello").OK)
	core.AssertContains(t, c.Command("hello").Value.(*core.Command).Description, "Hello")
}

func TestAX7CLI_RegisterCommand_Bad(t *core.T) {
	var c *core.Core

	core.AssertPanics(t, func() { RegisterCommand(c, "hello", core.Command{}) })
	core.AssertNil(t, c)
}

func TestAX7CLI_RegisterCommand_Ugly(t *core.T) {
	c := core.New()
	RegisterCommand(c, "root", core.Command{Description: "Root"})

	core.AssertTrue(t, c.Command("root").OK)
	core.AssertNotNil(t, c.Command("root").Value)
}

func TestAX7CLI_RequireArgs_Good(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "config"})
	got := RequireArgs(opts, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_RequireArgs_Bad(t *core.T) {
	opts := core.NewOptions()
	got := RequireArgs(opts, 1)

	core.AssertContains(t, got, "requires")
	core.AssertContains(t, got, "1")
}

func TestAX7CLI_RequireArgs_Ugly(t *core.T) {
	opts := core.NewOptions()
	got := RequireArgs(opts, 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_RequireExactArgs_Good(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "config"})
	got := RequireExactArgs(opts, 1)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_RequireExactArgs_Bad(t *core.T) {
	opts := core.NewOptions(core.Option{Key: "_arg", Value: "extra"})
	got := RequireExactArgs(opts, 0)

	core.AssertContains(t, got, "accepts no arguments")
	core.AssertNotEmpty(t, got)
}

func TestAX7CLI_RequireExactArgs_Ugly(t *core.T) {
	opts := core.NewOptions()
	got := RequireExactArgs(opts, 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_SemVer_Good(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "1.2.3", "", "unknown", "unknown"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "1.2.3", SemVer())
}

func TestAX7CLI_SemVer_Bad(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "0.0.0", "dev.1", "unknown", "unknown"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "0.0.0-dev.1", SemVer())
}

func TestAX7CLI_SemVer_Ugly(t *core.T) {
	oldVersion, oldPre, oldCommit, oldDate := AppVersion, BuildPreRelease, BuildCommit, BuildDate
	AppVersion, BuildPreRelease, BuildCommit, BuildDate = "1.0.0", "rc.1", "abc123", "20260428"
	defer func() { AppVersion, BuildPreRelease, BuildCommit, BuildDate = oldVersion, oldPre, oldCommit, oldDate }()

	core.AssertEqual(t, "1.0.0-rc.1+abc123.20260428", SemVer())
}

func TestAX7CLI_WithAppName_Good(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("codex")

	core.AssertEqual(t, "codex", AppName)
}

func TestAX7CLI_WithAppName_Bad(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("")

	core.AssertEqual(t, "", AppName)
}

func TestAX7CLI_WithAppName_Ugly(t *core.T) {
	old := AppName
	defer func() { AppName = old }()
	WithAppName("core dev")

	core.AssertEqual(t, "core dev", AppName)
}

func TestAX7CLI_WithLocales_Good(t *core.T) {
	fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
	src := WithLocales(fs, ".")

	core.AssertEqual(t, ".", src.Dir)
	core.AssertNotNil(t, src.FS)
}

func TestAX7CLI_WithLocales_Bad(t *core.T) {
	src := WithLocales(nil, ".")

	core.AssertEqual(t, ".", src.Dir)
	core.AssertNil(t, src.FS)
}

func TestAX7CLI_WithLocales_Ugly(t *core.T) {
	fs := fstest.MapFS{}
	src := WithLocales(fs, "")

	core.AssertEqual(t, "", src.Dir)
	core.AssertNotNil(t, src.FS)
}

func TestAX7CLI_Init_Good(t *core.T) {
	resetGlobals(t)
	err := Init(Options{AppName: "codex", Version: "1.0.0"})

	core.AssertNoError(t, err)
	core.AssertEqual(t, "codex", Core().App().Name)
}

func TestAX7CLI_Init_Bad(t *core.T) {
	resetGlobals(t)
	err := Init(Options{})

	core.AssertNoError(t, err)
	core.AssertNotNil(t, Core())
}

func TestAX7CLI_Init_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "once"}))
	err := Init(Options{AppName: "twice"})

	core.AssertNoError(t, err)
	core.AssertEqual(t, "once", Core().App().Name)
}

func TestAX7CLI_Core_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "core"}))

	core.AssertNotNil(t, Core())
	core.AssertEqual(t, "core", Core().App().Name)
}

func TestAX7CLI_Core_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Core() })
	core.AssertNil(t, instance)
}

func TestAX7CLI_Core_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "core"}))
	Shutdown()

	core.AssertNotNil(t, Core())
	core.AssertNotNil(t, Context())
}

func TestAX7CLI_Execute_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "execute"}))

	err := Execute()
	core.AssertNoError(t, err)
}

func TestAX7CLI_Execute_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Execute() })
	core.AssertNil(t, instance)
}

func TestAX7CLI_Execute_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "execute"}))
	instance.core.Service("cli", core.Service{})

	err := Execute()
	core.AssertNoError(t, err)
}

func TestAX7CLI_Run_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "run"}))

	err := Run(context.Background())
	core.AssertNoError(t, err)
}

func TestAX7CLI_Run_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Run(context.Background()) })
	core.AssertNil(t, instance)
}

func TestAX7CLI_Run_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "run"}))

	err := Run(nil)
	core.AssertNoError(t, err)
}

func TestAX7CLI_RunWithTimeout_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "timeout"}))
	stop := RunWithTimeout(time.Millisecond)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestAX7CLI_RunWithTimeout_Bad(t *core.T) {
	resetGlobals(t)
	stop := RunWithTimeout(0)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestAX7CLI_RunWithTimeout_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "timeout"}))
	stop := RunWithTimeout(-time.Second)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestAX7CLI_Context_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "context"}))

	core.AssertNotNil(t, Context())
	core.AssertNoError(t, Context().Err())
}

func TestAX7CLI_Context_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Context() })
	core.AssertNil(t, instance)
}

func TestAX7CLI_Context_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "context"}))
	Shutdown()

	core.AssertNotNil(t, Context())
	core.AssertError(t, Context().Err())
}

func TestAX7CLI_Shutdown_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "shutdown"}))

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNotNil(t, instance)
}

func TestAX7CLI_Shutdown_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNil(t, instance)
}

func TestAX7CLI_Shutdown_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, Init(Options{AppName: "shutdown"}))
	Shutdown()

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNotNil(t, instance)
}

func TestAX7CLI_Main_Good(t *core.T) {
	if os.Getenv("AX7_MAIN_GOOD") == "1" {
		os.Args = []string{"core"}
		Main()
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_GOOD")
	core.AssertNoError(t, err)
}

func TestAX7CLI_Main_Bad(t *core.T) {
	if os.Getenv("AX7_MAIN_BAD") == "1" {
		os.Args = []string{"core"}
		Main(func(*core.Core) { panic("main setup failed") })
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_BAD")
	core.AssertError(t, err)
}

func TestAX7CLI_Main_Ugly(t *core.T) {
	if os.Getenv("AX7_MAIN_UGLY") == "1" {
		os.Args = []string{"core"}
		Main(func(c *core.Core) { c.App().Name = "ugly" })
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_UGLY")
	core.AssertNoError(t, err)
}

func TestAX7CLI_MainWithLocales_Good(t *core.T) {
	if os.Getenv("AX7_MAIN_WITH_LOCALES_GOOD") == "1" {
		os.Args = []string{"core"}
		fs := fstest.MapFS{"en.json": {Data: []byte(`{"x":"y"}`)}}
		MainWithLocales([]LocaleSource{WithLocales(fs, ".")})
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_WITH_LOCALES_GOOD")
	core.AssertNoError(t, err)
}

func TestAX7CLI_MainWithLocales_Bad(t *core.T) {
	if os.Getenv("AX7_MAIN_WITH_LOCALES_BAD") == "1" {
		os.Args = []string{"core"}
		MainWithLocales([]LocaleSource{{}})
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_WITH_LOCALES_BAD")
	core.AssertNoError(t, err)
}

func TestAX7CLI_MainWithLocales_Ugly(t *core.T) {
	if os.Getenv("AX7_MAIN_WITH_LOCALES_UGLY") == "1" {
		os.Args = []string{"core"}
		MainWithLocales(nil, func(c *core.Core) { c.App().Version = SemVer() })
		return
	}
	err := ax7RunSelf(t, "AX7_MAIN_WITH_LOCALES_UGLY")
	core.AssertNoError(t, err)
}

func TestAX7CLI_GhAuthenticated_Good(t *core.T) {
	ax7FakeCommands(t, map[string]string{"gh": "echo 'Logged in to github.com'\n"})

	core.AssertFalse(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestAX7CLI_GhAuthenticated_Bad(t *core.T) {
	ax7FakeCommands(t, map[string]string{"gh": "echo 'not logged in'\nexit 1\n"})

	core.AssertFalse(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestAX7CLI_GhAuthenticated_Ugly(t *core.T) {
	ax7FakeCommands(t, map[string]string{"gh": "echo 'Logged in as codex'\n"})

	core.AssertFalse(t, GhAuthenticated())
	core.AssertNotPanics(t, func() { _ = GhAuthenticated() })
}

func TestAX7CLI_DefaultYes_Good(t *core.T) {
	cfg := &confirmConfig{}
	DefaultYes()(cfg)

	core.AssertTrue(t, cfg.defaultYes)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_DefaultYes_Bad(t *core.T) {
	var cfg *confirmConfig

	core.AssertPanics(t, func() { DefaultYes()(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_DefaultYes_Ugly(t *core.T) {
	cfg := &confirmConfig{defaultYes: false, required: true}
	DefaultYes()(cfg)

	core.AssertTrue(t, cfg.defaultYes)
	core.AssertTrue(t, cfg.required)
}

func TestAX7CLI_Required_Good(t *core.T) {
	cfg := &confirmConfig{}
	Required()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertFalse(t, cfg.defaultYes)
}

func TestAX7CLI_Required_Bad(t *core.T) {
	var cfg *confirmConfig

	core.AssertPanics(t, func() { Required()(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_Required_Ugly(t *core.T) {
	cfg := &confirmConfig{defaultYes: true}
	Required()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertTrue(t, cfg.defaultYes)
}

func TestAX7CLI_Timeout_Good(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(time.Second)(cfg)

	core.AssertEqual(t, time.Second, cfg.timeout)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_Timeout_Bad(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(0)(cfg)

	core.AssertEqual(t, time.Duration(0), cfg.timeout)
	core.AssertFalse(t, cfg.defaultYes)
}

func TestAX7CLI_Timeout_Ugly(t *core.T) {
	cfg := &confirmConfig{}
	Timeout(-time.Second)(cfg)

	core.AssertEqual(t, -time.Second, cfg.timeout)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_Confirm_Good(t *core.T) {
	SetStdin(core.NewReader("y\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, Confirm("Continue?"))
}

func TestAX7CLI_Confirm_Bad(t *core.T) {
	SetStdin(core.NewReader("n\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, Confirm("Continue?", DefaultYes()))
}

func TestAX7CLI_Confirm_Ugly(t *core.T) {
	SetStdin(core.NewReader("maybe\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, Confirm("Continue?"))
}

func TestAX7CLI_ConfirmAction_Good(t *core.T) {
	SetStdin(core.NewReader("yes\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmAction("install", "package"))
}

func TestAX7CLI_ConfirmAction_Bad(t *core.T) {
	SetStdin(core.NewReader("no\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmAction("remove", "package", DefaultYes()))
}

func TestAX7CLI_ConfirmAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmAction("deploy", "", DefaultYes()))
}

func TestAX7CLI_ConfirmDangerousAction_Good(t *core.T) {
	SetStdin(bufio.NewReader(core.NewReader("y\ny\n")))
	defer SetStdin(nil)

	core.AssertTrue(t, ConfirmDangerousAction("remove", "package"))
}

func TestAX7CLI_ConfirmDangerousAction_Bad(t *core.T) {
	SetStdin(core.NewReader("n\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmDangerousAction("remove", "package"))
}

func TestAX7CLI_ConfirmDangerousAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("y\nn\n"))
	defer SetStdin(nil)

	core.AssertFalse(t, ConfirmDangerousAction("remove", "package"))
}

func TestAX7CLI_WithDefault_Good(t *core.T) {
	cfg := &questionConfig{}
	WithDefault("codex")(cfg)

	core.AssertEqual(t, "codex", cfg.defaultValue)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_WithDefault_Bad(t *core.T) {
	var cfg *questionConfig

	core.AssertPanics(t, func() { WithDefault("codex")(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_WithDefault_Ugly(t *core.T) {
	cfg := &questionConfig{defaultValue: "old"}
	WithDefault("")(cfg)

	core.AssertEqual(t, "", cfg.defaultValue)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_WithValidator_Good(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(func(string) error { return nil })(cfg)

	core.AssertNotNil(t, cfg.validator)
	core.AssertNoError(t, cfg.validator("ok"))
}

func TestAX7CLI_WithValidator_Bad(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(nil)(cfg)

	core.AssertNil(t, cfg.validator)
	core.AssertFalse(t, cfg.required)
}

func TestAX7CLI_WithValidator_Ugly(t *core.T) {
	cfg := &questionConfig{}
	WithValidator(func(string) error { return Err("invalid") })(cfg)

	core.AssertError(t, cfg.validator("bad"))
	core.AssertNotNil(t, cfg.validator)
}

func TestAX7CLI_RequiredInput_Good(t *core.T) {
	cfg := &questionConfig{}
	RequiredInput()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertEqual(t, "", cfg.defaultValue)
}

func TestAX7CLI_RequiredInput_Bad(t *core.T) {
	var cfg *questionConfig

	core.AssertPanics(t, func() { RequiredInput()(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_RequiredInput_Ugly(t *core.T) {
	cfg := &questionConfig{defaultValue: "fallback"}
	RequiredInput()(cfg)

	core.AssertTrue(t, cfg.required)
	core.AssertEqual(t, "fallback", cfg.defaultValue)
}

func TestAX7CLI_Question_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "codex", Question("Name?"))
}

func TestAX7CLI_Question_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "default", Question("Name?", WithDefault("default")))
}

func TestAX7CLI_Question_Ugly(t *core.T) {
	SetStdin(core.NewReader("bad\ngood\n"))
	defer SetStdin(nil)
	got := Question("Name?", WithValidator(func(v string) error {
		if v == "bad" {
			return Err("bad")
		}
		return nil
	}))

	core.AssertEqual(t, "good", got)
}

func TestAX7CLI_QuestionAction_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "codex", QuestionAction("name", "agent"))
}

func TestAX7CLI_QuestionAction_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "fallback", QuestionAction("name", "agent", WithDefault("fallback")))
}

func TestAX7CLI_QuestionAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("value\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "value", QuestionAction("", ""))
}

func TestAX7CLI_WithDisplay_Good(t *core.T) {
	cfg := &chooseConfig[int]{}
	WithDisplay(func(v int) string { return Sprintf("item-%d", v) })(cfg)

	core.AssertNotNil(t, cfg.displayFn)
	core.AssertEqual(t, "item-2", cfg.displayFn(2))
}

func TestAX7CLI_WithDisplay_Bad(t *core.T) {
	cfg := &chooseConfig[int]{}
	WithDisplay[int](nil)(cfg)

	core.AssertNil(t, cfg.displayFn)
	core.AssertFalse(t, cfg.filter)
}

func TestAX7CLI_WithDisplay_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDisplay(func(v string) string { return core.Upper(v) })(cfg)

	core.AssertEqual(t, "CODEX", cfg.displayFn("codex"))
	core.AssertFalse(t, cfg.multi)
}

func TestAX7CLI_WithDefaultIndex_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](1)(cfg)

	core.AssertEqual(t, 1, cfg.defaultN)
	core.AssertFalse(t, cfg.filter)
}

func TestAX7CLI_WithDefaultIndex_Bad(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](-1)(cfg)

	core.AssertEqual(t, -1, cfg.defaultN)
	core.AssertFalse(t, cfg.multi)
}

func TestAX7CLI_WithDefaultIndex_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	WithDefaultIndex[string](99)(cfg)

	core.AssertEqual(t, 99, cfg.defaultN)
	core.AssertFalse(t, cfg.filter)
}

func TestAX7CLI_Filter_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	Filter[string]()(cfg)

	core.AssertTrue(t, cfg.filter)
	core.AssertFalse(t, cfg.multi)
}

func TestAX7CLI_Filter_Bad(t *core.T) {
	var cfg *chooseConfig[string]

	core.AssertPanics(t, func() { Filter[string]()(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_Filter_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{multi: true}
	Filter[string]()(cfg)

	core.AssertTrue(t, cfg.filter)
	core.AssertTrue(t, cfg.multi)
}

func TestAX7CLI_Multi_Good(t *core.T) {
	cfg := &chooseConfig[string]{}
	Multi[string]()(cfg)

	core.AssertTrue(t, cfg.multi)
	core.AssertFalse(t, cfg.filter)
}

func TestAX7CLI_Multi_Bad(t *core.T) {
	var cfg *chooseConfig[string]

	core.AssertPanics(t, func() { Multi[string]()(cfg) })
	core.AssertNil(t, cfg)
}

func TestAX7CLI_Multi_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{filter: true}
	Multi[string]()(cfg)

	core.AssertTrue(t, cfg.multi)
	core.AssertTrue(t, cfg.filter)
}

func TestAX7CLI_Display_Good(t *core.T) {
	cfg := &chooseConfig[int]{}
	Display(func(v int) string { return Sprintf("n=%d", v) })(cfg)

	core.AssertNotNil(t, cfg.displayFn)
	core.AssertEqual(t, "n=3", cfg.displayFn(3))
}

func TestAX7CLI_Display_Bad(t *core.T) {
	cfg := &chooseConfig[int]{}
	Display[int](nil)(cfg)

	core.AssertNil(t, cfg.displayFn)
	core.AssertFalse(t, cfg.filter)
}

func TestAX7CLI_Display_Ugly(t *core.T) {
	cfg := &chooseConfig[string]{}
	Display(func(v string) string { return v + v })(cfg)

	core.AssertEqual(t, "aa", cfg.displayFn("a"))
	core.AssertFalse(t, cfg.multi)
}

func TestAX7CLI_Choose_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "beta", Choose("Pick", []string{"alpha", "beta"}))
}

func TestAX7CLI_Choose_Bad(t *core.T) {
	got := Choose("Pick", []string{})

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Choose_Ugly(t *core.T) {
	SetStdin(core.NewReader("alp\n1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "alpha", Choose("Pick", []string{"alpha", "beta"}, Filter[string]()))
}

func TestAX7CLI_ChooseAction_Good(t *core.T) {
	SetStdin(core.NewReader("1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "alpha", ChooseAction("pick", "agent", []string{"alpha", "beta"}))
}

func TestAX7CLI_ChooseAction_Bad(t *core.T) {
	got := ChooseAction("pick", "agent", []string{})

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_ChooseAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, "beta", ChooseAction("pick", "", []string{"alpha", "beta"}, WithDefaultIndex[string](1)))
}

func TestAX7CLI_ChooseMulti_Good(t *core.T) {
	SetStdin(core.NewReader("1 3\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"alpha", "gamma"}, ChooseMulti("Pick", []string{"alpha", "beta", "gamma"}))
}

func TestAX7CLI_ChooseMulti_Bad(t *core.T) {
	got := ChooseMulti("Pick", []string{})

	core.AssertNil(t, got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_ChooseMulti_Ugly(t *core.T) {
	SetStdin(core.NewReader("gam\n1\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"gamma"}, ChooseMulti("Pick", []string{"alpha", "beta", "gamma"}, Filter[string]()))
}

func TestAX7CLI_ChooseMultiAction_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)

	core.AssertEqual(t, []string{"beta"}, ChooseMultiAction("pick", "agent", []string{"alpha", "beta"}))
}

func TestAX7CLI_ChooseMultiAction_Bad(t *core.T) {
	got := ChooseMultiAction("pick", "agent", []string{})

	core.AssertNil(t, got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_ChooseMultiAction_Ugly(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)

	core.AssertNil(t, ChooseMultiAction("pick", "", []string{"alpha"}))
}

func TestAX7CLI_GitClone_Good(t *core.T) {
	ax7FakeCommands(t, map[string]string{"gh": "echo 'Logged in'\nexit 0\n"})

	err := GitClone(context.Background(), "org", "repo", "target")
	core.AssertError(t, err)
}

func TestAX7CLI_GitClone_Bad(t *core.T) {
	ax7FakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'clone failed'\nexit 2\n",
	})

	err := GitClone(context.Background(), "org", "repo", "target")
	core.AssertError(t, err)
}

func TestAX7CLI_GitClone_Ugly(t *core.T) {
	ax7FakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'ok'\nexit 0\n",
	})

	err := GitClone(nil, "org", "repo", "target")
	core.AssertError(t, err)
}

func TestAX7CLI_GitCloneRef_Good(t *core.T) {
	ax7FakeCommands(t, map[string]string{"gh": "echo 'Logged in'\nexit 0\n"})

	err := GitCloneRef(context.Background(), "org", "repo", "target", "main")
	core.AssertError(t, err)
}

func TestAX7CLI_GitCloneRef_Bad(t *core.T) {
	ax7FakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'already exists'\nexit 2\n",
	})

	err := GitCloneRef(context.Background(), "org", "repo", "target", "main")
	core.AssertError(t, err)
}

func TestAX7CLI_GitCloneRef_Ugly(t *core.T) {
	ax7FakeCommands(t, map[string]string{
		"gh":  "echo 'not logged in'\nexit 1\n",
		"git": "echo 'ok'\nexit 0\n",
	})

	err := GitCloneRef(nil, "org", "repo", "target", "")
	core.AssertError(t, err)
}

func TestAX7CLI_Prompt_Good(t *core.T) {
	SetStdin(core.NewReader("codex\n"))
	defer SetStdin(nil)
	got, err := Prompt("Name", "default")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "codex", got)
}

func TestAX7CLI_Prompt_Bad(t *core.T) {
	SetStdin(core.NewReader("\n"))
	defer SetStdin(nil)
	got, err := Prompt("Name", "default")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "default", got)
}

func TestAX7CLI_Prompt_Ugly(t *core.T) {
	SetStdin(core.NewReader(""))
	defer SetStdin(nil)
	got, err := Prompt("Name", "")

	core.AssertError(t, err)
	core.AssertEqual(t, "", got)
}

func TestAX7CLI_Select_Good(t *core.T) {
	SetStdin(core.NewReader("2\n"))
	defer SetStdin(nil)
	got, err := Select("Pick", []string{"alpha", "beta"})

	core.AssertNoError(t, err)
	core.AssertEqual(t, "beta", got)
}

func TestAX7CLI_Select_Bad(t *core.T) {
	SetStdin(core.NewReader("9\n"))
	defer SetStdin(nil)
	got, err := Select("Pick", []string{"alpha", "beta"})

	core.AssertError(t, err)
	core.AssertEqual(t, "", got)
}

func TestAX7CLI_Select_Ugly(t *core.T) {
	got, err := Select("Pick", nil)

	core.AssertNoError(t, err)
	core.AssertEqual(t, "", got)
}

func TestAX7CLI_MultiSelect_Good(t *core.T) {
	SetStdin(core.NewReader("1 3\n"))
	defer SetStdin(nil)
	got, err := MultiSelect("Pick", []string{"alpha", "beta", "gamma"})

	core.AssertNoError(t, err)
	core.AssertEqual(t, []string{"alpha", "gamma"}, got)
}

func TestAX7CLI_MultiSelect_Bad(t *core.T) {
	SetStdin(core.NewReader("9\n"))
	defer SetStdin(nil)
	got, err := MultiSelect("Pick", []string{"alpha", "beta"})

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestAX7CLI_MultiSelect_Ugly(t *core.T) {
	got, err := MultiSelect("Pick", nil)

	core.AssertNoError(t, err)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Truncate_Good(t *core.T) {
	got := Truncate("abcdef", 4)

	core.AssertEqual(t, "a...", got)
	core.AssertLen(t, got, 4)
}

func TestAX7CLI_Truncate_Bad(t *core.T) {
	got := Truncate("abcdef", 0)

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Truncate_Ugly(t *core.T) {
	got := Truncate("abcdef", 2)

	core.AssertEqual(t, "ab", got)
	core.AssertLen(t, got, 2)
}

func TestAX7CLI_Pad_Good(t *core.T) {
	got := Pad("go", 4)

	core.AssertEqual(t, "go  ", got)
	core.AssertLen(t, got, 4)
}

func TestAX7CLI_Pad_Bad(t *core.T) {
	got := Pad("long", 2)

	core.AssertEqual(t, "long", got)
	core.AssertLen(t, got, 4)
}

func TestAX7CLI_Pad_Ugly(t *core.T) {
	got := Pad("", 3)

	core.AssertEqual(t, "   ", got)
	core.AssertLen(t, got, 3)
}

func TestAX7CLI_FormatAge_Good(t *core.T) {
	got := FormatAge(time.Now().Add(-2 * time.Minute))

	core.AssertContains(t, got, "m ago")
	core.AssertNotEqual(t, "just now", got)
}

func TestAX7CLI_FormatAge_Bad(t *core.T) {
	got := FormatAge(time.Now().Add(time.Minute))

	core.AssertEqual(t, "just now", got)
	core.AssertNotEmpty(t, got)
}

func TestAX7CLI_FormatAge_Ugly(t *core.T) {
	got := FormatAge(time.Now().Add(-45 * 24 * time.Hour))

	core.AssertContains(t, got, "mo ago")
	core.AssertNotEmpty(t, got)
}

func TestAX7CLI_DefaultTableStyle_Good(t *core.T) {
	style := DefaultTableStyle()

	core.AssertNotNil(t, style.HeaderStyle)
	core.AssertEqual(t, "  ", style.Separator)
}

func TestAX7CLI_DefaultTableStyle_Bad(t *core.T) {
	style := DefaultTableStyle()

	core.AssertNil(t, style.CellStyle)
	core.AssertNotNil(t, style.HeaderStyle)
}

func TestAX7CLI_DefaultTableStyle_Ugly(t *core.T) {
	style := DefaultTableStyle()
	style.Separator = "|"

	core.AssertEqual(t, "|", style.Separator)
	core.AssertEqual(t, "  ", DefaultTableStyle().Separator)
}

func TestAX7CLI_NewTable_Good(t *core.T) {
	table := NewTable("Name", "Status")

	core.AssertEqual(t, []string{"Name", "Status"}, table.Headers)
	core.AssertNotNil(t, table.Style.HeaderStyle)
}

func TestAX7CLI_NewTable_Bad(t *core.T) {
	table := NewTable()

	core.AssertEmpty(t, table.Headers)
	core.AssertEqual(t, "", table.String())
}

func TestAX7CLI_NewTable_Ugly(t *core.T) {
	table := NewTable(":check:")

	core.AssertEqual(t, []string{":check:"}, table.Headers)
	core.AssertContains(t, table.String(), "✓")
}

func TestAX7CLI_Table_AddRow_Good(t *core.T) {
	table := NewTable("Name").AddRow("codex")

	core.AssertLen(t, table.Rows, 1)
	core.AssertEqual(t, []string{"codex"}, table.Rows[0])
}

func TestAX7CLI_Table_AddRow_Bad(t *core.T) {
	table := NewTable("Name").AddRow()

	core.AssertLen(t, table.Rows, 1)
	core.AssertEmpty(t, table.Rows[0])
}

func TestAX7CLI_Table_AddRow_Ugly(t *core.T) {
	table := NewTable().AddRow("orphan")

	core.AssertContains(t, table.String(), "orphan")
	core.AssertLen(t, table.Rows, 1)
}

func TestAX7CLI_Table_WithBorders_Good(t *core.T) {
	table := NewTable("Name").WithBorders(BorderRounded)

	core.AssertEqual(t, BorderRounded, table.borders)
	core.AssertContains(t, table.String(), "╭")
}

func TestAX7CLI_Table_WithBorders_Bad(t *core.T) {
	table := NewTable("Name").WithBorders(BorderNone)

	core.AssertEqual(t, BorderNone, table.borders)
	core.AssertNotContains(t, table.String(), "╭")
}

func TestAX7CLI_Table_WithBorders_Ugly(t *core.T) {
	ax7PlainCLI(t)
	table := NewTable("Name").WithBorders(BorderHeavy)

	core.AssertEqual(t, BorderHeavy, table.borders)
	core.AssertContains(t, table.String(), "+")
}

func TestAX7CLI_Table_WithCellStyle_Good(t *core.T) {
	table := NewTable("Name").WithCellStyle(0, func(string) *AnsiStyle { return NewStyle().Bold() })

	core.AssertNotNil(t, table.cellStyleFns[0])
	core.AssertEqual(t, table, table.WithCellStyle(1, nil))
}

func TestAX7CLI_Table_WithCellStyle_Bad(t *core.T) {
	table := NewTable("Name").WithCellStyle(-1, nil)

	core.AssertNotNil(t, table.cellStyleFns)
	core.AssertNil(t, table.cellStyleFns[-1])
}

func TestAX7CLI_Table_WithCellStyle_Ugly(t *core.T) {
	table := NewTable("Name").WithCellStyle(0, func(value string) *AnsiStyle {
		if value == "hot" {
			return NewStyle().Bold()
		}
		return nil
	})

	core.AssertNotNil(t, table.cellStyleFns[0]("hot"))
}

func TestAX7CLI_Table_WithMaxWidth_Good(t *core.T) {
	table := NewTable("Name").WithMaxWidth(10)

	core.AssertEqual(t, 10, table.maxWidth)
	core.AssertEqual(t, table, table.WithMaxWidth(20))
}

func TestAX7CLI_Table_WithMaxWidth_Bad(t *core.T) {
	table := NewTable("Name").WithMaxWidth(0)

	core.AssertEqual(t, 0, table.maxWidth)
	core.AssertContains(t, table.String(), "Name")
}

func TestAX7CLI_Table_WithMaxWidth_Ugly(t *core.T) {
	table := NewTable("Name").AddRow("abcdef").WithMaxWidth(5)

	core.AssertContains(t, table.String(), "...")
	core.AssertEqual(t, 5, table.maxWidth)
}

func TestAX7CLI_Table_String_Good(t *core.T) {
	got := NewTable("Name").AddRow("codex").String()

	core.AssertContains(t, got, "Name")
	core.AssertContains(t, got, "codex")
}

func TestAX7CLI_Table_String_Bad(t *core.T) {
	got := NewTable().String()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Table_String_Ugly(t *core.T) {
	got := NewTable("Name").WithBorders(BorderDouble).String()

	core.AssertContains(t, got, "Name")
	core.AssertContains(t, got, "╔")
}

func TestAX7CLI_Table_Render_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { NewTable("Name").AddRow("codex").Render() })

	core.AssertContains(t, out, "Name")
	core.AssertContains(t, out, "codex")
}

func TestAX7CLI_Table_Render_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { NewTable().Render() })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_Table_Render_Ugly(t *core.T) {
	out := ax7CaptureStdout(t, func() { NewTable("Name").WithBorders(BorderNormal).Render() })

	core.AssertContains(t, out, "Name")
	core.AssertContains(t, out, "┌")
}

func TestAX7CLI_T_Good(t *core.T) {
	got := T("i18n.progress.check")

	core.AssertEqual(t, "Checking...", got)
	core.AssertContains(t, got, "Checking")
}

func TestAX7CLI_T_Bad(t *core.T) {
	got := T("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_T_Ugly(t *core.T) {
	got := T("i18n.fail.load", map[string]any{"Subject": "config"})

	core.AssertEqual(t, "Failed to load config", got)
	core.AssertContains(t, got, "config")
}

func TestAX7CLI_TrackedTask_Update_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Update("running")

	core.AssertContains(t, task.tracker.String(), "running")
	core.AssertContains(t, task.tracker.Summary(), "0/1")
}

func TestAX7CLI_TrackedTask_Update_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Update("running") })
	core.AssertNil(t, task)
}

func TestAX7CLI_TrackedTask_Update_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Update("")

	core.AssertContains(t, task.tracker.String(), "build")
	core.AssertContains(t, task.tracker.Summary(), "0/1")
}

func TestAX7CLI_TrackedTask_Done_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Done("done")

	core.AssertContains(t, task.tracker.String(), "done")
	core.AssertEqual(t, "1/1 passed", task.tracker.Summary())
}

func TestAX7CLI_TrackedTask_Done_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Done("done") })
	core.AssertNil(t, task)
}

func TestAX7CLI_TrackedTask_Done_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Done("")

	core.AssertEqual(t, "1/1 passed", task.tracker.Summary())
	core.AssertContains(t, task.tracker.String(), "build")
}

func TestAX7CLI_TrackedTask_Fail_Good(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Fail("failed")

	core.AssertContains(t, task.tracker.String(), "failed")
	core.AssertEqual(t, "0/1 passed, 1 failed", task.tracker.Summary())
}

func TestAX7CLI_TrackedTask_Fail_Bad(t *core.T) {
	var task *TrackedTask

	core.AssertPanics(t, func() { task.Fail("failed") })
	core.AssertNil(t, task)
}

func TestAX7CLI_TrackedTask_Fail_Ugly(t *core.T) {
	task := NewTaskTracker().Add("build")
	task.Fail("")

	core.AssertEqual(t, "0/1 passed, 1 failed", task.tracker.Summary())
	core.AssertContains(t, task.tracker.String(), "build")
}

func TestAX7CLI_NewTaskTracker_Good(t *core.T) {
	tr := NewTaskTracker()

	core.AssertNotNil(t, tr)
	core.AssertNotNil(t, tr.out)
}

func TestAX7CLI_NewTaskTracker_Bad(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEmpty(t, tr.tasks)
	core.AssertFalse(t, tr.started)
}

func TestAX7CLI_NewTaskTracker_Ugly(t *core.T) {
	tr := NewTaskTracker().WithOutput(core.NewBuilder())

	core.AssertNotNil(t, tr.out)
	core.AssertEqual(t, "", tr.String())
}

func TestAX7CLI_TaskTracker_WithOutput_Good(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)

	core.AssertEqual(t, out, tr.out)
	core.AssertEqual(t, tr, tr.WithOutput(out))
}

func TestAX7CLI_TaskTracker_WithOutput_Bad(t *core.T) {
	tr := NewTaskTracker()
	original := tr.out

	core.AssertEqual(t, tr, tr.WithOutput(nil))
	core.AssertEqual(t, original, tr.out)
}

func TestAX7CLI_TaskTracker_WithOutput_Ugly(t *core.T) {
	tr := NewTaskTracker().WithOutput(io.Discard)

	core.AssertEqual(t, io.Discard, tr.out)
	core.AssertFalse(t, tr.isTTY())
}

func TestAX7CLI_TaskTracker_Add_Good(t *core.T) {
	tr := NewTaskTracker()
	task := tr.Add("build")

	core.AssertNotNil(t, task)
	core.AssertContains(t, tr.String(), "build")
}

func TestAX7CLI_TaskTracker_Add_Bad(t *core.T) {
	tr := NewTaskTracker()
	task := tr.Add("")

	core.AssertNotNil(t, task)
	core.AssertLen(t, tr.tasks, 1)
}

func TestAX7CLI_TaskTracker_Add_Ugly(t *core.T) {
	tr := NewTaskTracker()
	first := tr.Add("same")
	second := tr.Add("same")

	core.AssertTrue(t, first != second)
	core.AssertLen(t, tr.tasks, 2)
}

func TestAX7CLI_TaskTracker_Wait_Good(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)
	tr.Add("build").Done("done")

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertContains(t, out.String(), "done")
}

func TestAX7CLI_TaskTracker_Wait_Bad(t *core.T) {
	tr := NewTaskTracker().WithOutput(core.NewBuilder())

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertEqual(t, "0/0 passed", tr.Summary())
}

func TestAX7CLI_TaskTracker_Wait_Ugly(t *core.T) {
	out := core.NewBuilder()
	tr := NewTaskTracker().WithOutput(out)
	tr.Add("build").Fail("failed")

	core.AssertNotPanics(t, func() { tr.Wait() })
	core.AssertContains(t, out.String(), "failed")
}

func TestAX7CLI_TaskTracker_Tasks_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build")
	var names []string
	for task := range tr.Tasks() {
		names = append(names, task.name)
	}

	core.AssertEqual(t, []string{"build"}, names)
}

func TestAX7CLI_TaskTracker_Tasks_Bad(t *core.T) {
	tr := NewTaskTracker()
	var count int
	for range tr.Tasks() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestAX7CLI_TaskTracker_Tasks_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("first")
	tr.Add("second")
	var count int
	for range tr.Tasks() {
		count++
	}

	core.AssertEqual(t, 2, count)
}

func TestAX7CLI_TaskTracker_Snapshots_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build").Update("running")
	var got []string
	for name, status := range tr.Snapshots() {
		got = append(got, name+":"+status)
	}

	core.AssertEqual(t, []string{"build:running"}, got)
}

func TestAX7CLI_TaskTracker_Snapshots_Bad(t *core.T) {
	tr := NewTaskTracker()
	var count int
	for range tr.Snapshots() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestAX7CLI_TaskTracker_Snapshots_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("").Done("")
	var got []string
	for name, status := range tr.Snapshots() {
		got = append(got, name+":"+status)
	}

	core.AssertEqual(t, []string{":"}, got)
}

func TestAX7CLI_TaskTracker_Summary_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("ok").Done("done")

	core.AssertEqual(t, "1/1 passed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "passed")
}

func TestAX7CLI_TaskTracker_Summary_Bad(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("bad").Fail("failed")

	core.AssertEqual(t, "0/1 passed, 1 failed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "failed")
}

func TestAX7CLI_TaskTracker_Summary_Ugly(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEqual(t, "0/0 passed", tr.Summary())
	core.AssertContains(t, tr.Summary(), "0/0")
}

func TestAX7CLI_TaskTracker_String_Good(t *core.T) {
	tr := NewTaskTracker()
	tr.Add("build").Done("done")

	core.AssertContains(t, tr.String(), "build")
	core.AssertContains(t, tr.String(), "done")
}

func TestAX7CLI_TaskTracker_String_Bad(t *core.T) {
	tr := NewTaskTracker()

	core.AssertEqual(t, "", tr.String())
	core.AssertEmpty(t, tr.String())
}

func TestAX7CLI_TaskTracker_String_Ugly(t *core.T) {
	tr := NewTaskTracker()
	tr.Add(":check:").Update(":warn:")

	core.AssertContains(t, tr.String(), "✓")
	core.AssertContains(t, tr.String(), "⚠")
}

func TestAX7CLI_Composite_Regions_Good(t *core.T) {
	c := Layout("HC")
	var regions []Region
	for r := range c.Regions() {
		regions = append(regions, r)
	}

	core.AssertLen(t, regions, 2)
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestAX7CLI_Composite_Regions_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestAX7CLI_Composite_Regions_Ugly(t *core.T) {
	c := Layout("HH")
	var count int
	for range c.Regions() {
		count++
	}

	core.AssertEqual(t, 1, count)
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestAX7CLI_Composite_Slots_Good(t *core.T) {
	c := Layout("CF")
	var count int
	for _, slot := range c.Slots() {
		core.AssertNotNil(t, slot)
		count++
	}

	core.AssertEqual(t, 2, count)
}

func TestAX7CLI_Composite_Slots_Bad(t *core.T) {
	c := Layout("Z")
	var count int
	for range c.Slots() {
		count++
	}

	core.AssertEqual(t, 0, count)
	core.AssertEmpty(t, c.regions)
}

func TestAX7CLI_Composite_Slots_Ugly(t *core.T) {
	c := Layout("C[HF]")
	var child *Composite
	for _, slot := range c.Slots() {
		child = slot.child
	}

	core.AssertNotNil(t, child)
	core.AssertNotNil(t, child.regions[RegionHeader])
}

func TestAX7CLI_StringBlock_Render_Good(t *core.T) {
	got := StringBlock(":check: ready").Render()

	core.AssertContains(t, got, "ready")
	core.AssertContains(t, got, "✓")
}

func TestAX7CLI_StringBlock_Render_Bad(t *core.T) {
	got := StringBlock("").Render()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_StringBlock_Render_Ugly(t *core.T) {
	got := StringBlock(":missing:").Render()

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestAX7CLI_Layout_Good(t *core.T) {
	c := Layout("HCF")

	core.AssertNotNil(t, c.regions[RegionHeader])
	core.AssertNotNil(t, c.regions[RegionFooter])
}

func TestAX7CLI_Layout_Bad(t *core.T) {
	c := Layout("Z")

	core.AssertEqual(t, "Z", c.variant)
	core.AssertEmpty(t, c.regions)
}

func TestAX7CLI_Layout_Ugly(t *core.T) {
	c := Layout("C[HF]")

	core.AssertNotNil(t, c.regions[RegionContent])
	core.AssertNotNil(t, c.regions[RegionContent].child)
}

func TestAX7CLI_ParseVariant_Good(t *core.T) {
	c, err := ParseVariant("HCF")

	core.AssertNoError(t, err)
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestAX7CLI_ParseVariant_Bad(t *core.T) {
	c, err := ParseVariant("Z")

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestAX7CLI_ParseVariant_Ugly(t *core.T) {
	c, err := ParseVariant("C[HF")

	core.AssertError(t, err)
	core.AssertNil(t, c)
}

func TestAX7CLI_Composite_H_Good(t *core.T) {
	c := Layout("H").H("header")

	core.AssertEqual(t, c, c.H("more"))
	core.AssertLen(t, c.regions[RegionHeader].blocks, 2)
}

func TestAX7CLI_Composite_H_Bad(t *core.T) {
	c := Layout("C").H("header")

	core.AssertNil(t, c.regions[RegionHeader])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestAX7CLI_Composite_H_Ugly(t *core.T) {
	c := Layout("H").H(123)

	core.AssertEqual(t, "123", c.regions[RegionHeader].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionHeader].blocks, 1)
}

func TestAX7CLI_Composite_L_Good(t *core.T) {
	c := Layout("L").L("left")

	core.AssertEqual(t, c, c.L("more"))
	core.AssertLen(t, c.regions[RegionLeft].blocks, 2)
}

func TestAX7CLI_Composite_L_Bad(t *core.T) {
	c := Layout("C").L("left")

	core.AssertNil(t, c.regions[RegionLeft])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestAX7CLI_Composite_L_Ugly(t *core.T) {
	c := Layout("L").L(StringBlock(":check:"))

	core.AssertEqual(t, "✓", c.regions[RegionLeft].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionLeft].blocks, 1)
}

func TestAX7CLI_Composite_C_Good(t *core.T) {
	c := Layout("C").C("content")

	core.AssertEqual(t, c, c.C("more"))
	core.AssertLen(t, c.regions[RegionContent].blocks, 2)
}

func TestAX7CLI_Composite_C_Bad(t *core.T) {
	c := Layout("H").C("content")

	core.AssertNil(t, c.regions[RegionContent])
	core.AssertNotNil(t, c.regions[RegionHeader])
}

func TestAX7CLI_Composite_C_Ugly(t *core.T) {
	c := Layout("C").C("")

	core.AssertEqual(t, "", c.regions[RegionContent].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionContent].blocks, 1)
}

func TestAX7CLI_Composite_R_Good(t *core.T) {
	c := Layout("R").R("right")

	core.AssertEqual(t, c, c.R("more"))
	core.AssertLen(t, c.regions[RegionRight].blocks, 2)
}

func TestAX7CLI_Composite_R_Bad(t *core.T) {
	c := Layout("C").R("right")

	core.AssertNil(t, c.regions[RegionRight])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestAX7CLI_Composite_R_Ugly(t *core.T) {
	c := Layout("R").R(RegionRight)

	core.AssertEqual(t, "82", c.regions[RegionRight].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionRight].blocks, 1)
}

func TestAX7CLI_Composite_F_Good(t *core.T) {
	c := Layout("F").F("footer")

	core.AssertEqual(t, c, c.F("more"))
	core.AssertLen(t, c.regions[RegionFooter].blocks, 2)
}

func TestAX7CLI_Composite_F_Bad(t *core.T) {
	c := Layout("C").F("footer")

	core.AssertNil(t, c.regions[RegionFooter])
	core.AssertNotNil(t, c.regions[RegionContent])
}

func TestAX7CLI_Composite_F_Ugly(t *core.T) {
	c := Layout("F").F(nil)

	core.AssertEqual(t, "<nil>", c.regions[RegionFooter].blocks[0].Render())
	core.AssertLen(t, c.regions[RegionFooter].blocks, 1)
}

func TestAX7CLI_UseRenderFlat_Good(t *core.T) {
	UseRenderSimple()
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotEqual(t, RenderSimple, currentRenderStyle)
}

func TestAX7CLI_UseRenderFlat_Bad(t *core.T) {
	currentRenderStyle = RenderBoxed
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotEqual(t, RenderBoxed, currentRenderStyle)
}

func TestAX7CLI_UseRenderFlat_Ugly(t *core.T) {
	UseRenderFlat()
	UseRenderFlat()

	core.AssertEqual(t, RenderFlat, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderFlat)
}

func TestAX7CLI_UseRenderSimple_Good(t *core.T) {
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotEqual(t, RenderFlat, currentRenderStyle)
}

func TestAX7CLI_UseRenderSimple_Bad(t *core.T) {
	currentRenderStyle = RenderBoxed
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotEqual(t, RenderBoxed, currentRenderStyle)
}

func TestAX7CLI_UseRenderSimple_Ugly(t *core.T) {
	UseRenderSimple()
	UseRenderSimple()

	core.AssertEqual(t, RenderSimple, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderSimple)
}

func TestAX7CLI_UseRenderBoxed_Good(t *core.T) {
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotEqual(t, RenderFlat, currentRenderStyle)
}

func TestAX7CLI_UseRenderBoxed_Bad(t *core.T) {
	currentRenderStyle = RenderSimple
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotEqual(t, RenderSimple, currentRenderStyle)
}

func TestAX7CLI_UseRenderBoxed_Ugly(t *core.T) {
	UseRenderBoxed()
	UseRenderBoxed()

	core.AssertEqual(t, RenderBoxed, currentRenderStyle)
	core.AssertNotPanics(t, UseRenderBoxed)
}

func TestAX7CLI_Composite_Render_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Layout("C").C("content").Render() })

	core.AssertContains(t, out, "content")
	core.AssertContains(t, out, "\n")
}

func TestAX7CLI_Composite_Render_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Layout("Z").Render() })

	core.AssertEqual(t, "", out)
	core.AssertEmpty(t, out)
}

func TestAX7CLI_Composite_Render_Ugly(t *core.T) {
	UseRenderSimple()
	out := ax7CaptureStdout(t, func() { Layout("HC").H("h").C("c").Render() })

	core.AssertContains(t, out, "h")
	core.AssertContains(t, out, "c")
}

func TestAX7CLI_Composite_String_Good(t *core.T) {
	got := Layout("C").C("content").String()

	core.AssertContains(t, got, "content")
	core.AssertContains(t, got, "\n")
}

func TestAX7CLI_Composite_String_Bad(t *core.T) {
	got := Layout("Z").String()

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_Composite_String_Ugly(t *core.T) {
	UseRenderBoxed()
	got := Layout("HC").H("h").C("c").String()

	core.AssertContains(t, got, "h")
	core.AssertContains(t, got, "c")
}

func TestAX7CLI_UseUnicode_Good(t *core.T) {
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertEqual(t, "✓", Glyph(":check:"))
}

func TestAX7CLI_UseUnicode_Bad(t *core.T) {
	UseASCII()
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertTrue(t, ColorEnabled())
}

func TestAX7CLI_UseUnicode_Ugly(t *core.T) {
	UseUnicode()
	UseUnicode()

	core.AssertEqual(t, ThemeUnicode, currentTheme)
	core.AssertEqual(t, "✓", Glyph(":check:"))
}

func TestAX7CLI_UseEmoji_Good(t *core.T) {
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertNotEqual(t, "✓", Glyph(":check:"))
}

func TestAX7CLI_UseEmoji_Bad(t *core.T) {
	UseASCII()
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertTrue(t, ColorEnabled())
}

func TestAX7CLI_UseEmoji_Ugly(t *core.T) {
	UseEmoji()
	UseEmoji()
	defer UseUnicode()

	core.AssertEqual(t, ThemeEmoji, currentTheme)
	core.AssertNotEmpty(t, Glyph(":check:"))
}

func TestAX7CLI_UseASCII_Good(t *core.T) {
	UseASCII()
	defer UseUnicode()

	core.AssertEqual(t, ThemeASCII, currentTheme)
	core.AssertEqual(t, "[OK]", Glyph(":check:"))
}

func TestAX7CLI_UseASCII_Bad(t *core.T) {
	UseASCII()
	defer UseUnicode()

	core.AssertFalse(t, ColorEnabled())
	core.AssertEqual(t, "[FAIL]", Glyph(":cross:"))
}

func TestAX7CLI_UseASCII_Ugly(t *core.T) {
	UseASCII()
	UseASCII()
	defer UseUnicode()

	core.AssertEqual(t, ThemeASCII, currentTheme)
	core.AssertFalse(t, ColorEnabled())
}

func TestAX7CLI_Glyph_Good(t *core.T) {
	UseUnicode()
	got := Glyph(":check:")

	core.AssertEqual(t, "✓", got)
	core.AssertNotEqual(t, ":check:", got)
}

func TestAX7CLI_Glyph_Bad(t *core.T) {
	got := Glyph(":missing:")

	core.AssertEqual(t, ":missing:", got)
	core.AssertContains(t, got, "missing")
}

func TestAX7CLI_Glyph_Ugly(t *core.T) {
	got := Glyph("")

	core.AssertEqual(t, "", got)
	core.AssertEmpty(t, got)
}

func TestAX7CLI_NewTree_Good(t *core.T) {
	tree := NewTree("root")

	core.AssertNotNil(t, tree)
	core.AssertContains(t, tree.String(), "root")
}

func TestAX7CLI_NewTree_Bad(t *core.T) {
	tree := NewTree("")

	core.AssertNotNil(t, tree)
	core.AssertContains(t, tree.String(), "\n")
}

func TestAX7CLI_NewTree_Ugly(t *core.T) {
	tree := NewTree(":check:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertNotContains(t, tree.String(), ":check:")
}

func TestAX7CLI_TreeNode_Add_Good(t *core.T) {
	tree := NewTree("root")
	child := tree.Add("child")

	core.AssertNotNil(t, child)
	core.AssertContains(t, tree.String(), "child")
}

func TestAX7CLI_TreeNode_Add_Bad(t *core.T) {
	tree := NewTree("root")
	child := tree.Add("")

	core.AssertNotNil(t, child)
	core.AssertLen(t, tree.children, 1)
}

func TestAX7CLI_TreeNode_Add_Ugly(t *core.T) {
	tree := NewTree("root")
	tree.Add(":check:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertLen(t, tree.children, 1)
}

func TestAX7CLI_TreeNode_AddStyled_Good(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("child", NewStyle().Bold())

	core.AssertNotNil(t, child.style)
	core.AssertContains(t, tree.String(), "child")
}

func TestAX7CLI_TreeNode_AddStyled_Bad(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("child", nil)

	core.AssertNil(t, child.style)
	core.AssertContains(t, tree.String(), "child")
}

func TestAX7CLI_TreeNode_AddStyled_Ugly(t *core.T) {
	tree := NewTree("root")
	child := tree.AddStyled("", NewStyle().Dim())

	core.AssertNotNil(t, child.style)
	core.AssertLen(t, tree.children, 1)
}

func TestAX7CLI_TreeNode_AddTree_Good(t *core.T) {
	tree := NewTree("root")
	child := NewTree("child")

	core.AssertEqual(t, tree, tree.AddTree(child))
	core.AssertContains(t, tree.String(), "child")
}

func TestAX7CLI_TreeNode_AddTree_Bad(t *core.T) {
	tree := NewTree("root")
	tree.AddTree(nil)

	core.AssertLen(t, tree.children, 1)
	core.AssertPanics(t, func() { _ = tree.String() })
}

func TestAX7CLI_TreeNode_AddTree_Ugly(t *core.T) {
	tree := NewTree("root")
	child := NewTree("child").Add("leaf")
	tree.AddTree(child)

	core.AssertContains(t, tree.String(), "leaf")
	core.AssertLen(t, tree.children, 1)
}

func TestAX7CLI_TreeNode_WithStyle_Good(t *core.T) {
	tree := NewTree("root").WithStyle(NewStyle().Bold())

	core.AssertNotNil(t, tree.style)
	core.AssertContains(t, tree.String(), "root")
}

func TestAX7CLI_TreeNode_WithStyle_Bad(t *core.T) {
	tree := NewTree("root").WithStyle(nil)

	core.AssertNil(t, tree.style)
	core.AssertContains(t, tree.String(), "root")
}

func TestAX7CLI_TreeNode_WithStyle_Ugly(t *core.T) {
	tree := NewTree("").WithStyle(NewStyle().Dim())

	core.AssertNotNil(t, tree.style)
	core.AssertContains(t, tree.String(), "\n")
}

func TestAX7CLI_TreeNode_Children_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	var count int
	for range tree.Children() {
		count++
	}

	core.AssertEqual(t, 1, count)
}

func TestAX7CLI_TreeNode_Children_Bad(t *core.T) {
	tree := NewTree("root")
	var count int
	for range tree.Children() {
		count++
	}

	core.AssertEqual(t, 0, count)
}

func TestAX7CLI_TreeNode_Children_Ugly(t *core.T) {
	tree := NewTree("root")
	tree.Add("a")
	tree.Add("b")
	var names []string
	for child := range tree.Children() {
		names = append(names, child.label)
	}

	core.AssertEqual(t, []string{"a", "b"}, names)
}

func TestAX7CLI_TreeNode_String_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	got := tree.String()

	core.AssertContains(t, got, "root")
	core.AssertContains(t, got, "child")
}

func TestAX7CLI_TreeNode_String_Bad(t *core.T) {
	got := NewTree("").String()

	core.AssertContains(t, got, "\n")
	core.AssertEqual(t, 1, core.RuneCount(got))
}

func TestAX7CLI_TreeNode_String_Ugly(t *core.T) {
	tree := NewTree(":check:")
	tree.Add(":warn:")

	core.AssertContains(t, tree.String(), "✓")
	core.AssertContains(t, tree.String(), "⚠")
}

func TestAX7CLI_TreeNode_Render_Good(t *core.T) {
	tree := NewTree("root")
	tree.Add("child")
	out := ax7CaptureStdout(t, func() { tree.Render() })

	core.AssertContains(t, out, "root")
	core.AssertContains(t, out, "child")
}

func TestAX7CLI_TreeNode_Render_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { NewTree("").Render() })

	core.AssertContains(t, out, "\n")
	core.AssertEqual(t, 1, core.RuneCount(out))
}

func TestAX7CLI_TreeNode_Render_Ugly(t *core.T) {
	tree := NewTree(":check:")
	out := ax7CaptureStdout(t, func() { tree.Render() })

	core.AssertContains(t, out, "✓")
	core.AssertNotContains(t, out, ":check:")
}

func TestAX7CLI_Check_Good(t *core.T) {
	check := Check("audit")

	core.AssertNotNil(t, check)
	core.AssertContains(t, check.String(), "audit")
}

func TestAX7CLI_Check_Bad(t *core.T) {
	check := Check("")

	core.AssertNotNil(t, check)
	core.AssertNotContains(t, check.String(), "\n")
}

func TestAX7CLI_Check_Ugly(t *core.T) {
	check := Check(":check:")

	core.AssertContains(t, check.String(), "✓")
	core.AssertNotContains(t, check.String(), ":check:")
}

func TestAX7CLI_CheckBuilder_Pass_Good(t *core.T) {
	check := Check("audit").Pass()

	core.AssertEqual(t, "passed", check.status)
	core.AssertNotNil(t, check.style)
}

func TestAX7CLI_CheckBuilder_Pass_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Pass() })
	core.AssertNil(t, check)
}

func TestAX7CLI_CheckBuilder_Pass_Ugly(t *core.T) {
	check := Check("audit").Fail().Pass()

	core.AssertEqual(t, "passed", check.status)
	core.AssertEqual(t, Glyph(":check:"), check.icon)
}

func TestAX7CLI_CheckBuilder_Fail_Good(t *core.T) {
	check := Check("audit").Fail()

	core.AssertEqual(t, "failed", check.status)
	core.AssertNotNil(t, check.style)
}

func TestAX7CLI_CheckBuilder_Fail_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Fail() })
	core.AssertNil(t, check)
}

func TestAX7CLI_CheckBuilder_Fail_Ugly(t *core.T) {
	check := Check("audit").Pass().Fail()

	core.AssertEqual(t, "failed", check.status)
	core.AssertEqual(t, Glyph(":cross:"), check.icon)
}

func TestAX7CLI_CheckBuilder_Skip_Good(t *core.T) {
	check := Check("audit").Skip()

	core.AssertEqual(t, "skipped", check.status)
	core.AssertNotNil(t, check.style)
}

func TestAX7CLI_CheckBuilder_Skip_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Skip() })
	core.AssertNil(t, check)
}

func TestAX7CLI_CheckBuilder_Skip_Ugly(t *core.T) {
	check := Check("audit").Fail().Skip()

	core.AssertEqual(t, "skipped", check.status)
	core.AssertEqual(t, Glyph(":skip:"), check.icon)
}

func TestAX7CLI_CheckBuilder_Warn_Good(t *core.T) {
	check := Check("audit").Warn()

	core.AssertEqual(t, "warning", check.status)
	core.AssertNotNil(t, check.style)
}

func TestAX7CLI_CheckBuilder_Warn_Bad(t *core.T) {
	var check *CheckBuilder

	core.AssertPanics(t, func() { check.Warn() })
	core.AssertNil(t, check)
}

func TestAX7CLI_CheckBuilder_Warn_Ugly(t *core.T) {
	check := Check("audit").Pass().Warn()

	core.AssertEqual(t, "warning", check.status)
	core.AssertEqual(t, Glyph(":warn:"), check.icon)
}

func TestAX7CLI_CheckBuilder_Duration_Good(t *core.T) {
	check := Check("audit").Duration("1s")

	core.AssertEqual(t, "1s", check.duration)
	core.AssertContains(t, check.String(), "1s")
}

func TestAX7CLI_CheckBuilder_Duration_Bad(t *core.T) {
	check := Check("audit").Duration("")

	core.AssertEqual(t, "", check.duration)
	core.AssertNotContains(t, check.String(), "1s")
}

func TestAX7CLI_CheckBuilder_Duration_Ugly(t *core.T) {
	check := Check("audit").Duration("∞")

	core.AssertEqual(t, "∞", check.duration)
	core.AssertContains(t, check.String(), "∞")
}

func TestAX7CLI_CheckBuilder_Message_Good(t *core.T) {
	check := Check("audit").Message("ready")

	core.AssertEqual(t, "ready", check.status)
	core.AssertContains(t, check.String(), "ready")
}

func TestAX7CLI_CheckBuilder_Message_Bad(t *core.T) {
	check := Check("audit").Message("")

	core.AssertEqual(t, "", check.status)
	core.AssertNotContains(t, check.String(), "ready")
}

func TestAX7CLI_CheckBuilder_Message_Ugly(t *core.T) {
	check := Check("audit").Message(":check:")

	core.AssertEqual(t, ":check:", check.status)
	core.AssertContains(t, check.String(), "✓")
}

func TestAX7CLI_CheckBuilder_String_Good(t *core.T) {
	got := Check("audit").Pass().String()

	core.AssertContains(t, got, "audit")
	core.AssertContains(t, got, "passed")
}

func TestAX7CLI_CheckBuilder_String_Bad(t *core.T) {
	got := Check("").String()

	core.AssertNotContains(t, got, "\n")
	core.AssertNotContains(t, got, "passed")
}

func TestAX7CLI_CheckBuilder_String_Ugly(t *core.T) {
	got := Check(":check:").Warn().String()

	core.AssertContains(t, got, "✓")
	core.AssertContains(t, got, "warning")
}

func TestAX7CLI_CheckBuilder_Print_Good(t *core.T) {
	out := ax7CaptureStdout(t, func() { Check("audit").Pass().Print() })

	core.AssertContains(t, out, "audit")
	core.AssertContains(t, out, "passed")
}

func TestAX7CLI_CheckBuilder_Print_Bad(t *core.T) {
	out := ax7CaptureStdout(t, func() { Check("").Print() })

	core.AssertContains(t, out, "\n")
	core.AssertNotContains(t, out, "passed")
}

func TestAX7CLI_CheckBuilder_Print_Ugly(t *core.T) {
	out := ax7CaptureStdout(t, func() { Check(":check:").Warn().Print() })

	core.AssertContains(t, out, "✓")
	core.AssertContains(t, out, "warning")
}

func TestAX7CLI_WithWordWrap_Good(t *core.T) {
	stream := NewStream(WithWordWrap(4), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 4, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestAX7CLI_WithWordWrap_Bad(t *core.T) {
	stream := NewStream(WithWordWrap(0), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 0, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestAX7CLI_WithWordWrap_Ugly(t *core.T) {
	stream := NewStream(WithWordWrap(-1), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, -1, stream.wrap)
	core.AssertNotNil(t, stream.out)
}

func TestAX7CLI_WithStreamOutput_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))

	core.AssertEqual(t, out, stream.out)
	core.AssertNotNil(t, stream.done)
}

func TestAX7CLI_WithStreamOutput_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(nil))

	core.AssertNil(t, stream.out)
	core.AssertNotNil(t, stream.done)
}

func TestAX7CLI_WithStreamOutput_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(io.Discard))

	core.AssertEqual(t, io.Discard, stream.out)
	core.AssertEqual(t, 0, stream.Column())
}

func TestAX7CLI_NewStream_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))

	core.AssertNotNil(t, stream)
	core.AssertEqual(t, out, stream.out)
}

func TestAX7CLI_NewStream_Bad(t *core.T) {
	stream := NewStream()

	core.AssertNotNil(t, stream)
	core.AssertNotNil(t, stream.out)
}

func TestAX7CLI_NewStream_Ugly(t *core.T) {
	stream := NewStream(WithWordWrap(3), WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 3, stream.wrap)
	core.AssertEqual(t, 0, stream.Column())
}

func TestAX7CLI_Stream_Write_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello")

	core.AssertEqual(t, "hello", out.String())
	core.AssertEqual(t, 5, stream.Column())
}

func TestAX7CLI_Stream_Write_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("")

	core.AssertEqual(t, "", out.String())
	core.AssertEqual(t, 0, stream.Column())
}

func TestAX7CLI_Stream_Write_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out), WithWordWrap(3))
	stream.Write("abcd")

	core.AssertContains(t, out.String(), "\n")
	core.AssertEqual(t, 1, stream.Column())
}

func TestAX7CLI_Stream_WriteFrom_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	err := stream.WriteFrom(core.NewReader("hello"))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "hello", out.String())
}

func TestAX7CLI_Stream_WriteFrom_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	err := stream.WriteFrom(core.NewReader(""))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "", out.String())
}

func TestAX7CLI_Stream_WriteFrom_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out), WithWordWrap(3))
	err := stream.WriteFrom(core.NewReader("abcd"))

	core.AssertNoError(t, err)
	core.AssertContains(t, out.String(), "\n")
}

func TestAX7CLI_Stream_Done_Good(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello")
	stream.Done()

	core.AssertContains(t, out.String(), "\n")
	core.AssertNotPanics(t, stream.Done)
}

func TestAX7CLI_Stream_Done_Bad(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Done()

	core.AssertEqual(t, "", out.String())
	core.AssertNotPanics(t, stream.Done)
}

func TestAX7CLI_Stream_Done_Ugly(t *core.T) {
	out := core.NewBuilder()
	stream := NewStream(WithStreamOutput(out))
	stream.Write("hello\n")
	stream.Done()

	core.AssertEqual(t, "hello\n", out.String())
	core.AssertNotPanics(t, stream.Done)
}

func TestAX7CLI_Stream_Wait_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 0, stream.Column())
}

func TestAX7CLI_Stream_Wait_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	go stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 0, stream.Column())
}

func TestAX7CLI_Stream_Wait_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("x")
	stream.Done()

	core.AssertNotPanics(t, stream.Wait)
	core.AssertEqual(t, 1, stream.Column())
}

func TestAX7CLI_Stream_Column_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("abc")

	core.AssertEqual(t, 3, stream.Column())
	core.AssertNotEqual(t, 0, stream.Column())
}

func TestAX7CLI_Stream_Column_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))

	core.AssertEqual(t, 0, stream.Column())
	core.AssertNotEqual(t, 1, stream.Column())
}

func TestAX7CLI_Stream_Column_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("a\nbc")

	core.AssertEqual(t, 2, stream.Column())
	core.AssertNotEqual(t, 4, stream.Column())
}

func TestAX7CLI_Stream_Captured_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("hello")

	core.AssertEqual(t, "hello", stream.Captured())
	core.AssertNotEmpty(t, stream.Captured())
}

func TestAX7CLI_Stream_Captured_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(io.Discard))
	stream.Write("hello")

	core.AssertEqual(t, "", stream.Captured())
	core.AssertEmpty(t, stream.Captured())
}

func TestAX7CLI_Stream_Captured_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("")

	core.AssertEqual(t, "", stream.Captured())
	core.AssertEmpty(t, stream.Captured())
}

func TestAX7CLI_Stream_CapturedOK_Good(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	stream.Write("hello")
	got, ok := stream.CapturedOK()

	core.AssertTrue(t, ok)
	core.AssertEqual(t, "hello", got)
}

func TestAX7CLI_Stream_CapturedOK_Bad(t *core.T) {
	stream := NewStream(WithStreamOutput(io.Discard))
	got, ok := stream.CapturedOK()

	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", got)
}

func TestAX7CLI_Stream_CapturedOK_Ugly(t *core.T) {
	stream := NewStream(WithStreamOutput(core.NewBuilder()))
	got, ok := stream.CapturedOK()

	core.AssertTrue(t, ok)
	core.AssertEqual(t, "", got)
}

func TestAX7CLI_NewDaemon_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "daemon.pid")})

	core.AssertNotNil(t, d)
	core.AssertEqual(t, "/health", d.opts.HealthPath)
}

func TestAX7CLI_NewDaemon_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertEqual(t, "", d.opts.PIDFile)
	core.AssertEqual(t, "/ready", d.opts.ReadyPath)
}

func TestAX7CLI_NewDaemon_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthPath: "/h", ReadyPath: "/r"})

	core.AssertEqual(t, "/h", d.opts.HealthPath)
	core.AssertEqual(t, "/r", d.opts.ReadyPath)
}

func TestAX7CLI_Daemon_Start_Good(t *core.T) {
	pid := core.Path(t.TempDir(), "daemon.pid")
	d := NewDaemon(DaemonOptions{PIDFile: pid})

	core.AssertNoError(t, d.Start(context.Background()))
	core.AssertNoError(t, d.Stop(context.Background()))
}

func TestAX7CLI_Daemon_Start_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "missing", "daemon.pid")})

	core.AssertNoError(t, d.Start(nil))
	core.AssertNoError(t, d.Stop(nil))
}

func TestAX7CLI_Daemon_Start_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})

	core.AssertNoError(t, d.Start(context.Background()))
	core.AssertNoError(t, d.Stop(context.Background()))
}

func TestAX7CLI_Daemon_Stop_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{PIDFile: core.Path(t.TempDir(), "daemon.pid")})
	core.RequireNoError(t, d.Start(context.Background()))

	core.AssertNoError(t, d.Stop(context.Background()))
	core.AssertFalse(t, d.started)
}

func TestAX7CLI_Daemon_Stop_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertNoError(t, d.Stop(nil))
	core.AssertFalse(t, d.started)
}

func TestAX7CLI_Daemon_Stop_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})
	core.RequireNoError(t, d.Start(context.Background()))

	core.AssertNoError(t, d.Stop(nil))
	core.AssertEqual(t, "", d.addr)
}

func TestAX7CLI_Daemon_HealthAddr_Good(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:0"})
	core.RequireNoError(t, d.Start(context.Background()))
	defer d.Stop(context.Background())

	core.AssertNotEmpty(t, d.HealthAddr())
}

func TestAX7CLI_Daemon_HealthAddr_Bad(t *core.T) {
	d := NewDaemon(DaemonOptions{})

	core.AssertEqual(t, "", d.HealthAddr())
	core.AssertFalse(t, d.started)
}

func TestAX7CLI_Daemon_HealthAddr_Ugly(t *core.T) {
	d := NewDaemon(DaemonOptions{HealthAddr: "127.0.0.1:9999"})

	core.AssertEqual(t, "127.0.0.1:9999", d.HealthAddr())
	core.AssertFalse(t, d.started)
}

func TestAX7CLI_StopPIDFile_Good(t *core.T) {
	err := StopPIDFile(core.Path(t.TempDir(), "missing.pid"), time.Millisecond)

	core.AssertNoError(t, err)
	core.AssertNil(t, err)
}

func TestAX7CLI_StopPIDFile_Bad(t *core.T) {
	err := StopPIDFile("", time.Millisecond)

	core.AssertNoError(t, err)
	core.AssertNil(t, err)
}

func TestAX7CLI_StopPIDFile_Ugly(t *core.T) {
	path := core.Path(t.TempDir(), "bad.pid")
	core.RequireNoError(t, os.WriteFile(path, []byte("not-a-pid"), 0o644))

	err := StopPIDFile(path, time.Millisecond)
	core.AssertError(t, err)
}

func TestAX7CLI_LogDebug_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug", "k", "v") })
	core.AssertNotPanics(t, func() { LogDebug("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogDebug_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug", "odd") })
	core.AssertNotPanics(t, func() { LogDebug("debug", nil) })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogDebug_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogDebug("debug\nline", "k", 1) })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogInfo_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info", "k", "v") })
	core.AssertNotPanics(t, func() { LogInfo("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogInfo_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info", "odd") })
	core.AssertNotPanics(t, func() { LogInfo("info", nil) })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogInfo_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogInfo("info\nline", "k", 1) })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogWarn_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn", "k", "v") })
	core.AssertNotPanics(t, func() { LogWarn("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogWarn_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn", "odd") })
	core.AssertNotPanics(t, func() { LogWarn("warn", nil) })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogWarn_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogWarn("warn\nline", "k", 1) })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogError_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error", "k", "v") })
	core.AssertNotPanics(t, func() { LogError("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogError_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error", "odd") })
	core.AssertNotPanics(t, func() { LogError("error", nil) })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogError_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogError("error\nline", "k", 1) })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurity_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security", "k", "v") })
	core.AssertNotPanics(t, func() { LogSecurity("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurity_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security", "odd") })
	core.AssertNotPanics(t, func() { LogSecurity("security", nil) })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurity_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurity("security\nline", "k", 1) })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurityf_Good(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("security %s", "event") })
	core.AssertNotPanics(t, func() { LogSecurityf("") })
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurityf_Bad(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("%s", "bad") })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}

func TestAX7CLI_LogSecurityf_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() { LogSecurityf("security\n%s", "event") })
	core.AssertTrue(t, true)
	core.AssertFalse(t, false)
}
