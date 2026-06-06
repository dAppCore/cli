// Package doctor provides environment validation commands.
//
// Checks for:
//   - Required tools: git, gh, php, composer, node
//   - Optional tools: pnpm, claude, docker
//   - GitHub access: SSH keys and CLI authentication
//   - Workspace: repos.yaml presence and clone status
//
// Run before 'core setup' to ensure your environment is ready.
// Provides platform-specific installation instructions for missing tools.
//
// `core doctor` runs the full scan; the focused subcommands are:
//   - core doctor environment — detected system information
//   - core doctor commands    — registered core commands
//   - core doctor checks       — tool and GitHub diagnostics
//   - core doctor install      — install instructions for missing tools
package doctor

import (
	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
)

// AddDoctorCommands registers the 'doctor' command and all subcommands.
//
//	doctor.AddDoctorCommands(c)
func AddDoctorCommands(c *core.Core) core.Result {
	commands := []struct {
		path        string
		description string
		action      core.CommandAction
	}{
		{"doctor", "Check development environment health", doctorAction},
		{"doctor/environment", "Show detected system information", doctorEnvironmentAction},
		{"doctor/commands", "List registered core commands", doctorCommandsAction},
		{"doctor/checks", "Run tool and GitHub diagnostics", doctorChecksAction},
		{"doctor/install", "Show install instructions for missing tools", doctorInstallAction},
	}
	for _, cmd := range commands {
		if r := c.Command(cmd.path, core.Command{
			Description: cmd.description,
			Action:      cmd.action,
		}); !r.OK {
			return r
		}
	}
	return core.Ok(nil)
}

// doctorCommandsAction lists the registered core commands and their descriptions.
func doctorCommandsAction(_ core.Options) core.Result {
	c := cli.Core()
	if c == nil {
		return cli.Err("CLI core not initialised")
	}

	cli.Section("commands")
	for _, path := range c.Commands() {
		if core.Contains(path, "/") {
			continue // top-level commands only
		}
		description := ""
		if r := c.Command(path); r.OK {
			if cmd, ok := r.Value.(*core.Command); ok {
				// Description is an i18n key by contract; cli.T resolves it when
				// the catalog has it and passes literal descriptions through.
				description = cli.T(cmd.Description)
			}
		}
		cli.Println("  %s %s", successStyle.Render(path), dimStyle.Render(description))
	}
	return core.Ok(nil)
}
