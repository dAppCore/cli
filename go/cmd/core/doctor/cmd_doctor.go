// Package doctor provides environment check commands.
package doctor

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// Style aliases from shared
var (
	successStyle = cli.SuccessStyle
	errorStyle   = cli.ErrorStyle
	dimStyle     = cli.DimStyle
)

func doctorAction(opts core.Options) core.Result {
	verbose := opts.Bool("verbose")
	if r := runDoctor(verbose); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func runDoctor(verbose bool) core.Result {
	cli.Println("%s", cli.T("common.progress.checking", map[string]any{"Item": "development environment"}))
	cli.Blank()

	failed := checkRequiredTools(verbose)
	checkOptionalTools(verbose)
	failed += checkGitHub()

	// Check workspace
	cli.Println("\n%s", cli.T("cmd.doctor.workspace"))
	checkWorkspace()

	// Summary
	cli.Blank()
	if failed > 0 {
		cli.Error(cli.T("cmd.doctor.issues", map[string]any{"Count": failed}))
		cli.Println("\n%s", cli.T("cmd.doctor.install_missing"))
		printInstallInstructions()
		return cli.Err("%s", cli.T("cmd.doctor.issues_error", map[string]any{"Count": failed}))
	}

	cli.Success(cli.T("cmd.doctor.ready"))
	return core.Ok(nil)
}

func formatCheckResult(ok bool, name, detail string) string {
	checkBuilder := cli.Check(name)
	if ok {
		checkBuilder.Pass()
	} else {
		checkBuilder.Fail()
	}
	if detail != "" {
		checkBuilder.Message(detail)
	} else {
		checkBuilder.Message("")
	}
	return checkBuilder.String()
}
