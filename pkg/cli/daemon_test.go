package cli

import "dappco.re/go"

func TestDetectMode_Good(t *core.T) {
	t.Setenv("CORE_DAEMON", "1")
	mode := DetectMode()
	core.AssertEqual(t, ModeDaemon, mode)
	core.AssertEqual(t, "daemon", mode.String())
}

func TestDetectMode_Bad(t *core.T) {
	t.Setenv("CORE_DAEMON", "0")
	mode := DetectMode()
	core.AssertNotEqual(t, ModeDaemon, mode)
}

func TestDetectMode_Ugly(t *core.T) {
	core.
		// Mode.String() covers all branches including the default unknown case.
		AssertEqual(t, "interactive", ModeInteractive.String())
	core.AssertEqual(t, "pipe", ModePipe.String())
	core.AssertEqual(t, "daemon", ModeDaemon.String())
	core.AssertEqual(t, "unknown", Mode(99).String())
}
