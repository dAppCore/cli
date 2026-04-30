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
package doctor

import (
	"dappco.re/go"
)

// AddDoctorCommands registers the 'doctor' command and all subcommands.
//
//	doctor.AddDoctorCommands(c)
func AddDoctorCommands(c *core.Core) core.Result {
	if r := c.Command("doctor", core.Command{
		Description: "Check development environment health",
		Action:      doctorAction,
	}); !r.OK {
		return r
	}
	return core.Ok(nil)
}
