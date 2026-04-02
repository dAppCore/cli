package pkgcmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T, dir, name string) string {
	t.Helper()

	repoPath := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(repoPath, 0755))

	gitCommand(t, repoPath, "init")
	gitCommand(t, repoPath, "config", "user.email", "test@test.com")
	gitCommand(t, repoPath, "config", "user.name", "Test")
	gitCommand(t, repoPath, "commit", "--allow-empty", "-m", "initial")

	return repoPath
}

func capturePkgStreams(t *testing.T, fn func()) (string, string) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	require.NoError(t, err)
	rErr, wErr, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = wOut
	os.Stderr = wErr

	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	fn()

	require.NoError(t, wOut.Close())
	require.NoError(t, wErr.Close())

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	_, err = io.Copy(&stdout, rOut)
	require.NoError(t, err)
	_, err = io.Copy(&stderr, rErr)
	require.NoError(t, err)

	return stdout.String(), stderr.String()
}

func TestCheckRepoSafety_Clean(t *testing.T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "clean-repo")

	blocked, reasons := checkRepoSafety(repoPath)
	assert.False(t, blocked)
	assert.Empty(t, reasons)
}

func TestCheckRepoSafety_UncommittedChanges(t *testing.T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "dirty-repo")

	require.NoError(t, os.WriteFile(filepath.Join(repoPath, "new.txt"), []byte("data"), 0644))

	blocked, reasons := checkRepoSafety(repoPath)
	assert.True(t, blocked)
	assert.NotEmpty(t, reasons)
	assert.Contains(t, reasons[0], "uncommitted changes")
}

func TestCheckRepoSafety_Stash(t *testing.T) {
	tmp := t.TempDir()
	repoPath := setupTestRepo(t, tmp, "stash-repo")

	require.NoError(t, os.WriteFile(filepath.Join(repoPath, "stash.txt"), []byte("data"), 0644))
	gitCommand(t, repoPath, "add", ".")
	gitCommand(t, repoPath, "stash")

	blocked, reasons := checkRepoSafety(repoPath)
	assert.True(t, blocked)

	found := false
	for _, r := range reasons {
		if strings.Contains(r, "stash") {
			found = true
		}
	}
	assert.True(t, found, "expected stash warning in reasons: %v", reasons)
}

func TestRunPkgRemove_RemovesRegistryEntry_Good(t *testing.T) {
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

	require.NoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))

	oldwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(oldwd))
	})

	require.NoError(t, runPkgRemove("core-alpha", false))

	_, err = os.Stat(repoPath)
	assert.True(t, os.IsNotExist(err))

	updated, err := os.ReadFile(filepath.Join(tmp, "repos.yaml"))
	require.NoError(t, err)
	assert.NotContains(t, string(updated), "core-alpha")
	assert.Contains(t, string(updated), "core-beta")
}

func TestRunPkgRemove_Bad_BlockedWarningsGoToStderr(t *testing.T) {
	tmp := t.TempDir()

	registry := strings.TrimSpace(`
org: host-uk
base_path: .
repos:
  core-alpha:
    type: foundation
    description: Alpha package
`) + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))

	repoPath := filepath.Join(tmp, "core-alpha")
	require.NoError(t, os.MkdirAll(repoPath, 0755))
	gitCommand(t, repoPath, "init")
	gitCommand(t, repoPath, "config", "user.email", "test@test.com")
	gitCommand(t, repoPath, "config", "user.name", "Test")
	commitGitRepo(t, repoPath, "file.txt", "v1\n", "initial")
	require.NoError(t, os.WriteFile(filepath.Join(repoPath, "file.txt"), []byte("v2\n"), 0644))

	withWorkingDir(t, tmp)

	stdout, stderr := capturePkgStreams(t, func() {
		err := runPkgRemove("core-alpha", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unresolved changes")
	})

	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "Cannot remove core-alpha")
	assert.Contains(t, stderr, "uncommitted changes")
	assert.Contains(t, stderr, "Resolve the issues above or use --force to override.")
}
