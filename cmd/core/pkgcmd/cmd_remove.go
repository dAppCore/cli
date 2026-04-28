// cmd_remove.go implements the 'pkg remove' command with safety checks.
//
// Before removing a package, it verifies:
// 1. No uncommitted changes exist
// 2. No unpushed branches exist
// This prevents accidental data loss from agents or tools that might
// attempt to remove packages without cleaning up first.
package pkgcmd

import (
	"os/exec"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/cli/pkg/i18n"
	coreio "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

func pkgRemoveAction(opts core.Options) core.Result {
	name := opts.String("_arg")
	if name == "" {
		return core.Fail(cli.Err(i18n.T("cmd.pkg.error.repo_required")))
	}
	force := opts.Bool("force")
	if err := runPkgRemove(name, force); err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}

func runPkgRemove(name string, force bool) error {
	// Find package path via registry.
	registryPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return cli.Err(i18n.T("cmd.pkg.error.no_repos_yaml"))
	}

	registry, err := repos.LoadRegistry(coreio.Local, registryPath)
	if err != nil {
		return cli.Wrap(err, i18n.T("i18n.fail.load", "registry"))
	}

	basePath := registry.BasePath
	if basePath == "" {
		basePath = "."
	}
	if !core.PathIsAbs(basePath) {
		basePath = core.Path(core.PathDir(registryPath), basePath)
	}

	repoPath := core.Path(basePath, name)

	if !coreio.Local.IsDir(core.Path(repoPath, ".git")) {
		return cli.Err("package %s is not installed at %s", name, repoPath)
	}

	if !force {
		blocked, reasons := checkRepoSafety(repoPath)
		if blocked {
			cli.Println("%s Cannot remove %s:", errorStyle.Render("Blocked:"), repoNameStyle.Render(name))
			for _, reason := range reasons {
				cli.Println("  %s %s", errorStyle.Render("*"), reason)
			}
			cli.Println("\nResolve the issues above or use --force to override.")
			return cli.Err("package has unresolved changes")
		}
	}

	// Remove the directory.
	cli.Print("%s %s... ", dimStyle.Render("Removing"), repoNameStyle.Render(name))

	if err := coreio.Local.DeleteAll(repoPath); err != nil {
		cli.Println("%s", errorStyle.Render("x "+err.Error()))
		return err
	}

	cli.Println("%s", successStyle.Render("ok"))
	return nil
}

// checkRepoSafety checks a git repo for uncommitted changes and unpushed branches.
//
//	blocked, reasons := checkRepoSafety("/path/to/repo")
//	if blocked { fmt.Println(reasons) }
func checkRepoSafety(repoPath string) (blocked bool, reasons []string) {
	// Check for uncommitted changes (staged, unstaged, untracked).
	proc := exec.Command("git", "-C", repoPath, "status", "--porcelain") // TODO: migrate to c.Process()
	output, err := proc.Output()
	if err == nil && core.Trim(string(output)) != "" {
		lines := core.Split(core.Trim(string(output)), "\n")
		blocked = true
		reasons = append(reasons, cli.Sprintf("has %d uncommitted changes", len(lines)))
	}

	// Check for unpushed commits on current branch.
	proc = exec.Command("git", "-C", repoPath, "log", "--oneline", "@{u}..HEAD") // TODO: migrate to c.Process()
	output, err = proc.Output()
	if err == nil && core.Trim(string(output)) != "" {
		lines := core.Split(core.Trim(string(output)), "\n")
		blocked = true
		reasons = append(reasons, cli.Sprintf("has %d unpushed commits on current branch", len(lines)))
	}

	// Check all local branches for unpushed work.
	proc = exec.Command("git", "-C", repoPath, "branch", "--no-merged", "origin/HEAD") // TODO: migrate to c.Process()
	output, _ = proc.Output()
	if trimmedOutput := core.Trim(string(output)); trimmedOutput != "" {
		branches := core.Split(trimmedOutput, "\n")
		var unmerged []string
		for _, branchName := range branches {
			branchName = core.Trim(branchName)
			branchName = core.TrimPrefix(branchName, "* ")
			if branchName != "" {
				unmerged = append(unmerged, branchName)
			}
		}
		if len(unmerged) > 0 {
			blocked = true
			reasons = append(reasons, cli.Sprintf("has %d unmerged branches: %s",
				len(unmerged), core.Join(", ", unmerged...)))
		}
	}

	// Check for stashed changes.
	proc = exec.Command("git", "-C", repoPath, "stash", "list") // TODO: migrate to c.Process()
	output, err = proc.Output()
	if err == nil && core.Trim(string(output)) != "" {
		lines := core.Split(core.Trim(string(output)), "\n")
		blocked = true
		reasons = append(reasons, cli.Sprintf("has %d stashed entries", len(lines)))
	}

	return blocked, reasons
}
