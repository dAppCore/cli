package main

import (
	"forge.lthn.ai/core/cli/cmd/config"
	"forge.lthn.ai/core/cli/cmd/doctor"
	"forge.lthn.ai/core/cli/cmd/gocmd"
	"forge.lthn.ai/core/cli/cmd/help"
	"forge.lthn.ai/core/cli/cmd/module"
	"forge.lthn.ai/core/cli/cmd/pkgcmd"
	"forge.lthn.ai/core/cli/cmd/plugin"
	"forge.lthn.ai/core/cli/cmd/service"
	"forge.lthn.ai/core/cli/cmd/session"
	"forge.lthn.ai/core/cli/pkg/cli"

	// Ecosystem command packages — self-register via init() + cli.RegisterCommands()
	_ "forge.lthn.ai/core/agent/cmd/agent"
	_ "forge.lthn.ai/core/agent/cmd/dispatch"
	_ "forge.lthn.ai/core/agent/cmd/taskgit"
	_ "forge.lthn.ai/core/go-ansible/cmd/ansible"
	_ "forge.lthn.ai/core/go-api/cmd/api"
	_ "forge.lthn.ai/core/go-build/cmd/build"
	_ "forge.lthn.ai/core/go-build/cmd/ci"
	_ "forge.lthn.ai/core/go-build/cmd/sdk"
	_ "forge.lthn.ai/core/go-container/cmd/vm"
	_ "forge.lthn.ai/core/go-crypt/cmd/crypt"
	_ "forge.lthn.ai/core/go-devops/cmd/deploy"
	_ "forge.lthn.ai/core/go-devops/cmd/dev"
	_ "forge.lthn.ai/core/go-devops/cmd/docs"
	_ "forge.lthn.ai/core/go-devops/cmd/gitcmd"
	_ "forge.lthn.ai/core/go-devops/cmd/setup"
	_ "forge.lthn.ai/core/go-infra/cmd/monitor"
	_ "forge.lthn.ai/core/go-infra/cmd/prod"
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
		cli.WithCommands("module", module.AddModuleCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),
		cli.WithCommands("plugin", plugin.AddPluginCommands),
		cli.WithCommands("session", session.AddSessionCommands),
		cli.WithCommands("go", gocmd.AddGoCommands),
		cli.WithCommands("service", service.AddServiceCommands),
	)
}
