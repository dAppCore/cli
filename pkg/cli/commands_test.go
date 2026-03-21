package cli

import (
	"sync"
	"testing"

	"github.com/spf13/cobra"
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

		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{Use: "hello", Short: "Say hello"})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		// The "hello" command should be on the root.
		cmd, _, err := RootCmd().Find([]string{"hello"})
		require.NoError(t, err)
		assert.Equal(t, "hello", cmd.Use)
	})

	t.Run("multiple groups compose", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{Use: "alpha", Short: "Alpha"})
		})
		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{Use: "beta", Short: "Beta"})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		for _, name := range []string{"alpha", "beta"} {
			cmd, _, err := RootCmd().Find([]string{name})
			require.NoError(t, err)
			assert.Equal(t, name, cmd.Use)
		}
	})

	t.Run("group with subcommands", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(root *cobra.Command) {
			grp := &cobra.Command{Use: "ml", Short: "ML commands"}
			grp.AddCommand(&cobra.Command{Use: "train", Short: "Train a model"})
			grp.AddCommand(&cobra.Command{Use: "serve", Short: "Serve a model"})
			root.AddCommand(grp)
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		cmd, _, err := RootCmd().Find([]string{"ml", "train"})
		require.NoError(t, err)
		assert.Equal(t, "train", cmd.Use)

		cmd, _, err = RootCmd().Find([]string{"ml", "serve"})
		require.NoError(t, err)
		assert.Equal(t, "serve", cmd.Use)
	})

	t.Run("executes registered command", func(t *testing.T) {
		resetGlobals(t)

		executed := false
		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{
				Use:   "ping",
				Short: "Ping",
				RunE: func(_ *cobra.Command, _ []string) error {
					executed = true
					return nil
				},
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		RootCmd().SetArgs([]string{"ping"})
		err = Execute()
		require.NoError(t, err)
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
		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{Use: "late", Short: "Late arrival"})
		})

		cmd, _, err := RootCmd().Find([]string{"late"})
		require.NoError(t, err)
		assert.Equal(t, "late", cmd.Use)
	})
}

// TestWithAppName_Good tests the app name override.
func TestWithAppName_Good(t *testing.T) {
	t.Run("overrides root command use", func(t *testing.T) {
		resetGlobals(t)

		WithAppName("lem")
		defer WithAppName("core") // restore

		err := Init(Options{AppName: AppName})
		require.NoError(t, err)

		assert.Equal(t, "lem", RootCmd().Use)
	})

	t.Run("default is core", func(t *testing.T) {
		resetGlobals(t)

		err := Init(Options{AppName: AppName})
		require.NoError(t, err)

		assert.Equal(t, "core", RootCmd().Use)
	})
}

