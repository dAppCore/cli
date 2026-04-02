package pkgcmd

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"forge.lthn.ai/core/go-i18n"
	coreio "forge.lthn.ai/core/go-io"
	"forge.lthn.ai/core/go-scm/repos"
	"github.com/spf13/cobra"
)

// addPkgListCommand adds the 'pkg list' command.
func addPkgListCommand(parent *cobra.Command) {
	var format string
	listCmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("cmd.pkg.list.short"),
		Long:  i18n.T("cmd.pkg.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			return runPkgList(format)
		},
	}

	listCmd.Flags().StringVar(&format, "format", "table", "Output format: table or json")
	parent.AddCommand(listCmd)
}

type pkgListEntry struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Installed   bool   `json:"installed"`
	Path        string `json:"path"`
}

type pkgListReport struct {
	Format    string         `json:"format"`
	Total     int            `json:"total"`
	Installed int            `json:"installed"`
	Missing   int            `json:"missing"`
	Packages  []pkgListEntry `json:"packages"`
}

func runPkgList(format string) error {
	regPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return errors.New(i18n.T("cmd.pkg.error.no_repos_yaml_workspace"))
	}

	reg, err := repos.LoadRegistry(coreio.Local, regPath)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.load", "registry"), err)
	}

	basePath := reg.BasePath
	if basePath == "" {
		basePath = "."
	}
	if !filepath.IsAbs(basePath) {
		basePath = filepath.Join(filepath.Dir(regPath), basePath)
	}

	allRepos := reg.List()
	if len(allRepos) == 0 {
		fmt.Println(i18n.T("cmd.pkg.list.no_packages"))
		return nil
	}

	slices.SortFunc(allRepos, func(a, b *repos.Repo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	var entries []pkgListEntry
	var installed, missing int
	for _, r := range allRepos {
		repoPath := filepath.Join(basePath, r.Name)
		exists := coreio.Local.Exists(filepath.Join(repoPath, ".git"))
		if exists {
			installed++
		} else {
			missing++
		}

		desc := r.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		if desc == "" {
			desc = i18n.T("cmd.pkg.no_description")
		}

		entries = append(entries, pkgListEntry{
			Name:        r.Name,
			Description: desc,
			Installed:   exists,
			Path:        repoPath,
		})
	}

	if format == "json" {
		report := pkgListReport{
			Format:    "json",
			Total:     len(entries),
			Installed: installed,
			Missing:   missing,
			Packages:  entries,
		}

		out, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format package list: %w", err)
		}

		fmt.Println(string(out))
		return nil
	}

	if format != "table" {
		return fmt.Errorf("unsupported format %q: expected table or json", format)
	}

	fmt.Printf("%s\n\n", repoNameStyle.Render(i18n.T("cmd.pkg.list.title")))

	for _, entry := range entries {
		status := successStyle.Render("✓")
		if !entry.Installed {
			status = dimStyle.Render("○")
		}

		desc := entry.Description
		if !entry.Installed {
			desc = dimStyle.Render(desc)
		}

		fmt.Printf("  %s %s\n", status, repoNameStyle.Render(entry.Name))
		fmt.Printf("      %s\n", desc)
	}

	fmt.Println()
	fmt.Printf("%s %s\n", dimStyle.Render(i18n.Label("total")), i18n.T("cmd.pkg.list.summary", map[string]int{"Installed": installed, "Missing": missing}))

	if missing > 0 {
		fmt.Printf("\n%s %s\n", i18n.T("cmd.pkg.list.install_missing"), dimStyle.Render("core setup"))
	}

	return nil
}

var updateAll bool
var updateFormat string

// addPkgUpdateCommand adds the 'pkg update' command.
func addPkgUpdateCommand(parent *cobra.Command) {
	updateCmd := &cobra.Command{
		Use:   "update [packages...]",
		Short: i18n.T("cmd.pkg.update.short"),
		Long:  i18n.T("cmd.pkg.update.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			return runPkgUpdate(args, updateAll, format)
		},
	}

	updateCmd.Flags().BoolVar(&updateAll, "all", false, i18n.T("cmd.pkg.update.flag.all"))
	updateCmd.Flags().StringVar(&updateFormat, "format", "table", "Output format: table or json")

	parent.AddCommand(updateCmd)
}

type pkgUpdateEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Installed bool   `json:"installed"`
	Status    string `json:"status"`
	Output    string `json:"output,omitempty"`
}

