package main

import (
	"forge.lthn.ai/core/cli/cmd/core/config"
	"forge.lthn.ai/core/cli/cmd/core/doctor"
	"forge.lthn.ai/core/cli/cmd/core/help"
	"forge.lthn.ai/core/cli/cmd/core/pkgcmd"
	"forge.lthn.ai/core/cli/pkg/cli"
)

func main() {
	cli.Main(
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
	)
}
