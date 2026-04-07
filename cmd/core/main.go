package main

import (
	"dappco.re/go/core/cli/cmd/core/config"
	"dappco.re/go/core/cli/cmd/core/doctor"
	"dappco.re/go/core/cli/cmd/core/help"
	"dappco.re/go/core/cli/cmd/core/pkgcmd"
	"dappco.re/go/core/cli/pkg/cli"

	// Ecosystem commands — self-register via init() + cli.RegisterCommands()
	_ "dappco.re/go/core/build/cmd/build"
	_ "dappco.re/go/core/build/cmd/ci"
	_ "dappco.re/go/core/build/cmd/sdk"
	_ "dappco.re/go/core/crypt/cmd/crypt"
	_ "dappco.re/go/core/devops/cmd/deploy"
	_ "dappco.re/go/core/devops/cmd/dev"
	_ "dappco.re/go/core/devops/cmd/docs"
	_ "dappco.re/go/core/devops/cmd/gitcmd"
	_ "dappco.re/go/core/devops/cmd/setup"
	_ "dappco.re/go/core/scm/cmd/collect"
	_ "dappco.re/go/core/scm/cmd/forge"
	_ "dappco.re/go/core/scm/cmd/gitea"
	_ "dappco.re/go/core/lint/cmd/qa"
)

func main() {
	cli.Main(
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
	)
}
