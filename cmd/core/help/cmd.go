package help

import (
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
	}
	return nil
}

func renderTopic(t *gohelp.Topic) {
	cli.Blank()
	cli.Println("%s", cli.TitleStyle.Render(t.Title))
	cli.Println("%s", strings.Repeat("-", len(t.Title)))
	cli.Blank()
	cli.Println("%s", t.Content)
	cli.Blank()
}
