package config

import "forge.lthn.ai/core/go/pkg/cli"

// AddConfigCommands registers the 'config' command group and all subcommands.
func AddConfigCommands(root *cli.Command) {
	configCmd := cli.NewGroup("config", "Manage configuration", "")
	root.AddCommand(configCmd)

	addGetCommand(configCmd)
	addSetCommand(configCmd)
	addListCommand(configCmd)
	addPathCommand(configCmd)
}
