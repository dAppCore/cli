package main

import (
	"forge.lthn.ai/core/cli/cmd/config"
	"forge.lthn.ai/core/cli/cmd/doctor"
	"forge.lthn.ai/core/cli/cmd/gocmd"
	"forge.lthn.ai/core/cli/cmd/help"
	"forge.lthn.ai/core/cli/cmd/module"
	"forge.lthn.ai/core/cli/cmd/pkgcmd"
	"forge.lthn.ai/core/cli/cmd/plugin"
	"forge.lthn.ai/core/cli/cmd/session"
	"forge.lthn.ai/core/cli/pkg/cli"
)

func main() {
	cli.Main(
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("module", module.AddModuleCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
		cli.WithCommands("plugin", plugin.AddPluginCommands),
		cli.WithCommands("session", session.AddSessionCommands),
		cli.WithCommands("go", gocmd.AddGoCommands),
	)
}
