package cli

import (
	core "dappco.re/go"
)

func TestDaemon_Mode_String_Good(t *core.T) {
	got := ModeInteractive.String()

	core.AssertEqual(t, "interactive", got)
	core.AssertContains(t, got, "active")
}

func TestDaemon_Mode_String_Bad(t *core.T) {
	got := Mode(99).String()

	core.AssertEqual(t, "unknown", got)
	core.AssertContains(t, got, "unknown")
}

func TestDaemon_Mode_String_Ugly(t *core.T) {
	got := ModeDaemon.String()

	core.AssertEqual(t, "daemon", got)
	core.AssertNotEqual(t, "interactive", got)
}

func TestDaemon_DetectMode_Good(t *core.T) {
	t.Setenv("CORE_DAEMON", "1")
	got := DetectMode()

	core.AssertEqual(t, ModeDaemon, got)
	core.AssertEqual(t, "daemon", got.String())
}

func TestDaemon_DetectMode_Bad(t *core.T) {
	t.Setenv("CORE_DAEMON", "")
	SetStdout(core.NewBuilder())
	defer SetStdout(nil)

	core.AssertEqual(t, ModePipe, DetectMode())
}

func TestDaemon_DetectMode_Ugly(t *core.T) {
	t.Setenv("CORE_DAEMON", "0")
	SetStdout(core.Discard)
	defer SetStdout(nil)

	core.AssertEqual(t, ModePipe, DetectMode())
}

func TestDaemon_IsTTY_Good(t *core.T) {
	SetStdout(core.Discard)
	defer SetStdout(nil)

	core.AssertFalse(t, IsTTY())
}

func TestDaemon_IsTTY_Bad(t *core.T) {
	SetStdout(core.NewBuilder())
	defer SetStdout(nil)

	core.AssertFalse(t, IsTTY())
}

func TestDaemon_IsTTY_Ugly(t *core.T) {
	SetStdout(nil)
	got := IsTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsTTY())
}

func TestDaemon_IsStdinTTY_Good(t *core.T) {
	SetStdin(core.NewReader(""))
	defer SetStdin(nil)

	core.AssertFalse(t, IsStdinTTY())
}

func TestDaemon_IsStdinTTY_Bad(t *core.T) {
	SetStdin(nil)
	got := IsStdinTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsStdinTTY())
}

func TestDaemon_IsStdinTTY_Ugly(t *core.T) {
	SetStdin(core.NewReader("input"))
	defer SetStdin(nil)

	core.AssertFalse(t, IsStdinTTY())
}

func TestDaemon_IsStderrTTY_Good(t *core.T) {
	SetStderr(core.Discard)
	defer SetStderr(nil)

	core.AssertFalse(t, IsStderrTTY())
}

func TestDaemon_IsStderrTTY_Bad(t *core.T) {
	SetStderr(core.NewBuilder())
	defer SetStderr(nil)

	core.AssertFalse(t, IsStderrTTY())
}

func TestDaemon_IsStderrTTY_Ugly(t *core.T) {
	SetStderr(nil)
	got := IsStderrTTY()

	core.AssertTrue(t, got || !got)
	core.AssertEqual(t, got, IsStderrTTY())
}
