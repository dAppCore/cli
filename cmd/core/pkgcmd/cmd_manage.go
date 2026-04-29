package pkgcmd

import (
	"syscall"

	"dappco.re/go"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/cli/pkg/i18n"
	coreio "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

func pkgListAction(_ core.Options) core.Result {
	if r := runPkgList(); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func runPkgList() core.Result {
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
		return core.Ok(nil)
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

	return core.Ok(nil)
}

func pkgUpdateAction(opts core.Options) core.Result {
	all := opts.Bool("all")
	pkg := opts.String("_arg")
	var packages []string
	if pkg != "" {
		packages = append(packages, pkg)
	}
	if !all && len(packages) == 0 {
		return cli.Err(i18n.T("cmd.pkg.error.specify_package"))
	}
	if r := runPkgUpdate(packages, all); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func runPkgUpdate(packages []string, all bool) core.Result {
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

		result := pkgRunGit(repoPath, "pull", "--ff-only")
		output := pkgProcessOutput(result.Value)
		if !result.OK {
			cli.Println("%s", errorStyle.Render("x"))
			cli.Println("      %s", core.Trim(output))
			failed++
			continue
		}

		if core.Contains(output, "Already up to date") {
			cli.Println("%s", dimStyle.Render(i18n.T("common.status.up_to_date")))
		} else {
			cli.Println("%s", successStyle.Render("ok"))
		}
		updated++
	}

	cli.Blank()
	cli.Println("%s %s",
		dimStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.update.summary", map[string]int{"Updated": updated, "Skipped": skipped, "Failed": failed}))

	return core.Ok(nil)
}

func pkgOutdatedAction(_ core.Options) core.Result {
	if r := runPkgOutdated(); !r.OK {
		return r
	}
	return core.Ok(nil)
}

func runPkgOutdated() core.Result {
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
		if r := pkgRunGit(repoPath, "fetch", "--quiet"); !r.OK {
			cli.LogWarn("failed to fetch package updates", "repo", repo.Name, "err", r.Error())
		}

		// Check commit count behind upstream.
		result := pkgRunGit(repoPath, "rev-list", "--count", "HEAD..@{u}")
		if !result.OK {
			continue
		}

		commitCount := core.Trim(pkgProcessOutput(result.Value))
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

	return core.Ok(nil)
}

func pkgRunGit(dir string, args ...string) core.Result {
	return pkgRunProcess(dir, "git", args...)
}

func pkgProcessOutput(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case nil:
		return ""
	default:
		return core.Sprint(v)
	}
}

func pkgRunProcess(dir, command string, args ...string) core.Result {
	commandResult := pkgFindExecutable(command)
	if !commandResult.OK {
		return commandResult
	}
	commandPath := commandResult.Value.(string)

	var pipe [2]int
	if err := syscall.Pipe(pipe[:]); err != nil {
		return core.Fail(err)
	}
	readFD, writeFD := pipe[0], pipe[1]
	defer syscall.Close(readFD)

	argv := append([]string{commandPath}, args...)
	pid, err := syscall.ForkExec(commandPath, argv, &syscall.ProcAttr{
		Dir:   dir,
		Env:   core.Environ(),
		Files: []uintptr{0, uintptr(writeFD), uintptr(writeFD)},
	})
	syscall.Close(writeFD)
	if err != nil {
		return core.Fail(err)
	}

	out := core.NewBuilder()
	buf := make([]byte, 4096)
	for {
		n, readErr := syscall.Read(readFD, buf)
		if n > 0 {
			out.WriteString(string(buf[:n]))
		}
		if readErr != nil {
			if readErr == syscall.EINTR {
				continue
			}
			break
		}
		if n == 0 {
			break
		}
	}

	var status syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &status, 0, nil); err != nil {
		return core.Fail(err)
	}
	output := out.String()
	if status.ExitStatus() == 0 {
		return core.Ok(output)
	}
	if output != "" {
		return core.Fail(core.NewError(output))
	}
	return core.Fail(core.E("pkg.process", core.Sprintf("%s exited with status %d", command, status.ExitStatus()), nil))
}

func pkgFindExecutable(command string) core.Result {
	if command == "" {
		return core.Fail(core.NewError("empty command"))
	}
	if core.Contains(command, string(core.PathSeparator)) {
		if r := core.Stat(command); r.OK {
			return core.Ok(command)
		}
		return core.Fail(core.E("pkg.process", core.Concat("command not found: ", command), nil))
	}
	for _, dir := range core.Split(core.Getenv("PATH"), string(core.PathListSeparator)) {
		if dir == "" {
			continue
		}
		candidate := core.PathJoin(dir, command)
		if r := core.Stat(candidate); r.OK {
			return core.Ok(candidate)
		}
	}
	return core.Fail(core.E("pkg.process", core.Concat("command not found: ", command), nil))
}
