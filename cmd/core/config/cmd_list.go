package config

import (
	"maps"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"gopkg.in/yaml.v3"
)

func configListAction(_ core.Options) core.Result {
	configuration, err := loadConfig()
	if err != nil {
		return core.Fail(err)
	}

	all := maps.Collect(configuration.All())
	if len(all) == 0 {
		cli.Dim("No configuration values set")
		return core.Ok(nil)
	}

	output, err := yaml.Marshal(all)
	if err != nil {
		return core.Fail(cli.Wrap(err, "failed to format config"))
	}

	cli.Print("%s", string(output))
	return core.Ok(nil)
}
