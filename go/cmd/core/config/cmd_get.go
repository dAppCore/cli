package config

import (
	"github.com/dappcore/go"
	"github.com/dappcore/cli/go/pkg/cli"
	"github.com/dappcore/config/go"
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
	if r := configuration.Get(key, &value); !r.OK {
		return cli.Err("key not found: %s", key)
	}

	cli.Println("%v", value)
	return core.Ok(nil)
}