type pkgUpdateReport struct {
	Format    string           `json:"format"`
	Total     int              `json:"total"`
	Installed int              `json:"installed"`
	Missing   int              `json:"missing"`
	Updated   int              `json:"updated"`
	UpToDate  int              `json:"upToDate"`
	Failed    int              `json:"failed"`
	Packages  []pkgUpdateEntry `json:"packages"`
}

func runPkgUpdate(packages []string, all bool, format string) error {
	regPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return errors.New(i18n.T("cmd.pkg.error.no_repos_yaml"))
	}

	reg, err := repos.LoadRegistry(coreio.Local, regPath)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.load", "registry"), err)
	}

	basePath := reg.BasePath
	if basePath == "" {
		basePath = "."
	}
	if !filepath.IsAbs(basePath) {
		basePath = filepath.Join(filepath.Dir(regPath), basePath)
	}

	jsonOutput := strings.EqualFold(format, "json")
	var toUpdate []string
	if all || len(packages) == 0 {
		for _, r := range reg.List() {
			toUpdate = append(toUpdate, r.Name)
		}
	} else {
		toUpdate = packages
	}

	if !jsonOutput {
		fmt.Printf("%s %s\n\n", dimStyle.Render(i18n.T("cmd.pkg.update.update_label")), i18n.T("cmd.pkg.update.updating", map[string]int{"Count": len(toUpdate)}))
	}

	var updated, upToDate, skipped, failed int
	var entries []pkgUpdateEntry
	for _, name := range toUpdate {
		repoPath := filepath.Join(basePath, name)

		if _, err := coreio.Local.List(filepath.Join(repoPath, ".git")); err != nil {
			if !jsonOutput {
				fmt.Printf("  %s %s (%s)\n", dimStyle.Render("○"), name, i18n.T("cmd.pkg.update.not_installed"))
			}
			if jsonOutput {
				entries = append(entries, pkgUpdateEntry{
					Name:      name,
					Path:      repoPath,
					Installed: false,
					Status:    "missing",
				})
			}
			skipped++
			continue
		}

		if !jsonOutput {
			fmt.Printf("  %s %s... ", dimStyle.Render("↓"), name)
		}

		cmd := exec.Command("git", "-C", repoPath, "pull", "--ff-only")
		output, err := cmd.CombinedOutput()
		if err != nil {
			if !jsonOutput {
				fmt.Printf("%s\n", errorStyle.Render("✗"))
				fmt.Printf("      %s\n", strings.TrimSpace(string(output)))
			}
			if jsonOutput {
				entries = append(entries, pkgUpdateEntry{
					Name:      name,
					Path:      repoPath,
					Installed: true,
					Status:    "failed",
					Output:    strings.TrimSpace(string(output)),
				})
			}
			failed++
			continue
		}

		if strings.Contains(string(output), "Already up to date") {
			if !jsonOutput {
				fmt.Printf("%s\n", dimStyle.Render(i18n.T("common.status.up_to_date")))
			}
			if jsonOutput {
				entries = append(entries, pkgUpdateEntry{
					Name:      name,
					Path:      repoPath,
					Installed: true,
					Status:    "up_to_date",
					Output:    strings.TrimSpace(string(output)),
				})
			}
			upToDate++
		} else {
			if !jsonOutput {
				fmt.Printf("%s\n", successStyle.Render("✓"))
			}
			if jsonOutput {
				entries = append(entries, pkgUpdateEntry{
					Name:      name,
					Path:      repoPath,
					Installed: true,
					Status:    "updated",
					Output:    strings.TrimSpace(string(output)),
				})
			}
			updated++
		}
	}

	if jsonOutput {
		report := pkgUpdateReport{
			Format:    "json",
			Total:     len(toUpdate),
			Installed: updated + upToDate + failed,
			Missing:   skipped,
			Updated:   updated,
			UpToDate:  upToDate,
			Failed:    failed,
			Packages:  entries,
		}
		return printPkgUpdateJSON(report)
	}

	fmt.Println()
	fmt.Printf("%s %s\n",
		dimStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.update.summary", map[string]int{"Updated": updated + upToDate, "Skipped": skipped, "Failed": failed}))

	return nil
}

func printPkgUpdateJSON(report pkgUpdateReport) error {
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.format", "update results"), err)
	}

	fmt.Println(string(out))
	return nil
}

