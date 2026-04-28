package help

import (
	"bytes"
	"io"
	"os"

	. "dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	gohelp "dappco.re/go/help"
)

func captureOutput(t *T, fn func()) string {
	t.Helper()

	oldOut := os.Stdout
	r, w, err := os.Pipe()
	RequireNoError(t, err)
	os.Stdout = w

	defer func() {
		os.Stdout = oldOut
	}()

	fn()
	RequireNoError(t, w.Close())

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	RequireNoError(t, err)
	return buf.String()
}

func newHelpCommand(t *T) *cli.Command {
	t.Helper()

	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help"})
	RequireNoError(t, err)
	return cmd
}

func searchableHelpQuery(t *T) string {
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

func TestAddHelpCommands_Good(t *T) {
	cmd := newHelpCommand(t)

	topics := gohelp.DefaultCatalog().List()
	RequireNotEmpty(t, topics)

	out := captureOutput(t, func() {
		err := cmd.RunE(cmd, nil)
		RequireNoError(t, err)
	})
	AssertContains(t, out, "AVAILABLE HELP TOPICS")
	AssertContains(t, out, topics[0].ID)
	AssertContains(t, out, "browse")
	AssertContains(t, out, "core help search <topic>")
}

func TestAddHelpCommands_Good_Serve(t *T) {
	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help", "serve"})
	RequireNoError(t, err)
	RequireTrue(t, cmd != nil, "RequireNotNil")

	oldStart := startHelpServer
	defer func() { startHelpServer = oldStart }()

	var gotAddr string
	startHelpServer = func(catalog *gohelp.Catalog, addr string) error {
		RequireTrue(t, catalog != nil, "RequireNotNil")
		gotAddr = addr
		return nil
	}
	RequireNoError(t, cmd.Flags().Set("addr", "127.0.0.1:9090"))
	err = cmd.RunE(cmd, nil)
	RequireNoError(t, err)
	AssertEqual(t, "127.0.0.1:9090", gotAddr)
}

func TestAddHelpCommands_Good_Search(t *T) {
	root := &cli.Command{Use: "core"}
	AddHelpCommands(root)

	cmd, _, err := root.Find([]string{"help", "search"})
	RequireNoError(t, err)
	RequireTrue(t, cmd != nil, "RequireNotNil")

	query := searchableHelpQuery(t)
	RequireNoError(t, cmd.Flags().Set("query", query))

	out := captureOutput(t, func() {
		err := cmd.RunE(cmd, nil)
		RequireNoError(t, err)
	})
	AssertContains(t, out, "SEARCH RESULTS")
	AssertContains(t, out, query)
	AssertContains(t, out, "browse")
	AssertContains(t, out, "core help search")
}

func TestRenderSearchResults_Good(t *T) {
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
		RequireNoError(t, err)
	})
	AssertContains(t, out, "SEARCH RESULTS")
	AssertContains(t, out, "config - Configuration")
	AssertContains(t, out, "Core is configured via environment variables.")
	AssertContains(t, out, "browse")
	AssertContains(t, out, "core help search \"config\"")
}

func TestRenderTopicList_Good(t *T) {
	out := captureOutput(t, func() {
		err := renderTopicList([]*gohelp.Topic{
			{
				ID:      "config",
				Title:   "Configuration",
				Content: "# Configuration\n\nCore is configured via environment variables.\n\nMore details follow.",
			},
		})
		RequireNoError(t, err)
	})
	AssertContains(t, out, "AVAILABLE HELP TOPICS")
	AssertContains(t, out, "config - Configuration")
	AssertContains(t, out, "Core is configured via environment variables.")
	AssertContains(t, out, "browse")
	AssertContains(t, out, "core help search <topic>")
}

func TestRenderTopic_Good(t *T) {
	out := captureOutput(t, func() {
		renderTopic(&gohelp.Topic{
			ID:      "config",
			Title:   "Configuration",
			Content: "Core is configured via environment variables.",
		})
	})
	AssertContains(t, out, "Configuration")
	AssertContains(t, out, "Core is configured via environment variables.")
	AssertContains(t, out, "browse")
	AssertContains(t, out, "core help search \"config\"")
}

func TestAddHelpCommands_Bad(t *T) {
	t.Run("missing search results", func(t *T) {
		cmd := newHelpCommand(t)
		RequireNoError(t, cmd.Flags().Set("search", "zzzyyyxxx"))

		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, nil)
			RequireTrue(t, err != nil, "RequireError")
			AssertContains(t, err.Error(), "no help topics matched")
		})
		AssertContains(t, out, "browse")
		AssertContains(t, out, "core help")
		AssertContains(t, out, "core help search")
	})

	t.Run("missing topic without suggestions shows hints", func(t *T) {
		cmd := newHelpCommand(t)

		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, []string{"definitely-not-a-real-topic"})
			RequireTrue(t, err != nil, "RequireError")
			AssertContains(t, err.Error(), "help topic")
		})
		AssertContains(t, out, "browse")
		AssertContains(t, out, "core help")
	})

	t.Run("missing search query", func(t *T) {
		root := &cli.Command{Use: "core"}
		AddHelpCommands(root)

		cmd, _, findErr := root.Find([]string{"help", "search"})
		RequireNoError(t, findErr)
		RequireTrue(t, cmd != nil, "RequireNotNil")

		var runErr error
		out := captureOutput(t, func() {
			runErr = cmd.RunE(cmd, nil)
		})
		RequireTrue(t, runErr != nil, "RequireError")
		AssertContains(t, runErr.Error(), "help search query is required")
		AssertContains(t, out, "browse")
		AssertContains(t, out, "core help")
	})

	t.Run("missing topic shows suggestions when available", func(t *T) {
		query := searchableHelpQuery(t)

		cmd := newHelpCommand(t)
		out := captureOutput(t, func() {
			err := cmd.RunE(cmd, []string{query})
			RequireTrue(t, err != nil, "RequireError")
			AssertContains(t, err.Error(), "help topic")
		})
		AssertContains(t, out, "SEARCH RESULTS")
	})
}
