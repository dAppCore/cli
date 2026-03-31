package pkgcmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestRenderPkgSearchResults_ShowsMetadata(t *testing.T) {
	out := capturePkgOutput(t, func() {
		renderPkgSearchResults([]ghRepo{
			{
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

	assert.Contains(t, out, "core-alpha")
	assert.Contains(t, out, "Alpha package")
	assert.Contains(t, out, "42 stars")
	assert.Contains(t, out, "Go")
	assert.Contains(t, out, "updated 2h ago")
}
