package pkgcmd

import (
	"os/exec"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/cli/pkg/i18n"
	coreio "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

func pkgListAction(_ core.Options) core.Result {
	if err := runPkgList(); err != nil {
		return core.Result{Value: err, OK: false}
	}
	return core.Result{OK: true}
}

func runPkgList() error {
	registryPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return cli.Err(i18n.T("cmd.pkg.error.no_repos_yaml_workspace"))
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

	allRepos := registry.List()
	if len(allRepos) == 0 {
		cli.Println("%s", i18n.T("cmd.pkg.list.no_packages"))
		return nil
	}

	cli.Println("%s\n", repoNameStyle.Render(i18n.T("cmd.pkg.list.title")))

	var installed, missing int
	for _, repo := range allRepos {
		repoPath := core.Path(basePath, repo.Name)
		exists := coreio.Local.Exists(core.Path(repoPath, ".git"))
		if exists {
			installed++
		} else {
			missing++
		}

		status := successStyle.Render("ok")
		if !exists {
			status = dimStyle.Render("o")
		}

		description := repo.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}
		if description == "" {
			description = dimStyle.Render(i18n.T("cmd.pkg.no_description"))
		}

		cli.Println("  %s %s", status, repoNameStyle.Render(repo.Name))
		cli.Println("      %s", description)
	}

	cli.Blank()
	cli.Println("%s %s", dimStyle.Render(i18n.Label("total")), i18n.T("cmd.pkg.list.summary", map[string]int{"Installed": installed, "Missing": missing}))

	if missing > 0 {
		cli.Println("\n%s %s", i18n.T("cmd.pkg.list.install_missing"), dimStyle.Render("core setup"))
	}

	return nil
}

func pkgUpdateAction(opts core.Options) core.Result {
	all := opts.Bool("all")
	pkg := opts.String("_arg")
	var packages []string
	if pkg != "" {
		packages = append(packages, pkg)
	}
	if !all && len(packages) == 0 {
		return core.Result{Value: cli.Err(i18n.T("cmd.pkg.error.specify_package")), OK: false}
	}
	if err := runPkgUpdate(packages, all); err != nil {
		return core.Result{Value: err, OK: false}
	}
	return core.Result{OK: true}
}

func runPkgUpdate(packages []string, all bool) error {
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

	var toUpdate []string
	if all {
		for _, repo := range registry.List() {
			toUpdate = append(toUpdate, repo.Name)
		}
	} else {
		toUpdate = packages
	}

	cli.Println("%s %s\n", dimStyle.Render(i18n.T("cmd.pkg.update.update_label")), i18n.T("cmd.pkg.update.updating", map[string]int{"Count": len(toUpdate)}))

	var updated, skipped, failed int
	for _, name := range toUpdate {
		repoPath := core.Path(basePath, name)

		if _, err := coreio.Local.List(core.Path(repoPath, ".git")); err != nil {
			cli.Println("  %s %s (%s)", dimStyle.Render("o"), name, i18n.T("cmd.pkg.update.not_installed"))
			skipped++
			continue
		}

		cli.Print("  %s %s... ", dimStyle.Render("v"), name)

		proc := exec.Command("git", "-C", repoPath, "pull", "--ff-only") // TODO: migrate to c.Process()
		output, err := proc.CombinedOutput()
		if err != nil {
			cli.Println("%s", errorStyle.Render("x"))
			cli.Println("      %s", core.Trim(string(output)))
			failed++
			continue
		}

		if core.Contains(string(output), "Already up to date") {
			cli.Println("%s", dimStyle.Render(i18n.T("common.status.up_to_date")))
		} else {
			cli.Println("%s", successStyle.Render("ok"))
		}
		updated++
	}

	cli.Blank()
	cli.Println("%s %s",
		dimStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.update.summary", map[string]int{"Updated": updated, "Skipped": skipped, "Failed": failed}))

	return nil
}

func pkgOutdatedAction(_ core.Options) core.Result {
	if err := runPkgOutdated(); err != nil {
		return core.Result{Value: err, OK: false}
	}
	return core.Result{OK: true}
}

func runPkgOutdated() error {
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

	cli.Println("%s %s\n", dimStyle.Render(i18n.T("cmd.pkg.outdated.outdated_label")), i18n.T("common.progress.checking_updates"))

	var outdated, upToDate, notInstalled int

	for _, repo := range registry.List() {
		repoPath := core.Path(basePath, repo.Name)

		if !coreio.Local.Exists(core.Path(repoPath, ".git")) {
			notInstalled++
			continue
		}

		// Fetch updates silently.
		_ = exec.Command("git", "-C", repoPath, "fetch", "--quiet").Run() // TODO: migrate to c.Process()

		// Check commit count behind upstream.
		proc := exec.Command("git", "-C", repoPath, "rev-list", "--count", "HEAD..@{u}") // TODO: migrate to c.Process()
		output, err := proc.Output()
		if err != nil {
			continue
		}

		commitCount := core.Trim(string(output))
		if commitCount != "0" {
			cli.Println("  %s %s (%s)",
				errorStyle.Render("v"), repoNameStyle.Render(repo.Name), i18n.T("cmd.pkg.outdated.commits_behind", map[string]string{"Count": commitCount}))
			outdated++
		} else {
			upToDate++
		}
	}

	_ = notInstalled

	cli.Blank()
	if outdated == 0 {
		cli.Println("%s %s", successStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.outdated.all_up_to_date"))
	} else {
		cli.Println("%s %s",
			dimStyle.Render(i18n.Label("summary")), i18n.T("cmd.pkg.outdated.summary", map[string]int{"Outdated": outdated, "UpToDate": upToDate}))
		cli.Println("\n%s %s", i18n.T("cmd.pkg.outdated.update_with"), dimStyle.Render("core pkg update --all"))
	}

	return nil
}
