package config

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/config"
)

// configSetAction handles 'config set --key=<key> --value=<value>'.
// Also accepts positional form via _arg for backwards compatibility when
// only one arg is passed (interpreted as key, value read from --value).
func configSetAction(opts core.Options) core.Result {
	key := opts.String("key")
	value := opts.String("value")

	// Fallback: first positional arg as key if --key not provided.
	if key == "" {
		key = opts.String("_arg")
	}

	if key == "" {
		return cli.Err("requires --key and --value arguments (e.g. config set --key=dev.editor --value=vim)")
	}
	if value == "" {
		return cli.Err("requires --value argument (e.g. config set --key=%s --value=<value>)", key)
	}

	configurationResult := loadConfig()
	if !configurationResult.OK {
		return configurationResult
	}
	configuration := configurationResult.Value.(*config.Config)

	if err := configuration.Set(key, value); err != nil {
		return cli.Wrap(err, "failed to set config value")
	}

	cli.Success(key + " = " + value)
	return core.Ok(nil)
}
