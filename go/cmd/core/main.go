// Command core is the Lethean/dAppCore developer CLI — the token-efficient
// build, lint, qa, and packaging tooling that the ecosystem grew out of.
//
// It composes on the modern CoreGo flow: cli.Main builds a *core.Core (the same
// substrate the GUI uses via core.New(...), just without the gui/Wails service)
// and each command surface registers onto it through an AddXCommands(*core.Core)
// registrar handed to cli.WithCommands. Ecosystem commands are wired EXPLICITLY
// — the lib package is imported and its registrar passed in, not blank-imported
// for init() self-registration (the cmd binaries are package main, not
// importable, so that older pattern is dead).
package main

import (
	"dappco.re/go/cli/cmd/core/config"
	"dappco.re/go/cli/cmd/core/doctor"
	"dappco.re/go/cli/cmd/core/help"
	"dappco.re/go/cli/cmd/core/pkgcmd"
	"dappco.re/go/cli/pkg/cli"

	build "dappco.re/go/build/cmd/build"
	gocmd "dappco.re/go/lint/cmd/gocmd"
	"dappco.re/go/lint/cmd/qa"
)

func main() {
	cli.Main(
		// Local CLI built-ins.
		cli.WithCommands("config", config.AddConfigCommands),
		cli.WithCommands("doctor", doctor.AddDoctorCommands),
		cli.WithCommands("help", help.AddHelpCommands),
		cli.WithCommands("pkg", pkgcmd.AddPkgCommands),

		// Ecosystem surface — core {build, go, qa, …}. lint/php follow once this
		// shape is signed off.
		cli.WithCommands("build", build.AddBuildCommands),
		cli.WithCommands("go", gocmd.AddGoCommands),
		cli.WithCommands("qa", qa.AddQACommands),

		// Demo capability so the action→CLI projection has something to mount.
		AddDemoActions,

		// Opt-in: project the action registry (the capability map) onto the CLI.
		// Runs last so the explicit commands above win any path collision. This
		// is the third surface for the same map the API and MCP layers project.
		cli.MountActions,
	)
}
