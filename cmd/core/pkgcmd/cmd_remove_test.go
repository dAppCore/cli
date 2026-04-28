package pkgcmd

import (
	"bytes"
	. "dappco.re/go"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func setupTestRepo(t *T, dir, name string) string {
	t.Helper()

	repoPath := filepath.Join(dir, name)
	RequireNoError(t, os.MkdirAll(repoPath, 0755))

	gitCommand(t, repoPath, "init")
	gitCommand(t, repoPath, "config", "user.email", "test@test.com")
	gitCommand(t, repoPath, "config", "user.name", "Test")
	gitCommand(t, repoPath, "commit", "--allow-empty", "-m", "initial")

	return repoPath
}

func capturePkgStreams(t *T, fn func()) (string, string) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	RequireNoError(t, err)
	rErr, wErr, err := os.Pipe()
	RequireNoError(t, err)

	os.Stdout = wOut
	os.Stderr = wErr

	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	fn()
	RequireNoError(t, wOut.Close())
	RequireNoError(t, wErr.Close())

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, err = io.Copy(&stdout, rOut)
	RequireNoError(t, err)
	_, err = io.Copy(&stderr, rErr)
	RequireNoError(t, err)

	return stdout.String(), stderr.String()
}

func TestCheckRepoSafety_Clean(t *T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "clean-repo")

	blocked, reasons := checkRepoSafety(repoPath)
	AssertFalse(t, blocked)
	AssertEmpty(t, reasons)
}

func TestCheckRepoSafety_UncommittedChanges(t *T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "dirty-repo")
	RequireNoError(t, os.WriteFile(filepath.Join(repoPath, "new.txt"), []byte("data"), 0644))

	blocked, reasons := checkRepoSafety(repoPath)
	AssertTrue(t, blocked)
	AssertNotEmpty(t, reasons)
	AssertContains(t, reasons[0], "uncommitted changes")
}

func TestCheckRepoSafety_Stash(t *T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "stash-repo")
	RequireNoError(t, os.WriteFile(filepath.Join(repoPath, "stash.txt"), []byte("data"), 0644))
	gitCommand(t, repoPath, "add", ".")
	gitCommand(t, repoPath, "stash")

	blocked, reasons := checkRepoSafety(repoPath)
	AssertTrue(t, blocked)

	found := false
	for _, r := range reasons {
		if strings.Contains(r, "stash") {
			found = true
		}
	}
	AssertTrue(t, found, Sprintf("expected stash warning in reasons: %v", reasons))
}

func TestRunPkgRemove_RemovesRegistryEntry_Good(t *T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "core-alpha")

	registry := strings.TrimSpace(`
version: 1
org: host-uk
base_path: .
repos:
  core-alpha:
    type: foundation
    description: Alpha package
  core-beta:
    type: module
    description: Beta package
`) + "\n"
	RequireNoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))

	oldwd, err := os.Getwd()
	RequireNoError(t, err)
	RequireNoError(t, os.Chdir(tmp))
	t.Cleanup(func() {
		RequireNoError(t, os.Chdir(oldwd))
	})
	RequireNoError(t, runPkgRemove("core-alpha", false))

	_, err = os.Stat(repoPath)
	AssertTrue(t, os.IsNotExist(err))

	updated, err := os.ReadFile(filepath.Join(tmp, "repos.yaml"))
	RequireNoError(t, err)
	AssertNotContains(t, string(updated), "core-alpha")
	AssertContains(t, string(updated), "core-beta")
}

func TestRunPkgRemove_Bad_BlockedWarningsGoToStderr(t *T) {
	tmp := t.TempDir()

	registry := strings.TrimSpace(`
org: host-uk
base_path: .
repos:
  core-alpha:
    type: foundation
    description: Alpha package
`) + "\n"
	RequireNoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))

	repoPath := filepath.Join(tmp, "core-alpha")
	RequireNoError(t, os.MkdirAll(repoPath, 0755))
	gitCommand(t, repoPath, "init")
	gitCommand(t, repoPath, "config", "user.email", "test@test.com")
	gitCommand(t, repoPath, "config", "user.name", "Test")
	commitGitRepo(t, repoPath, "file.txt", "v1\n", "initial")
	RequireNoError(t, os.WriteFile(filepath.Join(repoPath, "file.txt"), []byte("v2\n"), 0644))

	withWorkingDir(t, tmp)

	stdout, stderr := capturePkgStreams(t, func() {
		err := runPkgRemove("core-alpha", false)
		RequireTrue(t, err != nil, "RequireError")
		AssertContains(t, err.Error(), "unresolved changes")
	})
	AssertEmpty(t, stdout)
	AssertContains(t, stderr, "Cannot remove core-alpha")
	AssertContains(t, stderr, "uncommitted changes")
	AssertContains(t, stderr, "Resolve the issues above or use --force to override.")
}
