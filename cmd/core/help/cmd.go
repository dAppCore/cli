package help

import (
	"bufio"
	"fmt"
	"strings"

	"forge.lthn.ai/core/cli/pkg/cli"
	gohelp "forge.lthn.ai/core/go-help"
	"github.com/spf13/cobra"
)

var startHelpServer = func(catalog *gohelp.Catalog, addr string) error {
	return gohelp.NewServer(catalog, addr).ListenAndServe()
}

func AddHelpCommands(root *cli.Command) {
	var searchQuery string

	helpCmd := &cli.Command{
		Use:   "help [topic]",
		Short: "Display help documentation",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cli.Command, args []string) error {
			catalog := gohelp.DefaultCatalog()

			if searchQuery != "" {
				return searchHelpTopics(catalog, searchQuery)
			}

			if len(args) == 0 {
				return renderTopicList(catalog.List())
			}

			topic, err := catalog.Get(args[0])
			if err != nil {
				if suggestions := catalog.Search(args[0]); len(suggestions) > 0 {
					if suggestErr := renderSearchResults(suggestions, args[0]); suggestErr != nil {
						return suggestErr
					}
					cli.Blank()
					renderHelpHint(args[0])
					return cli.Err("help topic %q not found", args[0])
				}
				renderHelpHint(args[0])
				return cli.Err("help topic %q not found", args[0])
			}

			renderTopic(topic)
			return nil
		},
	}

	searchCmd := &cli.Command{
		Use:   "search [query]",
		Short: "Search help topics",
		Args:  cobra.ArbitraryArgs,
	}
	var searchCmdQuery string
	searchCmd.Flags().StringVarP(&searchCmdQuery, "query", "q", "", "Search query")
	searchCmd.RunE = func(cmd *cli.Command, args []string) error {
		catalog := gohelp.DefaultCatalog()
		query := strings.TrimSpace(searchCmdQuery)
		if query == "" {
			query = strings.TrimSpace(strings.Join(args, " "))
		}
		if query == "" {
			return cli.Err("help search query is required")
		}
		return searchHelpTopics(catalog, query)
	}

	var serveAddr string
	serveCmd := &cli.Command{
		Use:   "serve",
		Short: "Serve help documentation over HTTP",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cli.Command, args []string) error {
			return startHelpServer(gohelp.DefaultCatalog(), serveAddr)
		},
	}
	serveCmd.Flags().StringVar(&serveAddr, "addr", ":8080", "HTTP listen address")

	helpCmd.AddCommand(serveCmd)
	helpCmd.AddCommand(searchCmd)
	helpCmd.Flags().StringVarP(&searchQuery, "search", "s", "", "Search help topics")
	root.AddCommand(helpCmd)
}

func searchHelpTopics(catalog *gohelp.Catalog, query string) error {
	return renderSearchResults(catalog.Search(query), query)
}

func renderSearchResults(results []*gohelp.SearchResult, query string) error {
	if len(results) == 0 {
		renderHelpHint(query)
		return cli.Err("no help topics matched %q", query)
	}

	cli.Section("Search Results")
	for _, res := range results {
		cli.Println("  %s - %s", res.Topic.ID, res.Topic.Title)
		if snippet := strings.TrimSpace(res.Snippet); snippet != "" {
			cli.Println("%s", cli.DimStr("    "+snippet))
		}
	}
	return nil
}

func renderHelpHint(query string) {
	cli.Hint("browse", "core help")
	if trimmed := strings.TrimSpace(query); trimmed != "" {
		cli.Hint("search", fmt.Sprintf("core help search %q", trimmed))
	}
}

func renderTopicList(topics []*gohelp.Topic) error {
	if len(topics) == 0 {
		return cli.Err("no help topics available")
	}

	cli.Section("Available Help Topics")
	for _, topic := range topics {
		cli.Println("  %s - %s", topic.ID, topic.Title)
		if summary := topicSummary(topic); summary != "" {
			cli.Println("%s", cli.DimStr("    "+summary))
		}
	}
	return nil
}

func topicSummary(topic *gohelp.Topic) string {
	if topic == nil {
		return ""
	}

	content := strings.TrimSpace(topic.Content)
	if content == "" {
		return ""
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line
	}
	return ""
}

func renderTopic(t *gohelp.Topic) {
	cli.Blank()
	cli.Println("%s", cli.TitleStyle.Render(t.Title))
	cli.Println("%s", strings.Repeat("-", len(t.Title)))
	cli.Blank()
	cli.Println("%s", t.Content)
	cli.Blank()
}
