package cli

import (
	"sync"
	"testing"
	"testing/fstest"

	"forge.lthn.ai/core/go-i18n"
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

func resetI18nDefault(t *testing.T) {
	t.Helper()

	prev := i18n.Default()
	svc, err := i18n.New()
	require.NoError(t, err)
	i18n.SetDefault(svc)

	t.Cleanup(func() {
		i18n.SetDefault(prev)
	})
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

	t.Run("nested registration during startup does not deadlock", func(t *testing.T) {
		resetGlobals(t)

		RegisterCommands(func(root *cobra.Command) {
			root.AddCommand(&cobra.Command{Use: "outer", Short: "Outer"})
			RegisterCommands(func(root *cobra.Command) {
				root.AddCommand(&cobra.Command{Use: "inner", Short: "Inner"})
			})
		})

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		for _, name := range []string{"outer", "inner"} {
			cmd, _, err := RootCmd().Find([]string{name})
			require.NoError(t, err)
			assert.Equal(t, name, cmd.Use)
		}
	})
}

// TestLocaleLoading_Good verifies locale files become available to the active i18n service.
func TestLocaleLoading_Good(t *testing.T) {
	t.Run("Init loads I18nSources", func(t *testing.T) {
		resetGlobals(t)
		resetI18nDefault(t)

		localeFS := fstest.MapFS{
			"en.json": {
				Data: []byte(`{"custom":{"hello":"Hello from locale"}}`),
			},
		}

		err := Init(Options{
			AppName:     "test",
			I18nSources: []LocaleSource{WithLocales(localeFS, ".")},
		})
		require.NoError(t, err)

		assert.Equal(t, "Hello from locale", i18n.T("custom.hello"))
	})

	t.Run("WithCommands loads localeFS before registration", func(t *testing.T) {
		resetGlobals(t)
		resetI18nDefault(t)

		err := Init(Options{AppName: "test"})
		require.NoError(t, err)

		localeFS := fstest.MapFS{
			"en.json": {
				Data: []byte(`{"custom":{"immediate":"Loaded eagerly"}}`),
			},
		}

		var observed string
		setup := WithCommands("test", func(root *cobra.Command) {
			_ = root
			observed = i18n.T("custom.immediate")
		}, localeFS)

		setup(Core())

		assert.Equal(t, "Loaded eagerly", observed)
		assert.Equal(t, "Loaded eagerly", i18n.T("custom.immediate"))
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
