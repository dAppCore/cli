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
	"github.com/dappcore/cli/go/cmd/core/config"
	"github.com/dappcore/cli/go/cmd/core/doctor"
	"github.com/dappcore/cli/go/cmd/core/help"
	"github.com/dappcore/cli/go/cmd/core/pkgcmd"
	"github.com/dappcore/cli/go/pkg/cli"

	build "github.com/dappcore/go-build/go/cmd/build"
	buildlocales "github.com/dappcore/go-build/go/locales"
	gocmd "github.com/dappcore/lint/go/cmd/gocmd"
	"github.com/dappcore/lint/go/cmd/qa"
	lintlocales "github.com/dappcore/lint/go/locales"
)

func main() {
	// Each module that registers commands ships the strings for them, and the
	// binary composing those commands is the only place that knows which to
	// load. Dir "." because each embeds its locales at its own package root.
	locales := []cli.LocaleSource{
		cli.WithLocales(buildlocales.FS, "."),
		cli.WithLocales(lintlocales.FS, "."),
	}

	cli.MainWithLocales(locales,
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

		// Opt-in: project the action registry (the capability map) onto the CLI,
		// with go-i18n-generated help. Runs last so the explicit commands above
		// win any path collision. This is the third surface for the same map the
		// API and MCP layers project.
		mountActionsWithGrammar,
	)
}
