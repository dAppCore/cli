package cli

import (
	"context"
	core "dappco.re/go"
	"time"
)

func TestRuntime_Init_Good(t *core.T) {
	resetGlobals(t)
	err := cliResultError(Init(Options{AppName: "codex", Version: "1.0.0"}))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "codex", Core().App().Name)
}

func TestRuntime_Init_Bad(t *core.T) {
	resetGlobals(t)
	err := cliResultError(Init(Options{}))

	core.AssertNoError(t, err)
	core.AssertNotNil(t, Core())
}

func TestRuntime_Init_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "once"})))
	err := cliResultError(Init(Options{AppName: "twice"}))

	core.AssertNoError(t, err)
	core.AssertEqual(t, "once", Core().App().Name)
}

func TestRuntime_Core_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "core"})))

	core.AssertNotNil(t, Core())
	core.AssertEqual(t, "core", Core().App().Name)
}

func TestRuntime_Core_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Core() })
	core.AssertNil(t, instance)
}

func TestRuntime_Core_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "core"})))
	Shutdown()

	core.AssertNotNil(t, Core())
	core.AssertNotNil(t, Context())
}

func TestRuntime_Execute_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "execute"})))

	err := cliResultError(Execute())
	core.AssertError(t, err)
}

func TestRuntime_Execute_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Execute() })
	core.AssertNil(t, instance)
}

func TestRuntime_Execute_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "execute"})))
	instance.core.Service("cli", core.Service{})

	err := cliResultError(Execute())
	core.AssertError(t, err)
}

func TestRuntime_Run_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "run"})))

	err := cliResultError(Run(context.Background()))
	core.AssertError(t, err)
}

func TestRuntime_Run_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Run(context.Background()) })
	core.AssertNil(t, instance)
}

func TestRuntime_Run_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "run"})))

	err := cliResultError(Run(nil))
	core.AssertError(t, err)
}

func TestRuntime_RunWithTimeout_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "timeout"})))
	stop := RunWithTimeout(time.Millisecond)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestRuntime_RunWithTimeout_Bad(t *core.T) {
	resetGlobals(t)
	stop := RunWithTimeout(0)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestRuntime_RunWithTimeout_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "timeout"})))
	stop := RunWithTimeout(-time.Second)

	core.AssertNotPanics(t, stop)
	core.AssertNotNil(t, stop)
}

func TestRuntime_Context_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "context"})))

	core.AssertNotNil(t, Context())
	core.AssertNoError(t, Context().Err())
}

func TestRuntime_Context_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertPanics(t, func() { _ = Context() })
	core.AssertNil(t, instance)
}

func TestRuntime_Context_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "context"})))
	Shutdown()

	core.AssertNotNil(t, Context())
	core.AssertError(t, Context().Err())
}

func TestRuntime_Shutdown_Good(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "shutdown"})))

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNotNil(t, instance)
}

func TestRuntime_Shutdown_Bad(t *core.T) {
	resetGlobals(t)

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNil(t, instance)
}

func TestRuntime_Shutdown_Ugly(t *core.T) {
	resetGlobals(t)
	core.RequireNoError(t, cliResultError(Init(Options{AppName: "shutdown"})))
	Shutdown()

	core.AssertNotPanics(t, func() { Shutdown() })
	core.AssertNotNil(t, instance)
}
