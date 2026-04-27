package config

import (
	"dappco.re/go/core"
	"dappco.re/go/cli/pkg/cli"
)

func configPathAction(_ core.Options) core.Result {
	configuration, err := loadConfig()
	if err != nil {
		return core.Result{Value: err, OK: false}
	}

	cli.Println("%s", configuration.Path())
	return core.Result{OK: true}
}
