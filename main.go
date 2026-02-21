package main

import (
	"forge.lthn.ai/core/go/pkg/cli"

	// Commands — local to CLI
	_ "forge.lthn.ai/core/cli/cmd/config"
	_ "forge.lthn.ai/core/cli/cmd/doctor"
	_ "forge.lthn.ai/core/cli/cmd/help"
	_ "forge.lthn.ai/core/cli/cmd/module"
	_ "forge.lthn.ai/core/cli/cmd/pkgcmd"
	_ "forge.lthn.ai/core/cli/cmd/plugin"
	_ "forge.lthn.ai/core/cli/cmd/session"

	// Commands — from ecosystem repos
	_ "forge.lthn.ai/core/go/cmd/gocmd"
	_ "forge.lthn.ai/core/go-agentic/cmd/workspace"
	_ "forge.lthn.ai/core/go-ai/cmd/lab"
	_ "forge.lthn.ai/core/go-devops/cmd/dev"
	_ "forge.lthn.ai/core/go-devops/cmd/docs"
	_ "forge.lthn.ai/core/go-devops/cmd/gitcmd"
	_ "forge.lthn.ai/core/go-devops/cmd/monitor"
	_ "forge.lthn.ai/core/go-devops/cmd/qa"
	_ "forge.lthn.ai/core/go-devops/cmd/setup"
)

func main() {
	cli.Main()
}