// addPkgOutdatedCommand adds the 'pkg outdated' command.
func addPkgOutdatedCommand(parent *cobra.Command) {
	var format string
	outdatedCmd := &cobra.Command{
		Use:   "outdated",
		Short: i18n.T("cmd.pkg.outdated.short"),
		Long:  i18n.T("cmd.pkg.outdated.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			return runPkgOutdated(format)
		},
	}

	outdatedCmd.Flags().StringVar(&format, "format", "table", i18n.T("cmd.pkg.outdated.flag.format"))
	parent.AddCommand(outdatedCmd)
}

type pkgOutdatedEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Behind    int    `json:"behind"`
	UpToDate  bool   `json:"upToDate"`
	Installed bool   `json:"installed"`
}

type pkgOutdatedReport struct {
	Format    string             `json:"format"`
	Total     int                `json:"total"`
	Installed int                `json:"installed"`
	Missing   int                `json:"missing"`
	Outdated  int                `json:"outdated"`
	UpToDate  int                `json:"upToDate"`
	Packages  []pkgOutdatedEntry `json:"packages"`
}

func runPkgOutdated(format string) error {
	regPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return errors.New(i18n.T("cmd.pkg.error.no_repos_yaml"))
	}

	reg, err := repos.LoadRegistry(coreio.Local, regPath)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.load", "registry"), err)
	}

	basePath := reg.BasePath
	if basePath == "" {
		basePath = "."
	}
	if !filepath.IsAbs(basePath) {
		basePath = filepath.Join(filepath.Dir(regPath), basePath)
	}

	jsonOutput := strings.EqualFold(format, "json")
	if !jsonOutput {
		fmt.Printf("%s %s\n\n", dimStyle.Render(i18n.T("cmd.pkg.outdated.outdated_label")), i18n.T("common.progress.checking_updates"))
	}

	var installed, outdated, upToDate, notInstalled int
	var entries []pkgOutdatedEntry

	for _, r := range reg.List() {
		repoPath := filepath.Join(basePath, r.Name)

		if !coreio.Local.Exists(filepath.Join(repoPath, ".git")) {
			notInstalled++
			if jsonOutput {
				entries = append(entries, pkgOutdatedEntry{
					Name:      r.Name,
					Path:      repoPath,
					Behind:    0,
					UpToDate:  false,
					Installed: false,
				})
			}
			continue
		}
		installed++

		// Fetch updates
		_ = exec.Command("git", "-C", repoPath, "fetch", "--quiet").Run()

		// Check if behind
		cmd := exec.Command("git", "-C", repoPath, "rev-list", "--count", "HEAD..@{u}")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		count := strings.TrimSpace(string(output))
		behind := 0
		if count != "" {
			fmt.Sscanf(count, "%d", &behind)
		}
		if count != "0" {
			if !jsonOutput {
				fmt.Printf("  %s %s (%s)\n",
					errorStyle.Render("↓"), repoNameStyle.Render(r.Name), i18n.T("cmd.pkg.outdated.commits_behind", map[string]string{"Count": count}))
			}
			outdated++
			if jsonOutput {
				entries = append(entries, pkgOutdatedEntry{
					Name:      r.Name,
					Path:      repoPath,
					Behind:    behind,
					UpToDate:  false,
					Installed: true,
				})
			}
		} else {
			upToDate++
			if jsonOutput {
				entries = append(entries, pkgOutdatedEntry{
					Name:      r.Name,
					Path:      repoPath,
					Behind:    0,
					UpToDate:  true,
					Installed: true,
				})
			}
		}
	}

	if jsonOutput {
		report := pkgOutdatedReport{
			Format:    "json",
			Total:     len(reg.List()),
			Installed: installed,
			Missing:   notInstalled,
			Outdated:  outdated,
			UpToDate:  upToDate,
			Packages:  entries,
		}
		return printPkgOutdatedJSON(report)
	}

	fmt.Println()
	if outdated == 0 {
		fmt.Printf("%s %s\n", successStyle.Render(i18n.T("i18n.done.update")), i18n.T("cmd.pkg.outdated.all_up_to_date"))
	} else {
		fmt.Printf("%s %s\n",
			dimStyle.Render(i18n.Label("summary")), i18n.T("cmd.pkg.outdated.summary", map[string]int{"Outdated": outdated, "UpToDate": upToDate}))
		fmt.Printf("\n%s %s\n", i18n.T("cmd.pkg.outdated.update_with"), dimStyle.Render("core pkg update --all"))
	}

	return nil
}

func printPkgOutdatedJSON(report pkgOutdatedReport) error {
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.format", "outdated results"), err)
	}

	fmt.Println(string(out))
	return nil
}
