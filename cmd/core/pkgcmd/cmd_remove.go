// cmd_remove.go implements the 'pkg remove' command with safety checks.
//
// Before removing a package, it verifies:
// 1. No uncommitted changes exist
// 2. No unpushed branches exist
// This prevents accidental data loss from agents or tools that might
// attempt to remove packages without cleaning up first.
package pkgcmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge.lthn.ai/core/go-i18n"
	coreio "forge.lthn.ai/core/go-io"
	"forge.lthn.ai/core/go-scm/repos"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var removeForce bool

func addPkgRemoveCommand(parent *cobra.Command) {
	removeCmd := &cobra.Command{
		Use:   "remove <package>",
		Short: "Remove a package (with safety checks)",
		Long: `Removes a package directory after verifying it has no uncommitted
changes or unpushed branches. Use --force to skip safety checks.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New(i18n.T("cmd.pkg.error.repo_required"))
			}
			return runPkgRemove(args[0], removeForce)
		},
	}

	removeCmd.Flags().BoolVar(&removeForce, "force", false, "Skip safety checks (dangerous)")

	parent.AddCommand(removeCmd)
}

func runPkgRemove(name string, force bool) error {
	// Find package path via registry
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

	repoPath := filepath.Join(basePath, name)

	if !coreio.Local.IsDir(filepath.Join(repoPath, ".git")) {
		return fmt.Errorf("package %s is not installed at %s", name, repoPath)
	}

	if !force {
		blocked, reasons := checkRepoSafety(repoPath)
		if blocked {
			fmt.Fprintf(os.Stderr, "%s Cannot remove %s:\n", errorStyle.Render("Blocked:"), repoNameStyle.Render(name))
			for _, r := range reasons {
				fmt.Fprintf(os.Stderr, "  %s %s\n", errorStyle.Render("·"), r)
			}
			fmt.Fprintln(os.Stderr, "\nResolve the issues above or use --force to override.")
			return errors.New("package has unresolved changes")
		}
	}

	// Remove the directory
	fmt.Printf("%s %s... ", dimStyle.Render("Removing"), repoNameStyle.Render(name))

	if err := coreio.Local.DeleteAll(repoPath); err != nil {
		fmt.Printf("%s\n", errorStyle.Render("x "+err.Error()))
		return err
	}

	if err := removeRepoFromRegistry(regPath, name); err != nil {
		return fmt.Errorf("removed %s from disk, but failed to update registry: %w", name, err)
	}

	fmt.Printf("%s\n", successStyle.Render("ok"))
	return nil
}

func removeRepoFromRegistry(regPath, name string) error {
	content, err := coreio.Local.Read(regPath)
	if err != nil {
		return err
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return fmt.Errorf("failed to parse registry file: %w", err)
	}
	if len(doc.Content) == 0 {
		return errors.New("registry file is empty")
	}

	root := doc.Content[0]
	reposNode := mappingValue(root, "repos")
	if reposNode == nil {
		return errors.New("registry file has no repos section")
	}
	if reposNode.Kind != yaml.MappingNode {
		return errors.New("registry repos section is malformed")
	}

	if removeMappingEntry(reposNode, name) {
		out, err := yaml.Marshal(&doc)
		if err != nil {
			return fmt.Errorf("failed to format registry file: %w", err)
		}
		return coreio.Local.Write(regPath, string(out))
	}

	return nil
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	return nil
}

func removeMappingEntry(node *yaml.Node, key string) bool {
	if node == nil || node.Kind != yaml.MappingNode {
		return false
	}

	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != key {
			continue
		}
		node.Content = append(node.Content[:i], node.Content[i+2:]...)
		return true
	}

	return false
}

// checkRepoSafety checks a git repo for uncommitted changes and unpushed branches.
func checkRepoSafety(repoPath string) (blocked bool, reasons []string) {
	// Check for uncommitted changes (staged, unstaged, untracked)
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		blocked = true
		reasons = append(reasons, fmt.Sprintf("has %d uncommitted changes", len(lines)))
	}

	// Check for unpushed commits on current branch
	cmd = exec.Command("git", "-C", repoPath, "log", "--oneline", "@{u}..HEAD")
	output, err = cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		blocked = true
		reasons = append(reasons, fmt.Sprintf("has %d unpushed commits on current branch", len(lines)))
	}

	// Check all local branches for unpushed work
	cmd = exec.Command("git", "-C", repoPath, "branch", "--no-merged", "origin/HEAD")
	output, _ = cmd.Output()
	if trimmed := strings.TrimSpace(string(output)); trimmed != "" {
		branches := strings.Split(trimmed, "\n")
		var unmerged []string
		for _, b := range branches {
			b = strings.TrimSpace(b)
			b = strings.TrimPrefix(b, "* ")
			if b != "" {
				unmerged = append(unmerged, b)
			}
		}
		if len(unmerged) > 0 {
			blocked = true
			reasons = append(reasons, fmt.Sprintf("has %d unmerged branches: %s",
				len(unmerged), strings.Join(unmerged, ", ")))
		}
	}

	// Check for stashed changes
	cmd = exec.Command("git", "-C", repoPath, "stash", "list")
	output, err = cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) != "" {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		blocked = true
		reasons = append(reasons, fmt.Sprintf("has %d stashed entries", len(lines)))
	}

	return blocked, reasons
}
