package config

import (
	"maps"

	"forge.lthn.ai/core/cli/pkg/cli"
	"gopkg.in/yaml.v3"
)

func addListCommand(parent *cli.Command) {
	cmd := cli.NewCommand("list", "List all configuration values", "", func(cmd *cli.Command, args []string) error {
		configuration, err := loadConfig()
		if err != nil {
			return err
		}

		all := maps.Collect(configuration.All())
		if len(all) == 0 {
			cli.Dim("No configuration values set")
			return nil
		}

		output, err := yaml.Marshal(all)
		if err != nil {
			return cli.Wrap(err, "failed to format config")
		}

		cli.Print("%s", string(output))
		return nil
	})

	cli.WithArgs(cmd, cli.NoArgs())

	parent.AddCommand(cmd)
}
