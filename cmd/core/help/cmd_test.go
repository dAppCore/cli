package help

import (
	"bytes"
	"io"
	"os"
	"testing"

	"dappco.re/go/cli/pkg/cli"
	gohelp "dappco.re/go/help"
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

func searchableHelpQuery(t *testing.T) string {
	t.Helper()

	catalog := gohelp.DefaultCatalog()
	for _, candidate := range []string{"configuration", "docs", "search", "topic", "help"} {
		if _, err := catalog.Get(candidate); err == nil {
			continue
		}
		if len(catalog.Search(candidate)) > 0 {
			return candidate
		}
	}

	t.Skip("no suitable query found with suggestions")
	return ""
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
	assert.Contains(t, out, "browse")
	assert.Contains(t, out, "core help search <topic>")
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

func TestAddHelpCommands_Good_Search(t *testing.T) {
	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help", "search"})
	require.NoError(t, err)
	require.NotNil(t, cmd)

	query := searchableHelpQuery(t)
	require.NoError(t, cmd.Flags().Set("query", query))

	out := captureOutput(t, func() {
		err := cmd.RunE(cmd, nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "SEARCH RESULTS")
	assert.Contains(t, out, query)
	assert.Contains(t, out, "browse")
	assert.Contains(t, out, "core help search")
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
	assert.Contains(t, out, "browse")
	assert.Contains(t, out, "core help search \"config\"")
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
	assert.Contains(t, out, "browse")
	assert.Contains(t, out, "core help search <topic>")
}

func TestRenderTopic_Good(t *testing.T) {
	out := captureOutput(t, func() {
		renderTopic(&gohelp.Topic{
			ID:      "config",
			Title:   "Configuration",
			Content: "Core is configured via environment variables.",
		})
	})

	assert.Contains(t, out, "Configuration")
	assert.Contains(t, out, "Core is configured via environment variables.")
	assert.Contains(t, out, "browse")
	assert.Contains(t, out, "core help search \"config\"")
}

func TestAddHelpCommands_Bad(t *testing.T) {
	t.Run("missing search results", func(t *testing.T) {
		cmd := newHelpCommand(t)
		require.NoError(t, cmd.Flags().Set("search", "zzzyyyxxx"))

		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "no help topics matched")
		})

		assert.Contains(t, out, "browse")
		assert.Contains(t, out, "core help")
		assert.Contains(t, out, "core help search")
	})

	t.Run("missing topic without suggestions shows hints", func(t *testing.T) {
		cmd := newHelpCommand(t)

		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, []string{"definitely-not-a-real-topic"})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "help topic")
		})

		assert.Contains(t, out, "browse")
		assert.Contains(t, out, "core help")
	})

	t.Run("missing search query", func(t *testing.T) {
		root := &cli.Command{Use: "core"}
		AddHelpCommands(root)

		cmd, _, findErr := root.Find([]string{"help", "search"})
		require.NoError(t, findErr)
		require.NotNil(t, cmd)

		var runErr error
		out := captureOutput(t, func() {
			runErr = cmd.RunE(cmd, nil)
		})
		require.Error(t, runErr)
		assert.Contains(t, runErr.Error(), "help search query is required")
		assert.Contains(t, out, "browse")
		assert.Contains(t, out, "core help")
	})

	t.Run("missing topic shows suggestions when available", func(t *testing.T) {
		query := searchableHelpQuery(t)

		cmd := newHelpCommand(t)
		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, []string{query})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "help topic")
		})

		assert.Contains(t, out, "SEARCH RESULTS")
	})
}
