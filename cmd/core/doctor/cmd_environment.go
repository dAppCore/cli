package doctor

import (
	"os"
	"os/exec"

	"dappco.re/go/core"
	"dappco.re/go/core/cli/pkg/cli"
	"dappco.re/go/core/i18n"
	io "dappco.re/go/core/io"
	"dappco.re/go/core/scm/repos"
)

// checkGitHubSSH checks if SSH keys exist for GitHub access.
// Returns true if any standard SSH key file exists in ~/.ssh/.
func checkGitHubSSH() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	sshDirectory := core.Path(home, ".ssh")
	keyPatterns := []string{"id_rsa", "id_ed25519", "id_ecdsa", "id_dsa"}

	for _, keyName := range keyPatterns {
		keyPath := core.Path(sshDirectory, keyName)
		if _, err := os.Stat(keyPath); err == nil {
			return true
		}
	}

	return false
}

// checkGitHubCLI checks if the GitHub CLI is authenticated.
// Returns true when 'gh auth status' output contains "Logged in to".
func checkGitHubCLI() bool {
	proc := exec.Command("gh", "auth", "status")
	output, _ := proc.CombinedOutput()
	return core.Contains(string(output), "Logged in to")
}

// checkWorkspace checks for repos.yaml and counts cloned repos.
func checkWorkspace() {
	registryPath, err := repos.FindRegistry(io.Local)
	if err == nil {
		cli.Println("  %s %s", successStyle.Render("✓"), i18n.T("cmd.doctor.repos_yaml_found", map[string]any{"Path": registryPath}))

		registry, err := repos.LoadRegistry(io.Local, registryPath)
		if err == nil {
			basePath := registry.BasePath
			if basePath == "" {
				basePath = "./packages"
			}
			if !core.PathIsAbs(basePath) {
				basePath = core.Path(core.PathDir(registryPath), basePath)
			}
			if core.HasPrefix(basePath, "~/") {
				home, _ := os.UserHomeDir()
				basePath = core.Path(home, basePath[2:])
			}

			// Count existing repos.
			allRepos := registry.List()
			var cloned int
			for _, repo := range allRepos {
				repoPath := core.Path(basePath, repo.Name)
				if _, err := os.Stat(core.Path(repoPath, ".git")); err == nil {
					cloned++
				}
			}
			cli.Println("  %s %s", successStyle.Render("✓"), i18n.T("cmd.doctor.repos_cloned", map[string]any{"Cloned": cloned, "Total": len(allRepos)}))
		}
	} else {
		cli.Println("  %s %s", dimStyle.Render("○"), i18n.T("cmd.doctor.no_repos_yaml"))
	}
}
