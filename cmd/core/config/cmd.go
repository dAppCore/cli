package config

import (
	"dappco.re/go/core/cli/pkg/cli"
	"dappco.re/go/core/config"
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
