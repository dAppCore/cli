package config

import (
	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-config"
)

// AddConfigCommands registers the 'config' command group and all subcommands.
func AddConfigCommands(root *cli.Command) {
	configCmd := cli.NewGroup("config", "Manage configuration", "")
	root.AddCommand(configCmd)

	addGetCommand(configCmd)
	addSetCommand(configCmd)
	addListCommand(configCmd)
	addPathCommand(configCmd)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, cli.Wrap(err, "failed to load config")
	}
	return cfg, nil
}
