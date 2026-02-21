package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	// Commands via self-registration (local to CLI)
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/dev"
	_ "forge.lthn.ai/core/cli/cmd/docs"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/gitcmd"
	_ "forge.lthn.ai/core/cli/cmd/go"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/lab"
	_ "forge.lthn.ai/core/cli/cmd/module"
	_ "forge.lthn.ai/core/cli/cmd/monitor"
	_ "forge.lthn.ai/core/cli/cmd/pkgcmd"
	_ "forge.lthn.ai/core/cli/cmd/plugin"
	_ "forge.lthn.ai/core/cli/cmd/qa"
	_ "forge.lthn.ai/core/cli/cmd/session"
	_ "forge.lthn.ai/core/cli/cmd/setup"
	_ "forge.lthn.ai/core/cli/cmd/workspace"
)

func main() {
	cli.Main()
}
