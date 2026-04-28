package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/config"
)

// AddConfigCommands registers the 'config' command group and all subcommands.
//
//	config.AddConfigCommands(c)
func AddConfigCommands(c *core.Core) {
	c.Command("config/get", core.Command{
		Description: "Get a configuration value",
		Action:      configGetAction,
	})
	c.Command("config/set", core.Command{
		Description: "Set a configuration value",
		Action:      configSetAction,
	})
	c.Command("config/list", core.Command{
		Description: "List all configuration values",
		Action:      configListAction,
	})
	c.Command("config/path", core.Command{
		Description: "Show the configuration file path",
		Action:      configPathAction,
	})
}

func loadConfig() (*config.Config, error) {
	configuration, err := config.New()
	if err != nil {
		return nil, cli.Wrap(err, "failed to load config")
	}
	return configuration, nil
}
