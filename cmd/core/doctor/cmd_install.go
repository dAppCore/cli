package doctor

import (
	"runtime"

	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/i18n"
)

// printInstallInstructions prints operating-system-specific installation instructions.
func printInstallInstructions() {
	switch runtime.GOOS {
	case "darwin":
		cli.Println("  %s", i18n.T("cmd.doctor.install_macos"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_macos_cask"))
	case "linux":
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_header"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_git"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_gh"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_php"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_node"))
		cli.Println("  %s", i18n.T("cmd.doctor.install_linux_pnpm"))
	default:
		cli.Println("  %s", i18n.T("cmd.doctor.install_other"))
	}
}
