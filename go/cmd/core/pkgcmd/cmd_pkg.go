// Package pkgcmd provides package management commands for core-* repos.
package pkgcmd

import (
	"dappco.re/go"
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
func AddPkgCommands(c *core.Core) core.Result {
	if r := c.Command("pkg/search", core.Command{
		Description: "Search GitHub org for packages",
		Action:      pkgSearchAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("pkg/install", core.Command{
		Description: "Install a package from GitHub",
		Action:      pkgInstallAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("pkg/list", core.Command{
		Description: "List installed packages",
		Action:      pkgListAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("pkg/update", core.Command{
		Description: "Update installed packages",
		Action:      pkgUpdateAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("pkg/outdated", core.Command{
		Description: "Check for outdated packages",
		Action:      pkgOutdatedAction,
	}); !r.OK {
		return r
	}
	if r := c.Command("pkg/remove", core.Command{
		Description: "Remove a package (with safety checks)",
		Action:      pkgRemoveAction,
	}); !r.OK {
		return r
	}
	return core.Ok(nil)
}
