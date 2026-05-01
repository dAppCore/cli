package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/config"
)

// AddConfigCommands registers the 'config' command group and all subcommands.
//
//	config.AddConfigCommands(c)
func AddConfigCommands(c *core.Core) core.Result {
	if r := c.Command("config/get", core.Command{
		Description: "Get a configuration value",
		Action:      configGetAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("config/set", core.Command{
		Description: "Set a configuration value",
		Action:      configSetAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("config/list", core.Command{
		Description: "List all configuration values",
		Action:      configListAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("config/path", core.Command{
		Description: "Show the configuration file path",
		Action:      configPathAction,
	}); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func loadConfig() core.Result {
	configuration, err := config.New()
	if err != nil {
		return cli.Wrap(err, "failed to load config")
	}
	return core.Ok(configuration)
}
