package config

import (
	"forge.lthn.ai/core/cli/pkg/cli"
)

func addGetCommand(parent *cli.Command) {
	cmd := cli.NewCommand("get", "Get a configuration value", "", func(cmd *cli.Command, args []string) error {
		key := args[0]

		configuration, err := loadConfig()
		if err != nil {
			return err
		}

		var value any
		if err := configuration.Get(key, &value); err != nil {
			return cli.Err("key not found: %s", key)
		}

		cli.Println("%v", value)
		return nil
	})

	cli.WithArgs(cmd, cli.ExactArgs(1))
	cli.WithExample(cmd, "core config get dev.editor")

	parent.AddCommand(cmd)
}
