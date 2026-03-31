package help

import (
	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-help"
)

// AddHelpCommands registers the help command and subcommands.
//
//	help.AddHelpCommands(rootCmd)
func AddHelpCommands(root *cli.Command) {
	var searchFlag string

	helpCmd := &cli.Command{
		Use:   "help [topic]",
		Short: "Display help documentation",
		Run: func(cmd *cli.Command, args []string) {
			catalog := help.DefaultCatalog()

			if searchFlag != "" {
				results := catalog.Search(searchFlag)
				if len(results) == 0 {
					cli.Println("No topics found.")
					return
				}
				cli.Println("Search Results:")
				for _, result := range results {
					cli.Println("  %s - %s", result.Topic.ID, result.Topic.Title)
				}
				return
			}

			if len(args) == 0 {
				topics := catalog.List()
				cli.Println("Available Help Topics:")
				for _, topic := range topics {
					cli.Println("  %s - %s", topic.ID, topic.Title)
				}
				return
			}

			topic, err := catalog.Get(args[0])
			if err != nil {
				cli.Errorf("Error: %v", err)
				return
			}

			renderTopic(topic)
		},
	}

	helpCmd.Flags().StringVarP(&searchFlag, "search", "s", "", "Search help topics")
	root.AddCommand(helpCmd)
}

func renderTopic(topic *help.Topic) {
	cli.Println("\n%s", cli.TitleStyle.Render(topic.Title))
	cli.Println("----------------------------------------")
	cli.Println("%s", topic.Content)
	cli.Blank()
}
