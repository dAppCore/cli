package pkgcmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"forge.lthn.ai/core/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func capturePkgOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	require.NoError(t, w.Close())

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	return buf.String()
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	oldwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(oldwd))
	})
}

func writeTestRegistry(t *testing.T, dir string) {
	t.Helper()

	registry := strings.TrimSpace(`
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

	require.NoError(t, os.WriteFile(filepath.Join(dir, "repos.yaml"), []byte(registry), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "core-alpha", ".git"), 0755))
}

func gitCommand(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v failed: %s", args, string(out))
	return string(out)
}

func commitGitRepo(t *testing.T, dir, filename, content, message string) {
	t.Helper()

	require.NoError(t, os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644))
	gitCommand(t, dir, "add", filename)
	gitCommand(t, dir, "commit", "-m", message)
}

func setupOutdatedRegistry(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	remoteDir := filepath.Join(tmp, "remote.git")
	gitCommand(t, tmp, "init", "--bare", remoteDir)

	seedDir := filepath.Join(tmp, "seed")
	require.NoError(t, os.MkdirAll(seedDir, 0755))
	gitCommand(t, seedDir, "init")
	gitCommand(t, seedDir, "config", "user.email", "test@test.com")
	gitCommand(t, seedDir, "config", "user.name", "Test")
	commitGitRepo(t, seedDir, "repo.txt", "v1\n", "initial")
	gitCommand(t, seedDir, "remote", "add", "origin", remoteDir)
	gitCommand(t, seedDir, "push", "-u", "origin", "master")

	freshDir := filepath.Join(tmp, "core-fresh")
	gitCommand(t, tmp, "clone", remoteDir, freshDir)

	staleDir := filepath.Join(tmp, "core-stale")
	gitCommand(t, tmp, "clone", remoteDir, staleDir)

	commitGitRepo(t, seedDir, "repo.txt", "v2\n", "second")
	gitCommand(t, seedDir, "push")
	gitCommand(t, freshDir, "pull", "--ff-only")

	registry := strings.TrimSpace(`
org: host-uk
base_path: .
repos:
  core-fresh:
    type: foundation
    description: Fresh package
  core-stale:
    type: module
    description: Stale package
  core-missing:
    type: module
    description: Missing package
`) + "\n"

	require.NoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))
	return tmp
}

func TestRunPkgList_Good(t *testing.T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgList("table")
		require.NoError(t, err)
	})

	assert.Contains(t, out, "core-alpha")
	assert.Contains(t, out, "core-beta")
	assert.Contains(t, out, "core setup")
}

func TestRunPkgList_JSON(t *testing.T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgList("json")
		require.NoError(t, err)
	})

	var report pkgListReport
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &report))
	assert.Equal(t, "json", report.Format)
	assert.Equal(t, 2, report.Total)
	assert.Equal(t, 1, report.Installed)
	assert.Equal(t, 1, report.Missing)
	require.Len(t, report.Packages, 2)
	assert.Equal(t, "core-alpha", report.Packages[0].Name)
	assert.True(t, report.Packages[0].Installed)
	assert.Equal(t, filepath.Join(tmp, "core-alpha"), report.Packages[0].Path)
	assert.Equal(t, "core-beta", report.Packages[1].Name)
	assert.False(t, report.Packages[1].Installed)
}

func TestRunPkgList_UnsupportedFormat(t *testing.T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	err := runPkgList("yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestRunPkgOutdated_JSON(t *testing.T) {
	tmp := setupOutdatedRegistry(t)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgOutdated("json")
		require.NoError(t, err)
	})

	var report pkgOutdatedReport
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &report))
	assert.Equal(t, "json", report.Format)
	assert.Equal(t, 3, report.Total)
	assert.Equal(t, 2, report.Installed)
	assert.Equal(t, 1, report.Missing)
	assert.Equal(t, 1, report.Outdated)
	assert.Equal(t, 1, report.UpToDate)
	require.Len(t, report.Packages, 3)

	var staleFound, freshFound, missingFound bool
	for _, pkg := range report.Packages {
		switch pkg.Name {
		case "core-stale":
			staleFound = true
			assert.True(t, pkg.Installed)
			assert.False(t, pkg.UpToDate)
			assert.Equal(t, 1, pkg.Behind)
		case "core-fresh":
			freshFound = true
			assert.True(t, pkg.Installed)
			assert.True(t, pkg.UpToDate)
			assert.Equal(t, 0, pkg.Behind)
		case "core-missing":
			missingFound = true
			assert.False(t, pkg.Installed)
			assert.False(t, pkg.UpToDate)
			assert.Equal(t, 0, pkg.Behind)
		}
	}

	assert.True(t, staleFound)
	assert.True(t, freshFound)
	assert.True(t, missingFound)
}

func TestRenderPkgSearchResults_ShowsMetadata(t *testing.T) {
	out := capturePkgOutput(t, func() {
		renderPkgSearchResults([]ghRepo{
			{
				FullName:       "host-uk/core-alpha",
				Name:           "core-alpha",
				Description:    "Alpha package",
				Visibility:     "private",
				StargazerCount: 42,
				PrimaryLanguage: ghLanguage{
					Name: "Go",
				},
				UpdatedAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			},
		})
	})

	assert.Contains(t, out, "host-uk/core-alpha")
	assert.Contains(t, out, "Alpha package")
	assert.Contains(t, out, "42 stars")
	assert.Contains(t, out, "Go")
	assert.Contains(t, out, "updated 2h ago")
}

func TestRunPkgSearch_RespectsLimitWithCachedResults(t *testing.T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	c, err := cache.New(nil, filepath.Join(tmp, ".core", "cache"), 0)
	require.NoError(t, err)
	require.NoError(t, c.Set(cache.GitHubReposKey("host-uk"), []ghRepo{
		{
			FullName:       "host-uk/core-alpha",
			Name:           "core-alpha",
			Description:    "Alpha package",
			Visibility:     "public",
			UpdatedAt:      time.Now().Add(-time.Hour).Format(time.RFC3339),
			StargazerCount: 1,
			PrimaryLanguage: ghLanguage{
				Name: "Go",
			},
		},
		{
			FullName:       "host-uk/core-beta",
			Name:           "core-beta",
			Description:    "Beta package",
			Visibility:     "public",
			UpdatedAt:      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			StargazerCount: 2,
			PrimaryLanguage: ghLanguage{
				Name: "Go",
			},
		},
	}))

	out := capturePkgOutput(t, func() {
		err := runPkgSearch("host-uk", "*", "", 1, false, "table")
		require.NoError(t, err)
	})

	assert.Contains(t, out, "core-alpha")
	assert.NotContains(t, out, "core-beta")
}

func TestRunPkgUpdate_NoArgs_UpdatesAll(t *testing.T) {
	tmp := setupOutdatedRegistry(t)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgUpdate(nil, false)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "updating")
	assert.Contains(t, out, "core-fresh")
	assert.Contains(t, out, "core-stale")
}
