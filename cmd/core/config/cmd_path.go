package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/config"
)

func configPathAction(_ core.Options) core.Result {
	configurationResult := loadConfig()
	if !configurationResult.OK {
		return configurationResult
	}
	configuration := configurationResult.Value.(*config.Config)

	cli.Println("%s", configuration.Path())
	return core.Ok(nil)
}
