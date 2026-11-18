package config

import (
	"github.com/dappcore/go"
	"github.com/dappcore/cli/go/pkg/cli"
	"github.com/dappcore/config/go"
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
