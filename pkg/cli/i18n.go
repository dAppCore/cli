package cli

import (
	"forge.lthn.ai/core/go-i18n"
)

// T translates a key using the CLI's i18n service.
// Falls back to the global i18n.T if CLI not initialised.
func T(key string, args ...map[string]any) string {
	if len(args) > 0 {
		return i18n.T(key, args[0])
	}
	return i18n.T(key)
}
