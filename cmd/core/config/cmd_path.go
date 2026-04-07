package config

import (
	"dappco.re/go/core/cli/pkg/cli"
)

func addPathCommand(parent *cli.Command) {
	cmd := cli.NewCommand("path", "Show the configuration file path", "", func(cmd *cli.Command, args []string) error {
		configuration, err := loadConfig()
		if err != nil {
			return err
		}

		cli.Println("%s", configuration.Path())
		return nil
	})

	cli.WithArgs(cmd, cli.NoArgs())

	parent.AddCommand(cmd)
}
