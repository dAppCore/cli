// Package pkgcmd provides package management commands for core-* repos.
package pkgcmd

import (
	"dappco.re/go/core"
	"dappco.re/go/cli/pkg/cli"
)

// Style and utility aliases
var (
	repoNameStyle   = cli.RepoStyle
	successStyle    = cli.SuccessStyle
	errorStyle      = cli.ErrorStyle
	dimStyle        = cli.DimStyle
	ghAuthenticated = cli.GhAuthenticated
	gitClone        = cli.GitClone
	gitCloneRef     = cli.GitCloneRef
)

// AddPkgCommands adds the 'pkg' command and subcommands for package management.
func AddPkgCommands(c *core.Core) {
	c.Command("pkg/search", core.Command{
		Description: "Search GitHub org for packages",
		Action:      pkgSearchAction,
	})
	c.Command("pkg/install", core.Command{
		Description: "Install a package from GitHub",
		Action:      pkgInstallAction,
	})
	c.Command("pkg/list", core.Command{
		Description: "List installed packages",
		Action:      pkgListAction,
	})
	c.Command("pkg/update", core.Command{
		Description: "Update installed packages",
		Action:      pkgUpdateAction,
	})
	c.Command("pkg/outdated", core.Command{
		Description: "Check for outdated packages",
		Action:      pkgOutdatedAction,
	})
	c.Command("pkg/remove", core.Command{
		Description: "Remove a package (with safety checks)",
		Action:      pkgRemoveAction,
	})
}
