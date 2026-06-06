package doctor

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// check represents a tool check configuration
type check struct {
	name        string
	description string
	command     string
	args        []string
	versionFlag string
}

// requiredChecks returns tools that must be installed
func requiredChecks() []check {
	return []check{
		{
			name:        cli.T("cmd.doctor.check.git.name"),
			description: cli.T("cmd.doctor.check.git.description"),
			command:     "git",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        cli.T("cmd.doctor.check.go.name"),
			description: cli.T("cmd.doctor.check.go.description"),
			command:     "go",
			args:        []string{"version"},
			versionFlag: "version",
		},
		{
			name:        cli.T("cmd.doctor.check.gh.name"),
			description: cli.T("cmd.doctor.check.gh.description"),
			command:     "gh",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        cli.T("cmd.doctor.check.php.name"),
			description: cli.T("cmd.doctor.check.php.description"),
			command:     "php",
			args:        []string{"-v"},
			versionFlag: "-v",
		},
		{
			name:        cli.T("cmd.doctor.check.composer.name"),
			description: cli.T("cmd.doctor.check.composer.description"),
			command:     "composer",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        cli.T("cmd.doctor.check.node.name"),
			description: cli.T("cmd.doctor.check.node.description"),
			command:     "node",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
	}
}

// optionalChecks returns tools that are nice to have
func optionalChecks() []check {
	return []check{
		{
			name:        cli.T("cmd.doctor.check.pnpm.name"),
			description: cli.T("cmd.doctor.check.pnpm.description"),
			command:     "pnpm",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        cli.T("cmd.doctor.check.claude.name"),
			description: cli.T("cmd.doctor.check.claude.description"),
			command:     "claude",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        cli.T("cmd.doctor.check.docker.name"),
			description: cli.T("cmd.doctor.check.docker.description"),
			command:     "docker",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
	}
}

// runCheck executes a tool check and returns success status and version info.
//
//	ok, version := runCheck(check{command: "git", args: []string{"--version"}})
func runCheck(toolCheck check) (bool, string) {
	if !(core.App{}).Find(toolCheck.command, toolCheck.name).OK {
		return false, ""
	}
	return true, toolCheck.command
}

// checkRequiredTools prints the required-tool section and returns the number of
// missing required tools.
func checkRequiredTools(verbose bool) (failed int) {
	cli.Println("%s", cli.T("cmd.doctor.required"))
	for _, toolCheck := range requiredChecks() {
		ok, version := runCheck(toolCheck)
		if ok {
			detail := ""
			if verbose {
				detail = version
			}
			cli.Println("%s", formatCheckResult(true, toolCheck.name, detail))
		} else {
			cli.Println("  %s %s - %s", errorStyle.Render(cli.Glyph(":cross:")), toolCheck.name, toolCheck.description)
			failed++
		}
	}
	return failed
}

// checkOptionalTools prints the optional-tool section. Missing optional tools do
// not count as failures.
func checkOptionalTools(verbose bool) {
	cli.Println("\n%s", cli.T("cmd.doctor.optional"))
	for _, toolCheck := range optionalChecks() {
		ok, version := runCheck(toolCheck)
		if ok {
			detail := ""
			if verbose {
				detail = version
			}
			cli.Println("%s", formatCheckResult(true, toolCheck.name, detail))
		} else {
			cli.Println("  %s %s - %s", dimStyle.Render(cli.Glyph(":skip:")), toolCheck.name, dimStyle.Render(toolCheck.description))
		}
	}
}

// checkGitHub prints the GitHub-access section and returns the number of failed
// GitHub checks (missing SSH key, unauthenticated CLI).
func checkGitHub() (failed int) {
	cli.Println("\n%s", cli.T("cmd.doctor.github"))
	if checkGitHubSSH() {
		cli.Println("%s", formatCheckResult(true, cli.T("cmd.doctor.ssh_found"), ""))
	} else {
		cli.Println("  %s %s", errorStyle.Render(cli.Glyph(":cross:")), cli.T("cmd.doctor.ssh_missing"))
		failed++
	}
	if checkGitHubCLI() {
		cli.Println("%s", formatCheckResult(true, cli.T("cmd.doctor.cli_auth"), ""))
	} else {
		cli.Println("  %s %s", errorStyle.Render(cli.Glyph(":cross:")), cli.T("cmd.doctor.cli_auth_missing"))
		failed++
	}
	return failed
}

// doctorChecksAction runs the tool and GitHub diagnostics without the full-scan
// banner or workspace summary. It fails when any required tool or GitHub check
// is missing.
func doctorChecksAction(opts core.Options) core.Result {
	verbose := opts.Bool("verbose")

	failed := checkRequiredTools(verbose)
	checkOptionalTools(verbose)
	failed += checkGitHub()

	cli.Blank()
	if failed > 0 {
		return cli.Err("%s", cli.T("cmd.doctor.issues_error", map[string]any{"Count": failed}))
	}
	cli.Success(cli.T("cmd.doctor.ready"))
	return core.Ok(nil)
}
