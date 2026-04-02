package help

import (
	"bytes"
	"io"
	"os"
	"testing"

	"forge.lthn.ai/core/cli/pkg/cli"
	gohelp "forge.lthn.ai/core/go-help"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldOut := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	defer func() {
		os.Stdout = oldOut
	}()

	fn()

	require.NoError(t, w.Close())

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	return buf.String()
}

func newHelpCommand(t *testing.T) *cli.Command {
	t.Helper()

	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help"})
	require.NoError(t, err)
	return cmd
}

func TestAddHelpCommands_Good(t *testing.T) {
	cmd := newHelpCommand(t)

	topics := gohelp.DefaultCatalog().List()
	require.NotEmpty(t, topics)

	out := captureOutput(t, func() {
		err := cmd.RunE(cmd, nil)
		require.NoError(t, err)
	})
	assert.Contains(t, out, "AVAILABLE HELP TOPICS")
	assert.Contains(t, out, topics[0].ID)
}

func TestAddHelpCommands_Good_Serve(t *testing.T) {
	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help", "serve"})
	require.NoError(t, err)
	require.NotNil(t, cmd)

	oldStart := startHelpServer
	defer func() { startHelpServer = oldStart }()

	var gotAddr string
	startHelpServer = func(catalog *gohelp.Catalog, addr string) error {
		require.NotNil(t, catalog)
		gotAddr = addr
		return nil
	}

	require.NoError(t, cmd.Flags().Set("addr", "127.0.0.1:9090"))
	err = cmd.RunE(cmd, nil)
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1:9090", gotAddr)
}

func TestRenderSearchResults_Good(t *testing.T) {
	out := captureOutput(t, func() {
		err := renderSearchResults([]*gohelp.SearchResult{
			{
				Topic: &gohelp.Topic{
					ID:    "config",
					Title: "Configuration",
				},
				Snippet: "Core is configured via environment variables.",
			},
		}, "config")
		require.NoError(t, err)
	})

	assert.Contains(t, out, "SEARCH RESULTS")
	assert.Contains(t, out, "config - Configuration")
	assert.Contains(t, out, "Core is configured via environment variables.")
}

func TestRenderTopicList_Good(t *testing.T) {
	out := captureOutput(t, func() {
		err := renderTopicList([]*gohelp.Topic{
			{
				ID:      "config",
				Title:   "Configuration",
				Content: "# Configuration\n\nCore is configured via environment variables.\n\nMore details follow.",
			},
		})
		require.NoError(t, err)
	})

	assert.Contains(t, out, "AVAILABLE HELP TOPICS")
	assert.Contains(t, out, "config - Configuration")
	assert.Contains(t, out, "Core is configured via environment variables.")
}

func TestAddHelpCommands_Bad(t *testing.T) {
	t.Run("missing search results", func(t *testing.T) {
		cmd := newHelpCommand(t)
		require.NoError(t, cmd.Flags().Set("search", "zzzyyyxxx"))

		err := cmd.RunE(cmd, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no help topics matched")
	})

	t.Run("missing topic", func(t *testing.T) {
		cmd := newHelpCommand(t)
		err := cmd.RunE(cmd, []string{"definitely-not-a-real-topic"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "help topic")
	})
}
