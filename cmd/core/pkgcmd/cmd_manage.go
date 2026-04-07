package pkgcmd

import (
	"os/exec"

	"dappco.re/go/core"
	"dappco.re/go/core/cli/pkg/cli"
	"dappco.re/go/core/i18n"
	coreio "dappco.re/go/core/io"
	"dappco.re/go/core/scm/repos"
	"github.com/spf13/cobra"
)

// addPkgListCommand adds the 'pkg list' command.
func addPkgListCommand(parent *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("cmd.pkg.list.short"),
		Long:  i18n.T("cmd.pkg.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPkgList()
		},
	}

	parent.AddCommand(listCmd)
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

		status := successStyle.Render("✓")
		if !exists {
			status = dimStyle.Render("○")
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

var updateAll bool

// addPkgUpdateCommand adds the 'pkg update' command.
func addPkgUpdateCommand(parent *cobra.Command) {
	updateCmd := &cobra.Command{
		Use:   "update [packages...]",
		Short: i18n.T("cmd.pkg.update.short"),
		Long:  i18n.T("cmd.pkg.update.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !updateAll && len(args) == 0 {
				return cli.Err(i18n.T("cmd.pkg.error.specify_package"))
			}
			return runPkgUpdate(args, updateAll)
		},
	}

	updateCmd.Flags().BoolVar(&updateAll, "all", false, i18n.T("cmd.pkg.update.flag.all"))

	parent.AddCommand(updateCmd)
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
			cli.Println("  %s %s (%s)", dimStyle.Render("○"), name, i18n.T("cmd.pkg.update.not_installed"))
			skipped++
			continue
		}

		cli.Print("  %s %s... ", dimStyle.Render("↓"), name)

		proc := exec.Command("git", "-C", repoPath, "pull", "--ff-only")
		output, err := proc.CombinedOutput()
		if err != nil {
			cli.Println("%s", errorStyle.Render("✗"))
			cli.Println("      %s", core.Trim(string(output)))
			failed++
			continue
		}

		if core.Contains(string(output), "Already up to date") {
			cli.Println("%s", dimStyle.Render(i18n.T("common.status.up_to_date")))
		} else {
			cli.Println("%s", successStyle.Render("✓"))
		}
		updated++
	}

	cli.Blank()
	cli.Println("%s %s",
		dimStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.update.summary", map[string]int{"Updated": updated, "Skipped": skipped, "Failed": failed}))

	return nil
}

// addPkgOutdatedCommand adds the 'pkg outdated' command.
func addPkgOutdatedCommand(parent *cobra.Command) {
	outdatedCmd := &cobra.Command{
		Use:   "outdated",
		Short: i18n.T("cmd.pkg.outdated.short"),
		Long:  i18n.T("cmd.pkg.outdated.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPkgOutdated()
		},
	}

	parent.AddCommand(outdatedCmd)
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
		_ = exec.Command("git", "-C", repoPath, "fetch", "--quiet").Run()

		// Check commit count behind upstream.
		proc := exec.Command("git", "-C", repoPath, "rev-list", "--count", "HEAD..@{u}")
		output, err := proc.Output()
		if err != nil {
			continue
		}

		commitCount := core.Trim(string(output))
		if commitCount != "0" {
			cli.Println("  %s %s (%s)",
				errorStyle.Render("↓"), repoNameStyle.Render(repo.Name), i18n.T("cmd.pkg.outdated.commits_behind", map[string]string{"Count": commitCount}))
			outdated++
		} else {
			upToDate++
		}
	}

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
