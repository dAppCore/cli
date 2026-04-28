package pkgcmd

import (
	"context"
	. "dappco.re/go"
	"os"
	"path/filepath"
)

func TestRunPkgInstall_AllowsRepoShorthand_Good(t *T) {
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
	RequireNoError(t, err)
	AssertEqual(t, "host-uk", gotOrg)
	AssertEqual(t, "core-api", gotRepo)
	AssertEqual(t, filepath.Join(targetDir, "core-api"), gotPath)
	_, err = os.Stat(targetDir)
	RequireNoError(t, err)
}

func TestRunPkgInstall_AllowsExplicitOrgRepo_Good(t *T) {
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
	RequireNoError(t, err)
	AssertEqual(t, "myorg", gotOrg)
	AssertEqual(t, "core-api", gotRepo)
	AssertEqual(t, filepath.Join(targetDir, "core-api"), gotPath)
}

func TestRunPkgInstall_InvalidRepoFormat_Bad(t *T) {
	err := runPkgInstall("a/b/c", t.TempDir(), false)
	RequireTrue(t, err != nil, "RequireError")
	AssertContains(t, err.Error(), "invalid repo format")
}

func TestParsePkgInstallSource_Good(t *T) {
	t.Run("default org and repo", func(t *T) {
		org, repo, ref, err := parsePkgInstallSource("core-api")
		RequireNoError(t, err)
		AssertEqual(t, "host-uk", org)
		AssertEqual(t, "core-api", repo)
		AssertEmpty(t, ref)
	})

	t.Run("explicit org and ref", func(t *T) {
		org, repo, ref, err := parsePkgInstallSource("myorg/core-api@v1.2.3")
		RequireNoError(t, err)
		AssertEqual(t, "myorg", org)
		AssertEqual(t, "core-api", repo)
		AssertEqual(t, "v1.2.3", ref)
	})
}

func TestRunPkgInstall_WithRef_UsesRefClone_Good(t *T) {
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
	RequireNoError(t, err)
	AssertEqual(t, "myorg", gotOrg)
	AssertEqual(t, "core-api", gotRepo)
	AssertEqual(t, filepath.Join(targetDir, "core-api"), gotPath)
	AssertEqual(t, "v1.2.3", gotRef)
}
