package main

import (
	"dappco.re/go/cli/cmd/core/config"
	"dappco.re/go/cli/cmd/core/doctor"
	"dappco.re/go/cli/cmd/core/help"
	"dappco.re/go/cli/cmd/core/pkgcmd"
	"dappco.re/go/cli/pkg/cli"

	// Ecosystem commands — self-register via init() + cli.RegisterCommands()
	_ "dappco.re/go/build/cmd/build"
	_ "dappco.re/go/build/cmd/ci"
	_ "dappco.re/go/build/cmd/sdk"
	_ "dappco.re/go/crypt/cmd/crypt"
	_ "dappco.re/go/devops/cmd/deploy"
	_ "dappco.re/go/devops/cmd/dev"
	_ "dappco.re/go/devops/cmd/docs"
	_ "dappco.re/go/devops/cmd/gitcmd"
	_ "dappco.re/go/devops/cmd/setup"
	_ "dappco.re/go/scm/cmd/collect"
	_ "dappco.re/go/scm/cmd/forge"
	_ "dappco.re/go/scm/cmd/gitea"
	_ "dappco.re/go/lint/cmd/qa"
)

func main() {
	cli.Main(
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
	)
}
