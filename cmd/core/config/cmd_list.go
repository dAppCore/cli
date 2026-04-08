package config

import (
	"maps"

	"dappco.re/go/core"
	"dappco.re/go/core/cli/pkg/cli"
	"gopkg.in/yaml.v3"
)

func configListAction(_ core.Options) core.Result {
	configuration, err := loadConfig()
	if err != nil {
		return core.Result{Value: err, OK: false}
	}

	all := maps.Collect(configuration.All())
	if len(all) == 0 {
		cli.Dim("No configuration values set")
		return core.Result{OK: true}
	}

	output, err := yaml.Marshal(all)
	if err != nil {
		return core.Result{Value: cli.Wrap(err, "failed to format config"), OK: false}
	}

	cli.Print("%s", string(output))
	return core.Result{OK: true}
}
