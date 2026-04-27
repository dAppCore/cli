package cli

import (
	"sync"
	"testing"

	"dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetGlobals clears the CLI singleton and command registry for test isolation.
func resetGlobals(t *testing.T) {
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
func TestRegisterCommands_Good(t *testing.T) {
	t.Run("registers on startup", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("hello", core.Command{
				Description: "Say hello",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		// The "hello" command should be registered.
		r := Core().Command("hello")
		assert.True(t, r.OK, "hello command should be registered")
	})

	t.Run("multiple groups compose", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("alpha", core.Command{
				Description: "Alpha",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
		})
		RegisterCommands(func(c *core.Core) {
			c.Command("beta", core.Command{
				Description: "Beta",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		for _, name := range []string{"alpha", "beta"} {
			r := Core().Command(name)
			assert.True(t, r.OK, name+" command should be registered")
		}
	})

	t.Run("nested commands via path", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(c *core.Core) {
			c.Command("ml/train", core.Command{
				Description: "Train a model",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
			c.Command("ml/serve", core.Command{
				Description: "Serve a model",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		r := Core().Command("ml/train")
		assert.True(t, r.OK, "ml/train command should be registered")

		r = Core().Command("ml/serve")
		assert.True(t, r.OK, "ml/serve command should be registered")
	})

	t.Run("executes registered command", func(t *testing.T) {
		resetGlobals(t)

		executed := false
		RegisterCommands(func(c *core.Core) {
			c.Command("ping", core.Command{
				Description: "Ping",
				Action: func(_ core.Options) core.Result {
					executed = true
					return core.Result{OK: true}
				},
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		cl := Core().Cli()
		require.NotNil(t, cl)
		result := cl.Run("ping")
		assert.True(t, result.OK, "ping command should execute successfully")
		assert.True(t, executed, "registered command should have been executed")
	})
}

// TestRegisterCommands_Bad tests expected error conditions.
func TestRegisterCommands_Bad(t *testing.T) {
	t.Run("late registration attaches immediately", func(t *testing.T) {
		resetGlobals(t)

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		// Register after Init — should attach immediately.
		RegisterCommands(func(c *core.Core) {
			c.Command("late", core.Command{
				Description: "Late arrival",
				Action: func(_ core.Options) core.Result {
					return core.Result{OK: true}
				},
			})
		})

		r := Core().Command("late")
		assert.True(t, r.OK, "late command should be registered")
	})
}

// TestWithAppName_Good tests the app name override.
func TestWithAppName_Good(t *testing.T) {
	t.Run("overrides app name", func(t *testing.T) {
		resetGlobals(t)

		WithAppName("lem")
		defer WithAppName("core") // restore

		err := Init(Options{AppName: AppName})
		require.NoError(t, err)

		assert.Equal(t, "lem", Core().App().Name)
	})

	t.Run("default is core", func(t *testing.T) {
		resetGlobals(t)

		err := Init(Options{AppName: AppName})
		require.NoError(t, err)

		assert.Equal(t, "core", Core().App().Name)
	})
}

// TestRegisterCommands_Ugly tests edge cases and concurrent registration.
func TestRegisterCommands_Ugly(t *testing.T) {
	t.Run("register nil function does not panic", func(t *testing.T) {
		resetGlobals(t)

		// Registering a nil function should not panic at registration time.
		assert.NotPanics(t, func() {
			RegisterCommands(nil)
		})
	})

	t.Run("re-init after shutdown is idempotent", func(t *testing.T) {
		resetGlobals(t)

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)
		Shutdown()

		resetGlobals(t)
		err = Init(Options{AppName: "test"})
		require.NoError(t, err)
		assert.NotNil(t, Core())
	})
}
