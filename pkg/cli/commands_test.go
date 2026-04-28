package cli

import (
	"sync"

	"dappco.re/go"
)

// resetGlobals clears the CLI singleton and command registry for test isolation.
func resetGlobals(t *core.T) {
	t.Helper()
	doReset()
	t.Cleanup(doReset)
}

// doReset clears all package-level state. Only safe from a single goroutine
// with no concurrent RegisterCommands calls in flight (i.e. test setup/teardown).
func doReset() {
	registeredCommandsMu.Lock()
	registeredCommands = nil
	registeredLocales = nil
	commandsAttached = false
	registeredCommandsMu.Unlock()
	if instance != nil {
		Shutdown()
	}
	instance = nil
	once = sync.Once{}
}

// TestRegisterCommands_Good tests the happy path for command registration.
func TestRegisterCommands_Good(t *core.T) {
	t.Run("registers on startup", func(t *core.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("hello", core.Command{
				Description: "Say hello",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
		})

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)

		// The "hello" command should be registered.
		r := Core().Command("hello")
		core.AssertTrue(t, r.OK, "hello command should be registered")
	})

	t.Run("multiple groups compose", func(t *core.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("alpha", core.Command{
				Description: "Alpha",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
		})
		RegisterCommands(func(c *core.Core) {
			c.Command("beta", core.Command{
				Description: "Beta",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
		})

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)

		for _, name := range []string{"alpha", "beta"} {
			r := Core().Command(name)
			core.AssertTrue(t, r.OK, name+" command should be registered")
		}
	})

	t.Run("nested commands via path", func(t *core.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("ml/train", core.Command{
				Description: "Train a model",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
			c.Command("ml/serve", core.Command{
				Description: "Serve a model",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
		})

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)

		r := Core().Command("ml/train")
		core.AssertTrue(t, r.OK, "ml/train command should be registered")

		r = Core().Command("ml/serve")
		core.AssertTrue(t, r.OK, "ml/serve command should be registered")
	})

	t.Run("executes registered command", func(t *core.T) {
		resetGlobals(t)

		executed := false
		RegisterCommands(func(c *core.Core) {
			c.Command("ping", core.Command{
				Description: "Ping",
				Action: func(_ core.Options) core.Result {
					executed = true
					return core.Ok(nil)
				},
			})
		})

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)

		cl := Core().Cli()
		core.RequireTrue(t, cl != nil, "RequireNotNil")
		result := cl.Run("ping")
		core.AssertTrue(t, result.OK, "ping command should execute successfully")
		core.AssertTrue(t, executed, "registered command should have been executed")
	})
}

// TestRegisterCommands_Bad tests expected error conditions.
func TestRegisterCommands_Bad(t *core.T) {
	t.Run("late registration attaches immediately", func(t *core.T) {
		resetGlobals(t)

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)

		// Register after Init — should attach immediately.
		RegisterCommands(func(c *core.Core) {
			c.Command("late", core.Command{
				Description: "Late arrival",
				Action: func(_ core.Options) core.Result {
					return core.Ok(nil)
				},
			})
		})

		r := Core().Command("late")
		core.AssertTrue(t, r.OK, "late command should be registered")
	})
}

// TestWithAppName_Good tests the app name override.
func TestWithAppName_Good(t *core.T) {
	t.Run("overrides app name", func(t *core.T) {
		resetGlobals(t)

		WithAppName("lem")
		defer WithAppName("core") // restore

		err := Init(Options{AppName: AppName})
		core.RequireNoError(t, err)
		core.AssertEqual(t, "lem", Core().App().Name)
	})

	t.Run("default is core", func(t *core.T) {
		resetGlobals(t)

		err := Init(Options{AppName: AppName})
		core.RequireNoError(t, err)
		core.AssertEqual(t, "core", Core().App().Name)
	})
}

// TestRegisterCommands_Ugly tests edge cases and concurrent registration.
func TestRegisterCommands_Ugly(t *core.T) {
	t.Run("register nil function does not panic", func(t *core.T) {
		resetGlobals(t)
		core.

			// Registering a nil function should not panic at registration time.
			AssertNotPanics(t, func() {
				RegisterCommands(nil)
			})
	})

	t.Run("re-init after shutdown is idempotent", func(t *core.T) {
		resetGlobals(t)

		err := Init(Options{AppName: "test"})
		core.RequireNoError(t, err)
		Shutdown()

		resetGlobals(t)
		err = Init(Options{AppName: "test"})
		core.RequireNoError(t, err)
		core.AssertNotNil(t, Core())
	})
}
