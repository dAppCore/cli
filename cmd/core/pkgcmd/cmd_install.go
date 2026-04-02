package pkgcmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge.lthn.ai/core/go-i18n"
	coreio "forge.lthn.ai/core/go-io"
	"forge.lthn.ai/core/go-scm/repos"
	"github.com/spf13/cobra"
)

import (
	"errors"
)

var (
	installTargetDir string
	installAddToReg  bool
)

var errInvalidPkgInstallSource = errors.New("invalid repo format: use org/repo or org/repo@ref")

// addPkgInstallCommand adds the 'pkg install' command.
func addPkgInstallCommand(parent *cobra.Command) {
	installCmd := &cobra.Command{
		Use:   "install [org/]repo[@ref]",
		Short: i18n.T("cmd.pkg.install.short"),
		Long:  i18n.T("cmd.pkg.install.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New(i18n.T("cmd.pkg.error.repo_required"))
			}
			return runPkgInstall(args[0], installTargetDir, installAddToReg)
		},
	}

	installCmd.Flags().StringVar(&installTargetDir, "dir", "", i18n.T("cmd.pkg.install.flag.dir"))
	installCmd.Flags().BoolVar(&installAddToReg, "add", false, i18n.T("cmd.pkg.install.flag.add"))

	parent.AddCommand(installCmd)
}

func runPkgInstall(repoArg, targetDir string, addToRegistry bool) error {
	ctx := context.Background()

	// Parse repo shorthand:
	// - repoName -> defaults to host-uk/repoName
	// - org/repo  -> uses the explicit org
	org, repoName, ref, err := parsePkgInstallSource(repoArg)
	if err != nil {
		return err
	}

	// Determine target directory
	if targetDir == "" {
		if regPath, err := repos.FindRegistry(coreio.Local); err == nil {
			if reg, err := repos.LoadRegistry(coreio.Local, regPath); err == nil {
				targetDir = reg.BasePath
				if targetDir == "" {
					targetDir = "./packages"
				}
				if !filepath.IsAbs(targetDir) {
					targetDir = filepath.Join(filepath.Dir(regPath), targetDir)
				}
			}
		}
		if targetDir == "" {
			targetDir = "."
		}
	}

	if strings.HasPrefix(targetDir, "~/") {
		home, _ := os.UserHomeDir()
		targetDir = filepath.Join(home, targetDir[2:])
	}

	repoPath := filepath.Join(targetDir, repoName)

	if coreio.Local.Exists(filepath.Join(repoPath, ".git")) {
		fmt.Printf("%s %s\n", dimStyle.Render(i18n.Label("skip")), i18n.T("cmd.pkg.install.already_exists", map[string]string{"Name": repoName, "Path": repoPath}))
		return nil
	}

	if err := coreio.Local.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.create", "directory"), err)
	}

	fmt.Printf("%s %s/%s\n", dimStyle.Render(i18n.T("cmd.pkg.install.installing_label")), org, repoName)
	if ref != "" {
		fmt.Printf("%s %s\n", dimStyle.Render(i18n.Label("ref")), ref)
	}
	fmt.Printf("%s %s\n", dimStyle.Render(i18n.Label("target")), repoPath)
	fmt.Println()

	fmt.Printf("  %s... ", dimStyle.Render(i18n.T("common.status.cloning")))
	if ref == "" {
		err = gitClone(ctx, org, repoName, repoPath)
	} else {
		err = gitCloneRef(ctx, org, repoName, repoPath, ref)
	}
	if err != nil {
		fmt.Printf("%s\n", errorStyle.Render("✗ "+err.Error()))
		return err
	}
	fmt.Printf("%s\n", successStyle.Render("✓"))

	if addToRegistry {
		if err := addToRegistryFile(org, repoName); err != nil {
			fmt.Printf("  %s %s: %s\n", errorStyle.Render("✗"), i18n.T("cmd.pkg.install.add_to_registry"), err)
		} else {
			fmt.Printf("  %s %s\n", successStyle.Render("✓"), i18n.T("cmd.pkg.install.added_to_registry"))
		}
	}

	fmt.Println()
	fmt.Printf("%s %s\n", successStyle.Render(i18n.T("i18n.done.install")), i18n.T("cmd.pkg.install.installed", map[string]string{"Name": repoName}))

	return nil
}

func parsePkgInstallSource(repoArg string) (org, repoName, ref string, err error) {
	org = "host-uk"
	repoName = strings.TrimSpace(repoArg)
	if repoName == "" {
		return "", "", "", errors.New("repository argument required")
	}

	if at := strings.LastIndex(repoName, "@"); at >= 0 {
		ref = strings.TrimSpace(repoName[at+1:])
		repoName = strings.TrimSpace(repoName[:at])
		if ref == "" || repoName == "" {
			return "", "", "", errInvalidPkgInstallSource
		}
	}

	if strings.Contains(repoName, "/") {
		parts := strings.Split(repoName, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", "", errInvalidPkgInstallSource
		}
		org, repoName = parts[0], parts[1]
	}

	if strings.Contains(repoName, "/") {
		return "", "", "", errInvalidPkgInstallSource
	}

	return org, repoName, ref, nil
}

func addToRegistryFile(org, repoName string) error {
	regPath, err := repos.FindRegistry(coreio.Local)
	if err != nil {
		return errors.New(i18n.T("cmd.pkg.error.no_repos_yaml"))
	}

	reg, err := repos.LoadRegistry(coreio.Local, regPath)
	if err != nil {
		return err
	}

	if _, exists := reg.Get(repoName); exists {
		return nil
	}

	content, err := coreio.Local.Read(regPath)
	if err != nil {
		return err
	}

	repoType := detectRepoType(repoName)
	entry := fmt.Sprintf("\n  %s:\n    type: %s\n    description: (installed via core pkg install)\n",
		repoName, repoType)

	content += entry
	return coreio.Local.Write(regPath, content)
}

func clonePackageAtRef(ctx context.Context, org, repo, path, ref string) error {
	if ghAuthenticated() {
		httpsURL := fmt.Sprintf("https://github.com/%s/%s.git", org, repo)
		args := []string{"repo", "clone", httpsURL, path, "--", "--branch", ref, "--single-branch"}
		cmd := exec.CommandContext(ctx, "gh", args...)
		output, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		errStr := strings.TrimSpace(string(output))
		if strings.Contains(errStr, "already exists") {
			return errors.New(errStr)
		}
	}

	args := []string{"clone", "--branch", ref, "--single-branch", fmt.Sprintf("git@github.com:%s/%s.git", org, repo), path}
	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(strings.TrimSpace(string(output)))
	}
	return nil
}

func detectRepoType(name string) string {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "-mod-") || strings.HasSuffix(lower, "-mod") {
		return "module"
	}
	if strings.Contains(lower, "-plug-") || strings.HasSuffix(lower, "-plug") {
		return "plugin"
	}
	if strings.Contains(lower, "-services-") || strings.HasSuffix(lower, "-services") {
		return "service"
	}
	if strings.Contains(lower, "-website-") || strings.HasSuffix(lower, "-website") {
		return "website"
	}
	if strings.HasPrefix(lower, "core-") {
		return "package"
	}
	return "package"
}
