package config

import (
	"dappco.re/go/core"
	"dappco.re/go/core/cli/pkg/cli"
)

func configGetAction(opts core.Options) core.Result {
	key := opts.String("_arg")
	if key == "" {
		return core.Result{Value: cli.Err("requires a configuration key argument"), OK: false}
	}

	configuration, err := loadConfig()
	if err != nil {
		return core.Result{Value: err, OK: false}
	}

	var value any
	if err := configuration.Get(key, &value); err != nil {
		return core.Result{Value: cli.Err("key not found: %s", key), OK: false}
	}

	cli.Println("%v", value)
	return core.Result{OK: true}
}
