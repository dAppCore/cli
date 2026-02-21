package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/module"
	_ "forge.lthn.ai/core/cli/cmd/pkgcmd"
	_ "forge.lthn.ai/core/cli/cmd/plugin"
	_ "forge.lthn.ai/core/cli/cmd/session"
)

func main() {
	cli.Main()
}
