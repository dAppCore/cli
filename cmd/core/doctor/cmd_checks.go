package doctor

import (
	"context"

	"dappco.re/go/core"
	"dappco.re/go/core/process"
	"dappco.re/go/i18n"
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
	ctx := context.Background()
	processCore := core.New(core.WithService(process.Register))
	if startup := processCore.ServiceStartup(ctx, nil); !startup.OK {
		return false, ""
	}
	defer processCore.ServiceShutdown(context.Background())

	result := processCore.Process().Run(ctx, toolCheck.command, toolCheck.args...)
	if !result.OK {
		return false, ""
	}

	output, ok := result.Value.(string)
	if !ok {
		return false, ""
	}

	// Extract first line as version info.
	lines := core.Split(core.Trim(output), "\n")
	if len(lines) > 0 {
		return true, core.Trim(lines[0])
	}
	return true, ""
}
