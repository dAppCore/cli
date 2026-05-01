package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/config"
)

func configGetAction(opts core.Options) core.Result {
	key := opts.String("_arg")
	if key == "" {
		return cli.Err("requires a configuration key argument")
	}

	configurationResult := loadConfig()
	if !configurationResult.OK {
		return configurationResult
	}
	configuration := configurationResult.Value.(*config.Config)

	var value any
	if err := configuration.Get(key, &value); err != nil {
		return cli.Err("key not found: %s", key)
	}

	cli.Println("%v", value)
	return core.Ok(nil)
}
