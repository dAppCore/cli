package pkgcmd

import (
	"context"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/cli/pkg/i18n"
	coreio "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

func pkgInstallAction(opts core.Options) core.Result {
	repoArg := opts.String("_arg")
	if repoArg == "" {
		return cli.Err(i18n.T("cmd.pkg.error.repo_required"))
	}
	targetDir := opts.String("dir")
	addToReg := opts.Bool("add")
	if r := runPkgInstall(repoArg, targetDir, addToReg); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func runPkgInstall(repoArg, targetDirectory string, addToRegistry bool) core.Result {
	ctx := context.Background()

	// Parse org/repo argument.
	parts := core.Split(repoArg, "/")
	if len(parts) != 2 {
		return cli.Err(i18n.T("cmd.pkg.error.invalid_repo_format"))
	}
	org, repoName := parts[0], parts[1]

	// Determine target directory from registry or default.
	if targetDirectory == "" {
		if registryPath, err := repos.FindRegistry(coreio.Local); err == nil {
			if registry, err := repos.LoadRegistry(coreio.Local, registryPath); err == nil {
				targetDirectory = registry.BasePath
				if targetDirectory == "" {
					targetDirectory = "./packages"
				}
				if !core.PathIsAbs(targetDirectory) {
					targetDirectory = core.Path(core.PathDir(registryPath), targetDirectory)
				}
			}
		}
		if targetDirectory == "" {
			targetDirectory = "."
		}
	}

	if core.HasPrefix(targetDirectory, "~/") {
		homeResult := core.UserHomeDir()
		if homeResult.OK {
			targetDirectory = core.Path(homeResult.Value.(string), targetDirectory[2:])
		}
	}

	repoPath := core.Path(targetDirectory, repoName)

	if coreio.Local.Exists(core.Path(repoPath, ".git")) {
		cli.Println("%s %s", dimStyle.Render(i18n.Label("skip")), i18n.T("cmd.pkg.install.already_exists", map[string]string{"Name": repoName, "Path": repoPath}))
		return core.Ok(nil)
	}

	if err := coreio.Local.EnsureDir(targetDirectory); err != nil {
		return cli.Wrap(err, i18n.T("i18n.fail.create", "directory"))
	}

	cli.Println("%s %s/%s", dimStyle.Render(i18n.T("cmd.pkg.install.installing_label")), org, repoName)
	cli.Println("%s %s", dimStyle.Render(i18n.Label("target")), repoPath)
	cli.Blank()

	cli.Print("  %s... ", dimStyle.Render(i18n.T("common.status.cloning")))
	cloneResult := gitClone(ctx, org, repoName, repoPath)
	if !cloneResult.OK {
		cli.Println("%s", errorStyle.Render("x "+cloneResult.Error()))
		return cloneResult
	}
	cli.Println("%s", successStyle.Render("ok"))

	if addToRegistry {
		if r := addToRegistryFile(org, repoName); !r.OK {
			cli.Println("  %s %s: %s", errorStyle.Render("x"), i18n.T("cmd.pkg.install.add_to_registry"), r.Error())
		} else {
			cli.Println("  %s %s", successStyle.Render("ok"), i18n.T("cmd.pkg.install.added_to_registry"))
		}
	}

	cli.Blank()
	cli.Println("%s %s", successStyle.Render(i18n.T("i18n.done.install")), i18n.T("cmd.pkg.install.installed", map[string]string{"Name": repoName}))

	return core.Ok(nil)
}

func addToRegistryFile(org, repoName string) core.Result {
	registryPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return cli.Err(i18n.T("cmd.pkg.error.no_repos_yaml"))
	}

	registry, err := repos.LoadRegistry(coreio.Local, registryPath)
	if err != nil {
		return core.Fail(err)
	}

	if _, exists := registry.Get(repoName); exists {
		return core.Ok(nil)
	}

	content, err := coreio.Local.Read(registryPath)
	if err != nil {
		return core.Fail(err)
	}

	repoType := detectRepoType(repoName)
	entry := cli.Sprintf("\n  %s:\n    type: %s\n    description: (installed via core pkg install)\n",
		repoName, repoType)

	content += entry
	if err := coreio.Local.Write(registryPath, content); err != nil {
		return core.Fail(err)
	}
	return core.Ok(nil)
}

func detectRepoType(name string) string {
	lowerName := core.Lower(name)
	if core.Contains(lowerName, "-mod-") || core.HasSuffix(lowerName, "-mod") {
		return "module"
	}
	if core.Contains(lowerName, "-plug-") || core.HasSuffix(lowerName, "-plug") {
		return "plugin"
	}
	if core.Contains(lowerName, "-services-") || core.HasSuffix(lowerName, "-services") {
		return "service"
	}
	if core.Contains(lowerName, "-website-") || core.HasSuffix(lowerName, "-website") {
		return "website"
	}
	if core.HasPrefix(lowerName, "core-") {
		return "package"
	}
	return "package"
}
