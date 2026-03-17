package main

import (
	"forge.lthn.ai/core/cli/cmd/core/config"
	"forge.lthn.ai/core/cli/cmd/core/doctor"
	"forge.lthn.ai/core/cli/cmd/core/help"
	"forge.lthn.ai/core/cli/cmd/core/pkgcmd"
	"forge.lthn.ai/core/cli/pkg/cli"

	// Ecosystem commands — self-register via init() + cli.RegisterCommands()
	// TODO: go-build has SDK dep conflict (kin-openapi vs oasdiff), uncomment when fixed
	// _ "forge.lthn.ai/core/go-build/cmd/build"
	// _ "forge.lthn.ai/core/go-build/cmd/ci"
	// _ "forge.lthn.ai/core/go-build/cmd/sdk"
	_ "forge.lthn.ai/core/go-crypt/cmd/crypt"
	_ "forge.lthn.ai/core/go-devops/cmd/deploy"
	_ "forge.lthn.ai/core/go-devops/cmd/dev"
	_ "forge.lthn.ai/core/go-devops/cmd/docs"
	_ "forge.lthn.ai/core/go-devops/cmd/gitcmd"
	_ "forge.lthn.ai/core/go-devops/cmd/setup"
	_ "forge.lthn.ai/core/go-scm/cmd/collect"
	_ "forge.lthn.ai/core/go-scm/cmd/forge"
	_ "forge.lthn.ai/core/go-scm/cmd/gitea"
	_ "forge.lthn.ai/core/lint/cmd/qa"
)

func main() {
	cli.Main(
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
	)
}
