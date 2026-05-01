package doctor

import (
	"runtime"

	"dappco.re/go/cli/pkg/cli"
)

// printInstallInstructions prints operating-system-specific installation instructions.
func printInstallInstructions() {
	switch runtime.GOOS {
	case "darwin":
		cli.Println("  %s", cli.T("cmd.doctor.install_macos"))
		cli.Println("  %s", cli.T("cmd.doctor.install_macos_cask"))
	case "linux":
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_header"))
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_git"))
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_gh"))
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_php"))
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_node"))
		cli.Println("  %s", cli.T("cmd.doctor.install_linux_pnpm"))
	default:
		cli.Println("  %s", cli.T("cmd.doctor.install_other"))
	}
}
