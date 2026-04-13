package doctor

import (
	"os/exec"

	"dappco.re/go/core"
	"dappco.re/go/core/i18n"
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
			name:        i18n.T("cmd.doctor.check.git.name"),
			description: i18n.T("cmd.doctor.check.git.description"),
			command:     "git",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        i18n.T("cmd.doctor.check.go.name"),
			description: i18n.T("cmd.doctor.check.go.description"),
			command:     "go",
			args:        []string{"version"},
			versionFlag: "version",
		},
		{
			name:        i18n.T("cmd.doctor.check.gh.name"),
			description: i18n.T("cmd.doctor.check.gh.description"),
			command:     "gh",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        i18n.T("cmd.doctor.check.php.name"),
			description: i18n.T("cmd.doctor.check.php.description"),
			command:     "php",
			args:        []string{"-v"},
			versionFlag: "-v",
		},
		{
			name:        i18n.T("cmd.doctor.check.composer.name"),
			description: i18n.T("cmd.doctor.check.composer.description"),
			command:     "composer",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        i18n.T("cmd.doctor.check.node.name"),
			description: i18n.T("cmd.doctor.check.node.description"),
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
			name:        i18n.T("cmd.doctor.check.pnpm.name"),
			description: i18n.T("cmd.doctor.check.pnpm.description"),
			command:     "pnpm",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        i18n.T("cmd.doctor.check.claude.name"),
			description: i18n.T("cmd.doctor.check.claude.description"),
			command:     "claude",
			args:        []string{"--version"},
			versionFlag: "--version",
		},
		{
			name:        i18n.T("cmd.doctor.check.docker.name"),
			description: i18n.T("cmd.doctor.check.docker.description"),
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
	proc := exec.Command(toolCheck.command, toolCheck.args...) // TODO: migrate to c.Process()
	output, err := proc.CombinedOutput()
	if err != nil {
		return false, ""
	}

	// Extract first line as version info.
	lines := core.Split(core.Trim(string(output)), "\n")
	if len(lines) > 0 {
		return true, core.Trim(lines[0])
	}
	return true, ""
}
