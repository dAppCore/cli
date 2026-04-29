package cli

import (
	core "dappco.re/go"
)

func TestErrors_Err_Good(t *core.T) {
	err := cliResultError(Err("missing %s", "config"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "missing config")
}

func TestErrors_Err_Bad(t *core.T) {
	err := cliResultError(Err(""))

	core.AssertError(t, err)
	core.AssertEqual(t, "cli: ", err.Error())
}

func TestErrors_Err_Ugly(t *core.T) {
	err := cliResultError(Err("line\n%s", "two"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "line\ntwo")
}

func TestErrors_Wrap_Good(t *core.T) {
	base := core.NewError("connection refused")
	err := cliResultError(Wrap(base, "connect"))

	core.AssertError(t, err)
	core.AssertTrue(t, Is(err, base))
}

func TestErrors_Wrap_Bad(t *core.T) {
	err := cliResultError(Wrap(nil, "connect"))

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestErrors_Wrap_Ugly(t *core.T) {
	base := cliResultError(Err("root"))
	err := cliResultError(Wrap(base, ""))

	core.AssertError(t, err)
	core.AssertTrue(t, Is(err, base))
}

func TestErrors_WrapVerb_Good(t *core.T) {
	base := core.NewError("denied")
	err := cliResultError(WrapVerb(base, "load", "config"))

	core.AssertContains(t, err.Error(), "Failed to load config")
	core.AssertTrue(t, Is(err, base))
}

func TestErrors_WrapVerb_Bad(t *core.T) {
	err := cliResultError(WrapVerb(nil, "load", "config"))

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestErrors_WrapVerb_Ugly(t *core.T) {
	err := cliResultError(WrapVerb(core.NewError("denied"), "", "config"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "denied")
}

func TestErrors_WrapAction_Good(t *core.T) {
	base := core.NewError("offline")
	err := cliResultError(WrapAction(base, "connect"))

	core.AssertContains(t, err.Error(), "Failed to connect")
	core.AssertTrue(t, Is(err, base))
}

func TestErrors_WrapAction_Bad(t *core.T) {
	err := cliResultError(WrapAction(nil, "connect"))

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestErrors_WrapAction_Ugly(t *core.T) {
	err := cliResultError(WrapAction(core.NewError("offline"), ""))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "offline")
}

func TestErrors_Is_Good(t *core.T) {
	target := core.NewError("target")
	err := cliResultError(Wrap(target, "wrap"))

	core.AssertTrue(t, Is(err, target))
	core.AssertFalse(t, Is(err, core.NewError("other")))
}

func TestErrors_Is_Bad(t *core.T) {
	left := core.NewError("left")
	right := core.NewError("right")

	core.AssertFalse(t, Is(left, right))
	core.AssertFalse(t, Is(nil, right))
}

func TestErrors_Is_Ugly(t *core.T) {
	got := Is(nil, nil)
	core.AssertTrue(t, got)
	core.AssertTrue(t, Is(nil, nil))
	core.AssertFalse(t, Is(core.NewError("x"), nil))
}

func TestErrors_As_Good(t *core.T) {
	err := cliResultError(Exit(7, cliResultError(Err("exit"))))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 7, exitErr.Code)
}

func TestErrors_As_Bad(t *core.T) {
	var exitErr *ExitError
	err := core.NewError("plain")

	core.AssertFalse(t, As(err, &exitErr))
	core.AssertNil(t, exitErr)
}

func TestErrors_As_Ugly(t *core.T) {
	var exitErr *ExitError

	core.AssertFalse(t, As(nil, &exitErr))
	core.AssertNil(t, exitErr)
}

func TestErrors_Join_Good(t *core.T) {
	first := cliResultError(Err("first"))
	second := cliResultError(Err("second"))
	err := cliResultError(Join(first, second))

	core.AssertTrue(t, Is(err, first))
	core.AssertTrue(t, Is(err, second))
}

func TestErrors_Join_Bad(t *core.T) {
	err := cliResultError(Join(nil, nil))

	core.AssertNil(t, err)
	core.AssertFalse(t, Is(err, core.NewError("x")))
}

func TestErrors_Join_Ugly(t *core.T) {
	err := cliResultError(Join(nil, cliResultError(Err("only"))))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "only")
}

func TestErrors_ExitError_Error_Good(t *core.T) {
	err := &ExitError{Code: 2, Err: cliResultError(Err("failed"))}

	core.AssertEqual(t, "cli: failed", err.Error())
	core.AssertEqual(t, 2, err.Code)
}

func TestErrors_ExitError_Error_Bad(t *core.T) {
	err := &ExitError{Code: 2}

	core.AssertEqual(t, "", err.Error())
	core.AssertNil(t, err.Err)
}

func TestErrors_ExitError_Error_Ugly(t *core.T) {
	err := &ExitError{Code: 255, Err: cliResultError(Err("line\nfail"))}

	core.AssertContains(t, err.Error(), "line\nfail")
	core.AssertEqual(t, 255, err.Code)
}

func TestErrors_Exit_Good(t *core.T) {
	err := cliResultError(Exit(2, cliResultError(Err("bad args"))))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 2, exitErr.Code)
}

func TestErrors_Exit_Bad(t *core.T) {
	err := cliResultError(Exit(2, nil))

	core.AssertNil(t, err)
	core.AssertFalse(t, As(err, new(*ExitError)))
}

func TestErrors_Exit_Ugly(t *core.T) {
	err := cliResultError(Exit(0, cliResultError(Err("zero"))))
	var exitErr *ExitError

	core.AssertTrue(t, As(err, &exitErr))
	core.AssertEqual(t, 0, exitErr.Code)
}

func TestErrors_Fatal_Good(t *core.T) {
	if core.Getenv("AX7_FATAL_GOOD") == "1" {
		SetStderr(core.Discard)
		Fatal(cliResultError(Err("fatal")))
		return
	}
	err := cliRunSelf(t, "AX7_FATAL_GOOD")
	core.AssertError(t, err)
}

func TestErrors_Fatal_Bad(t *core.T) {
	cliPlainCLI(t)

	core.AssertNotPanics(t, func() { Fatal(nil) })
	core.AssertEqual(t, "", cliCaptureStderr(t, func() { Fatal(nil) }))
}

func TestErrors_Fatal_Ugly(t *core.T) {
	if core.Getenv("AX7_FATAL_UGLY") == "1" {
		SetStderr(core.Discard)
		Fatal(cliResultError(Err("fatal\nline")))
		return
	}
	err := cliRunSelf(t, "AX7_FATAL_UGLY")
	core.AssertError(t, err)
}

func TestErrors_Fatalf_Good(t *core.T) {
	if core.Getenv("AX7_FATALF_GOOD") == "1" {
		SetStderr(core.Discard)
		Fatalf("fatal %s", "format")
		return
	}
	err := cliRunSelf(t, "AX7_FATALF_GOOD")
	core.AssertError(t, err)
}

func TestErrors_Fatalf_Bad(t *core.T) {
	if core.Getenv("AX7_FATALF_BAD") == "1" {
		SetStderr(core.Discard)
		Fatalf("")
		return
	}
	err := cliRunSelf(t, "AX7_FATALF_BAD")
	core.AssertError(t, err)
}

func TestErrors_Fatalf_Ugly(t *core.T) {
	if core.Getenv("AX7_FATALF_UGLY") == "1" {
		SetStderr(core.Discard)
		Fatalf("fatal %d", 42)
		return
	}
	err := cliRunSelf(t, "AX7_FATALF_UGLY")
	core.AssertError(t, err)
}

func TestErrors_FatalWrap_Good(t *core.T) {
	if core.Getenv("AX7_FATALWRAP_GOOD") == "1" {
		SetStderr(core.Discard)
		FatalWrap(cliResultError(Err("root")), "wrap")
		return
	}
	err := cliRunSelf(t, "AX7_FATALWRAP_GOOD")
	core.AssertError(t, err)
}

func TestErrors_FatalWrap_Bad(t *core.T) {
	cliPlainCLI(t)

	core.AssertNotPanics(t, func() { FatalWrap(nil, "wrap") })
	core.AssertEqual(t, "", cliCaptureStderr(t, func() { FatalWrap(nil, "wrap") }))
}

func TestErrors_FatalWrap_Ugly(t *core.T) {
	if core.Getenv("AX7_FATALWRAP_UGLY") == "1" {
		SetStderr(core.Discard)
		FatalWrap(cliResultError(Err("root")), "")
		return
	}
	err := cliRunSelf(t, "AX7_FATALWRAP_UGLY")
	core.AssertError(t, err)
}

func TestErrors_FatalWrapVerb_Good(t *core.T) {
	if core.Getenv("AX7_FATALWRAPVERB_GOOD") == "1" {
		SetStderr(core.Discard)
		FatalWrapVerb(cliResultError(Err("root")), "load", "config")
		return
	}
	err := cliRunSelf(t, "AX7_FATALWRAPVERB_GOOD")
	core.AssertError(t, err)
}

func TestErrors_FatalWrapVerb_Bad(t *core.T) {
	cliPlainCLI(t)

	core.AssertNotPanics(t, func() { FatalWrapVerb(nil, "load", "config") })
	core.AssertEqual(t, "", cliCaptureStderr(t, func() { FatalWrapVerb(nil, "load", "config") }))
}

func TestErrors_FatalWrapVerb_Ugly(t *core.T) {
	if core.Getenv("AX7_FATALWRAPVERB_UGLY") == "1" {
		SetStderr(core.Discard)
		FatalWrapVerb(cliResultError(Err("root")), "", "")
		return
	}
	err := cliRunSelf(t, "AX7_FATALWRAPVERB_UGLY")
	core.AssertError(t, err)
}
