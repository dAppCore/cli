package pkgcmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"time"

	. "dappco.re/go"
	"dappco.re/go/cache"
)

func capturePkgOutput(t *T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	RequireNoError(t, err)
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	fn()
	RequireNoError(t, w.Close())

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	RequireNoError(t, err)
	return buf.String()
}

func withWorkingDir(t *T, dir string) {
	t.Helper()

	oldwd, err := os.Getwd()
	RequireNoError(t, err)
	RequireNoError(t, os.Chdir(dir))

	t.Cleanup(func() {
		RequireNoError(t, os.Chdir(oldwd))
	})
}

func writeTestRegistry(t *T, dir string) {
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
	RequireNoError(t, os.WriteFile(filepath.Join(dir, "repos.yaml"), []byte(registry), 0644))
	RequireNoError(t, os.MkdirAll(filepath.Join(dir, "core-alpha", ".git"), 0755))
}

func gitCommand(t *T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	RequireNoError(t, err, Sprintf("git %v failed: %s", args, string(out)))
	return string(out)
}

func commitGitRepo(t *T, dir, filename, content, message string) {
	t.Helper()
	RequireNoError(t, os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644))
	gitCommand(t, dir, "add", filename)
	gitCommand(t, dir, "commit", "-m", message)
}

func setupOutdatedRegistry(t *T) string {
	t.Helper()

	tmp := t.TempDir()

	remoteDir := filepath.Join(tmp, "remote.git")
	gitCommand(t, tmp, "init", "--bare", remoteDir)

	seedDir := filepath.Join(tmp, "seed")
	RequireNoError(t, os.MkdirAll(seedDir, 0755))
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
	RequireNoError(t, os.WriteFile(filepath.Join(tmp, "repos.yaml"), []byte(registry), 0644))
	return tmp
}

func TestRunPkgList_Good(t *T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgList("table")
		RequireNoError(t, err)
	})
	AssertContains(t, out, "core-alpha")
	AssertContains(t, out, "core-beta")
	AssertContains(t, out, "core setup")
}

func TestRunPkgList_JSON(t *T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgList("json")
		RequireNoError(t, err)
	})

	var report pkgListReport
	RequireNoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &report))
	AssertEqual(t, "json", report.Format)
	AssertEqual(t, 2, report.Total)
	AssertEqual(t, 1, report.Installed)
	AssertEqual(t, 1, report.Missing)
	RequireTrue(t, len(report.Packages) == 2, Sprintf("len mismatch want=%#v", 2))
	AssertEqual(t, "core-alpha", report.Packages[0].Name)
	AssertTrue(t, report.Packages[0].Installed)
	AssertEqual(t, filepath.Join(tmp, "core-alpha"), report.Packages[0].Path)
	AssertEqual(t, "core-beta", report.Packages[1].Name)
	AssertFalse(t, report.Packages[1].Installed)
}

func TestRunPkgList_UnsupportedFormat(t *T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	err := runPkgList("yaml")
	RequireTrue(t, err != nil, "RequireError")
	AssertContains(t, err.Error(), "unsupported format")
}

func TestRunPkgOutdated_JSON(t *T) {
	tmp := setupOutdatedRegistry(t)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgOutdated("json")
		RequireNoError(t, err)
	})

	var report pkgOutdatedReport
	RequireNoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &report))
	AssertEqual(t, "json", report.Format)
	AssertEqual(t, 3, report.Total)
	AssertEqual(t, 2, report.Installed)
	AssertEqual(t, 1, report.Missing)
	AssertEqual(t, 1, report.Outdated)
	AssertEqual(t, 1, report.UpToDate)
	RequireTrue(t, len(report.Packages) == 3, Sprintf("len mismatch want=%#v", 3))

	var staleFound, freshFound, missingFound bool
	for _, pkg := range report.Packages {
		switch pkg.Name {
		case "core-stale":
			staleFound = true
			AssertTrue(t, pkg.Installed)
			AssertFalse(t, pkg.UpToDate)
			AssertEqual(t, 1, pkg.Behind)
		case "core-fresh":
			freshFound = true
			AssertTrue(t, pkg.Installed)
			AssertTrue(t, pkg.UpToDate)
			AssertEqual(t, 0, pkg.Behind)
		case "core-missing":
			missingFound = true
			AssertFalse(t, pkg.Installed)
			AssertFalse(t, pkg.UpToDate)
			AssertEqual(t, 0, pkg.Behind)
		}
	}
	AssertTrue(t, staleFound)
	AssertTrue(t, freshFound)
	AssertTrue(t, missingFound)
}

func TestRenderPkgSearchResults_ShowsMetadata(t *T) {
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
	AssertContains(t, out, "host-uk/core-alpha")
	AssertContains(t, out, "Alpha package")
	AssertContains(t, out, "42 stars")
	AssertContains(t, out, "Go")
	AssertContains(t, out, "updated 2h ago")
}

func TestRunPkgSearch_RespectsLimitWithCachedResults(t *T) {
	tmp := t.TempDir()
	writeTestRegistry(t, tmp)
	withWorkingDir(t, tmp)

	c, err := cache.New(nil, filepath.Join(tmp, ".core", "cache"), 0)
	RequireNoError(t, err)
	RequireNoError(t, c.Set(cache.GitHubReposKey("host-uk"), []ghRepo{
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
		RequireNoError(t, err)
	})
	AssertContains(t, out, "core-alpha")
	AssertNotContains(t, out, "core-beta")
}

func TestRunPkgUpdate_NoArgs_UpdatesAll(t *T) {
	tmp := setupOutdatedRegistry(t)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgUpdate(nil, false, "table")
		RequireNoError(t, err)
	})
	AssertContains(t, out, "updating")
	AssertContains(t, out, "core-fresh")
	AssertContains(t, out, "core-stale")
}

func TestRunPkgUpdate_JSON(t *T) {
	tmp := setupOutdatedRegistry(t)
	withWorkingDir(t, tmp)

	out := capturePkgOutput(t, func() {
		err := runPkgUpdate(nil, false, "json")
		RequireNoError(t, err)
	})

	var report pkgUpdateReport
	RequireNoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &report))
	AssertEqual(t, "json", report.Format)
	AssertEqual(t, 3, report.Total)
	AssertEqual(t, 2, report.Installed)
	AssertEqual(t, 1, report.Missing)
	AssertEqual(t, 1, report.Updated)
	AssertEqual(t, 1, report.UpToDate)
	AssertEqual(t, 0, report.Failed)
	RequireTrue(t, len(report.Packages) == 3, Sprintf("len mismatch want=%#v", 3))

	var updatedFound, upToDateFound, missingFound bool
	for _, pkg := range report.Packages {
		switch pkg.Name {
		case "core-stale":
			updatedFound = true
			AssertTrue(t, pkg.Installed)
			AssertEqual(t, "updated", pkg.Status)
		case "core-fresh":
			upToDateFound = true
			AssertTrue(t, pkg.Installed)
			AssertEqual(t, "up_to_date", pkg.Status)
		case "core-missing":
			missingFound = true
			AssertFalse(t, pkg.Installed)
			AssertEqual(t, "missing", pkg.Status)
		}
	}
	AssertTrue(t, updatedFound)
	AssertTrue(t, upToDateFound)
	AssertTrue(t, missingFound)
}
