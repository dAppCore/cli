package config

import (
	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/config"
)

// AddConfigCommands registers the 'config' command group and all subcommands.
//
//	config.AddConfigCommands(rootCmd)
func AddConfigCommands(root *cli.Command) {
	configCmd := cli.NewGroup("config", "Manage configuration", "")
	root.AddCommand(configCmd)

	addGetCommand(configCmd)
	addSetCommand(configCmd)
	addListCommand(configCmd)
	addPathCommand(configCmd)
}

func loadConfig() (*config.Config, error) {
	configuration, err := config.New()
	if err != nil {
		return nil, cli.Wrap(err, "failed to load config")
	}
	return configuration, nil
}
