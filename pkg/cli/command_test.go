package cli

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistentFlagHelpers_Good(t *testing.T) {
	t.Run("persistent flags inherit through subcommands", func(t *testing.T) {
		parent := NewGroup("parent", "Parent", "")

		var (
			str   string
			b     bool
			i     int
			i64   int64
			f64   float64
			dur   time.Duration
			slice []string
		)

		PersistentStringFlag(parent, &str, "name", "n", "default", "Name")
		PersistentBoolFlag(parent, &b, "debug", "d", false, "Debug")
		PersistentIntFlag(parent, &i, "count", "c", 1, "Count")
		PersistentInt64Flag(parent, &i64, "seed", "", 2, "Seed")
		PersistentFloat64Flag(parent, &f64, "ratio", "", 3.5, "Ratio")
		PersistentDurationFlag(parent, &dur, "timeout", "t", 4*time.Second, "Timeout")
		PersistentStringSliceFlag(parent, &slice, "tag", "", nil, "Tags")

		child := NewCommand("child", "Child", "", func(_ *Command, _ []string) error {
			assert.Equal(t, "override", str)
			assert.True(t, b)
			assert.Equal(t, 9, i)
			assert.Equal(t, int64(42), i64)
			assert.InDelta(t, 7.25, f64, 1e-9)
			assert.Equal(t, 15*time.Second, dur)
			assert.Equal(t, []string{"alpha", "beta"}, slice)
			return nil
		})
		parent.AddCommand(child)

		parent.SetArgs([]string{
			"child",
			"--name", "override",
			"--debug",
			"--count", "9",
			"--seed", "42",
			"--ratio", "7.25",
			"--timeout", "15s",
			"--tag", "alpha",
			"--tag", "beta",
		})

		require.NoError(t, parent.Execute())
	})

	t.Run("persistent string array flags inherit through subcommands", func(t *testing.T) {
		parent := NewGroup("parent", "Parent", "")

		var tags []string
		PersistentStringArrayFlag(parent, &tags, "tag", "t", nil, "Tags")

		child := NewCommand("child", "Child", "", func(_ *Command, _ []string) error {
			assert.Equal(t, []string{"alpha", "beta"}, tags)
			return nil
		})
		parent.AddCommand(child)
		parent.SetArgs([]string{"child", "--tag", "alpha", "-t", "beta"})

		require.NoError(t, parent.Execute())
	})

	t.Run("persistent helpers use short flags when provided", func(t *testing.T) {
		parent := NewGroup("parent", "Parent", "")
		var value int
		PersistentIntFlag(parent, &value, "count", "c", 1, "Count")

		var seen bool
		child := &cobra.Command{
			Use: "child",
			RunE: func(_ *cobra.Command, _ []string) error {
				seen = true
				assert.Equal(t, 5, value)
				return nil
			},
		}
		parent.AddCommand(child)
		parent.SetArgs([]string{"child", "-c", "5"})

		require.NoError(t, parent.Execute())
		assert.True(t, seen)
	})
}

func TestFlagHelpers_Good(t *testing.T) {
	t.Run("string array flags collect repeated values", func(t *testing.T) {
		cmd := NewCommand("child", "Child", "", func(_ *Command, _ []string) error {
			return nil
		})

		var tags []string
		StringArrayFlag(cmd, &tags, "tag", "t", nil, "Tags")
		cmd.SetArgs([]string{"--tag", "alpha", "-t", "beta"})

		require.NoError(t, cmd.Execute())
		assert.Equal(t, []string{"alpha", "beta"}, tags)
	})

	t.Run("string array flags use short flags when provided", func(t *testing.T) {
		cmd := NewCommand("child", "Child", "", func(_ *Command, _ []string) error {
			return nil
		})

		var tags []string
		StringArrayFlag(cmd, &tags, "tag", "t", nil, "Tags")
		cmd.SetArgs([]string{"-t", "alpha"})

		require.NoError(t, cmd.Execute())
		assert.Equal(t, []string{"alpha"}, tags)
	})
}
