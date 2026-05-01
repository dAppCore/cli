package help

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/help"
)

// AddHelpCommands registers the help command and subcommands.
//
//	help.AddHelpCommands(c)
func AddHelpCommands(c *core.Core) core.Result {
	if r := c.Command("help", core.Command{
		Description: "Display help documentation",
		Action:      helpAction,
	}); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func helpAction(opts core.Options) core.Result {
	catalog := help.DefaultCatalog()
	search := opts.String("search")

	if search != "" {
		results := catalog.Search(search)
		if len(results) == 0 {
			cli.Println("No topics found.")
			return core.Ok(nil)
		}
		cli.Println("Search Results:")
		for _, result := range results {
			cli.Println("  %s - %s", result.Topic.ID, result.Topic.Title)
		}
		return core.Ok(nil)
	}

	// Check for topic argument
	topicID := opts.String("_arg")
	if topicID == "" {
		topics := catalog.List()
		cli.Println("Available Help Topics:")
		for _, topic := range topics {
			cli.Println("  %s - %s", topic.ID, topic.Title)
		}
		return core.Ok(nil)
	}

	topic, err := catalog.Get(topicID)
	if err != nil {
		return cli.Err("Error: %v", err)
	}

	renderTopic(topic)
	return core.Ok(nil)
}

func renderTopic(topic *help.Topic) {
	cli.Println("\n%s", cli.TitleStyle.Render(topic.Title))
	cli.Println("----------------------------------------")
	cli.Println("%s", topic.Content)
	cli.Blank()
}
