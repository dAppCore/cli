// Package doctor provides environment check commands.
package doctor

import (
	"dappco.re/go/core/cli/pkg/cli"
	"dappco.re/go/core/i18n"
	"github.com/spf13/cobra"
)

// Style aliases from shared
var (
	successStyle = cli.SuccessStyle
	errorStyle   = cli.ErrorStyle
	dimStyle     = cli.DimStyle
)

// Flag variable for doctor command
var doctorVerbose bool

var doctorCmd = &cobra.Command{
	Use: "doctor",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor(doctorVerbose)
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorVerbose, "verbose", false, i18n.T("cmd.doctor.verbose_flag"))
}

func runDoctor(verbose bool) error {
	cli.Println("%s", i18n.T("common.progress.checking", map[string]any{"Item": "development environment"}))
	cli.Blank()

	var passed, failed, optional int

	// Check required tools
	cli.Println("%s", i18n.T("cmd.doctor.required"))
	for _, toolCheck := range requiredChecks() {
		ok, version := runCheck(toolCheck)
		if ok {
			if verbose {
				cli.Println("%s", formatCheckResult(true, toolCheck.name, version))
			} else {
				cli.Println("%s", formatCheckResult(true, toolCheck.name, ""))
			}
			passed++
		} else {
			cli.Println("  %s %s - %s", errorStyle.Render(cli.Glyph(":cross:")), toolCheck.name, toolCheck.description)
			failed++
		}
	}

	// Check optional tools
	cli.Println("\n%s", i18n.T("cmd.doctor.optional"))
	for _, toolCheck := range optionalChecks() {
		ok, version := runCheck(toolCheck)
		if ok {
			if verbose {
				cli.Println("%s", formatCheckResult(true, toolCheck.name, version))
			} else {
				cli.Println("%s", formatCheckResult(true, toolCheck.name, ""))
			}
			passed++
		} else {
			cli.Println("  %s %s - %s", dimStyle.Render(cli.Glyph(":skip:")), toolCheck.name, dimStyle.Render(toolCheck.description))
			optional++
		}
	}

	// Check GitHub access
	cli.Println("\n%s", i18n.T("cmd.doctor.github"))
	if checkGitHubSSH() {
		cli.Println("%s", formatCheckResult(true, i18n.T("cmd.doctor.ssh_found"), ""))
	} else {
		cli.Println("  %s %s", errorStyle.Render(cli.Glyph(":cross:")), i18n.T("cmd.doctor.ssh_missing"))
		failed++
	}

	if checkGitHubCLI() {
		cli.Println("%s", formatCheckResult(true, i18n.T("cmd.doctor.cli_auth"), ""))
	} else {
		cli.Println("  %s %s", errorStyle.Render(cli.Glyph(":cross:")), i18n.T("cmd.doctor.cli_auth_missing"))
		failed++
	}

	// Check workspace
	cli.Println("\n%s", i18n.T("cmd.doctor.workspace"))
	checkWorkspace()

	// Summary
	cli.Blank()
	if failed > 0 {
		cli.Error(i18n.T("cmd.doctor.issues", map[string]any{"Count": failed}))
		cli.Println("\n%s", i18n.T("cmd.doctor.install_missing"))
		printInstallInstructions()
		return cli.Err("%s", i18n.T("cmd.doctor.issues_error", map[string]any{"Count": failed}))
	}

	cli.Success(i18n.T("cmd.doctor.ready"))
	return nil
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
