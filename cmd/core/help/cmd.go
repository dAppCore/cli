package help

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/help"
)

// AddHelpCommands registers the help command and subcommands.
//
//	help.AddHelpCommands(c)
func AddHelpCommands(c *core.Core) {
	c.Command("help", core.Command{
		Description: "Display help documentation",
		Action:      helpAction,
	})
}

func helpAction(opts core.Options) core.Result {
	catalog := help.DefaultCatalog()
	search := opts.String("search")

	if search != "" {
		results := catalog.Search(search)
		if len(results) == 0 {
			cli.Println("No topics found.")
			return core.Result{OK: true}
		}
		cli.Println("Search Results:")
		for _, result := range results {
			cli.Println("  %s - %s", result.Topic.ID, result.Topic.Title)
		}
		return core.Result{OK: true}
	}

	// Check for topic argument
	topicID := opts.String("_arg")
	if topicID == "" {
		topics := catalog.List()
		cli.Println("Available Help Topics:")
		for _, topic := range topics {
			cli.Println("  %s - %s", topic.ID, topic.Title)
		}
		return core.Result{OK: true}
	}

	topic, err := catalog.Get(topicID)
	if err != nil {
		return core.Result{Value: cli.Err("Error: %v", err), OK: false}
	}

	renderTopic(topic)
	return core.Result{OK: true}
}

func renderTopic(topic *help.Topic) {
	cli.Println("\n%s", cli.TitleStyle.Render(topic.Title))
	cli.Println("----------------------------------------")
	cli.Println("%s", topic.Content)
	cli.Blank()
}
