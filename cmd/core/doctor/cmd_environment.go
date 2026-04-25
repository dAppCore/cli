package doctor

import (
	"os" // Note: AX-6 — os.UserHomeDir; released core has no home-directory primitive.

	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/core"
	"dappco.re/go/i18n"
	io "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

var environmentFS = (&core.Fs{}).New("/")

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
		if environmentFS.Stat(keyPath).OK {
			return true
		}
	}

	return false
}

// checkGitHubCLI checks if the GitHub CLI is authenticated.
// Returns true when GitHub CLI reports an authenticated session.
func checkGitHubCLI() bool {
	if !(core.App{}).Find("gh", "GitHub CLI").OK {
		return false
	}
	return cli.GhAuthenticated()
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
				if environmentFS.Stat(core.Path(repoPath, ".git")).OK {
					cloned++
				}
			}
			cli.Println("  %s %s", successStyle.Render("✓"), i18n.T("cmd.doctor.repos_cloned", map[string]any{"Cloned": cloned, "Total": len(allRepos)}))
		}
	} else {
		cli.Println("  %s %s", dimStyle.Render("○"), i18n.T("cmd.doctor.no_repos_yaml"))
	}
}
