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

func TestParsePkgInstallSource_Good(t *testing.T) {
	t.Run("default org and repo", func(t *testing.T) {
		org, repo, ref, err := parsePkgInstallSource("core-api")
		require.NoError(t, err)
		assert.Equal(t, "host-uk", org)
		assert.Equal(t, "core-api", repo)
		assert.Empty(t, ref)
	})

	t.Run("explicit org and ref", func(t *testing.T) {
		org, repo, ref, err := parsePkgInstallSource("myorg/core-api@v1.2.3")
		require.NoError(t, err)
		assert.Equal(t, "myorg", org)
		assert.Equal(t, "core-api", repo)
		assert.Equal(t, "v1.2.3", ref)
	})
}

func TestRunPkgInstall_WithRef_UsesRefClone_Good(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "packages")

	originalGitCloneRef := gitCloneRef
	t.Cleanup(func() {
		gitCloneRef = originalGitCloneRef
	})

	var gotOrg, gotRepo, gotPath, gotRef string
	gitCloneRef = func(_ context.Context, org, repoName, repoPath, ref string) error {
		gotOrg = org
		gotRepo = repoName
		gotPath = repoPath
		gotRef = ref
		return nil
	}

	err := runPkgInstall("myorg/core-api@v1.2.3", targetDir, false)
	require.NoError(t, err)

	assert.Equal(t, "myorg", gotOrg)
	assert.Equal(t, "core-api", gotRepo)
	assert.Equal(t, filepath.Join(targetDir, "core-api"), gotPath)
	assert.Equal(t, "v1.2.3", gotRef)
}
