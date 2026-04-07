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
	"dappco.re/go/core/i18n"
	"github.com/spf13/cobra"
)

// AddDoctorCommands registers the 'doctor' command and all subcommands.
//
//	doctor.AddDoctorCommands(rootCmd)
func AddDoctorCommands(root *cobra.Command) {
	doctorCmd.Short = i18n.T("cmd.doctor.short")
	doctorCmd.Long = i18n.T("cmd.doctor.long")
	root.AddCommand(doctorCmd)
}
