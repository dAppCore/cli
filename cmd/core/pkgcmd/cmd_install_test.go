package pkgcmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunPkgInstall_AllowsRepoShorthand_Good(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "packages")

	originalGitClone := gitClone
	t.Cleanup(func() {
		gitClone = originalGitClone
	})

	var gotOrg, gotRepo, gotPath string
	gitClone = func(_ context.Context, org, repoName, repoPath string) error {
		gotOrg = org
		gotRepo = repoName
		gotPath = repoPath
		return nil
	}

	err := runPkgInstall("core-api", targetDir, false)
	require.NoError(t, err)

	assert.Equal(t, "host-uk", gotOrg)
	assert.Equal(t, "core-api", gotRepo)
	assert.Equal(t, filepath.Join(targetDir, "core-api"), gotPath)
	_, err = os.Stat(targetDir)
	require.NoError(t, err)
}

func TestRunPkgInstall_AllowsExplicitOrgRepo_Good(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "packages")

	originalGitClone := gitClone
	t.Cleanup(func() {
		gitClone = originalGitClone
	})

	var gotOrg, gotRepo, gotPath string
	gitClone = func(_ context.Context, org, repoName, repoPath string) error {
		gotOrg = org
		gotRepo = repoName
		gotPath = repoPath
		return nil
	}

	err := runPkgInstall("myorg/core-api", targetDir, false)
	require.NoError(t, err)

	assert.Equal(t, "myorg", gotOrg)
	assert.Equal(t, "core-api", gotRepo)
	assert.Equal(t, filepath.Join(targetDir, "core-api"), gotPath)
}

func TestRunPkgInstall_InvalidRepoFormat_Bad(t *testing.T) {
	err := runPkgInstall("a/b/c", t.TempDir(), false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid repo format")
}
