package config

import (
	"dappco.re/go/core"
	"dappco.re/go/core/cli/pkg/cli"
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
		return core.Result{Value: cli.Err("requires --key and --value arguments (e.g. config set --key=dev.editor --value=vim)"), OK: false}
	}
	if value == "" {
		return core.Result{Value: cli.Err("requires --value argument (e.g. config set --key=%s --value=<value>)", key), OK: false}
	}

	configuration, err := loadConfig()
	if err != nil {
		return core.Result{Value: err, OK: false}
	}

	if err := configuration.Set(key, value); err != nil {
		return core.Result{Value: cli.Wrap(err, "failed to set config value"), OK: false}
	}

	cli.Success(key + " = " + value)
	return core.Result{OK: true}
}
