package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

func configGetAction(opts core.Options) core.Result {
	key := opts.String("_arg")
	if key == "" {
		return core.Fail(cli.Err("requires a configuration key argument"))
	}

	configuration, err := loadConfig()
	if err != nil {
		return core.Fail(err)
	}

	var value any
	if err := configuration.Get(key, &value); err != nil {
		return core.Fail(cli.Err("key not found: %s", key))
	}

	cli.Println("%v", value)
	return core.Ok(nil)
}
