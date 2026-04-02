package help

import (
	"bufio"
	"strings"

	"forge.lthn.ai/core/cli/pkg/cli"
	gohelp "forge.lthn.ai/core/go-help"
	"github.com/spf13/cobra"
)

func AddHelpCommands(root *cli.Command) {
	var searchQuery string

	helpCmd := &cli.Command{
		Use:   "help [topic]",
		Short: "Display help documentation",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cli.Command, args []string) error {
			catalog := gohelp.DefaultCatalog()

			if searchQuery != "" {
				return renderSearchResults(catalog.Search(searchQuery), searchQuery)
			}

			if len(args) == 0 {
				return renderTopicList(catalog.List())
			}

			topic, err := catalog.Get(args[0])
			if err != nil {
				return cli.Err("help topic %q not found", args[0])
			}

			renderTopic(topic)
			return nil
		},
	}

	helpCmd.Flags().StringVarP(&searchQuery, "search", "s", "", "Search help topics")
	root.AddCommand(helpCmd)
}

func renderSearchResults(results []*gohelp.SearchResult, query string) error {
	if len(results) == 0 {
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
