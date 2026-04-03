package pkgcmd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePkgSearchPattern_Good(t *testing.T) {
	t.Run("uses flag pattern when set", func(t *testing.T) {
		got := resolvePkgSearchPattern("core-*", []string{"api"})
		assert.Equal(t, "core-*", got)
	})

	t.Run("uses positional pattern when flag is empty", func(t *testing.T) {
		got := resolvePkgSearchPattern("", []string{"api"})
		assert.Equal(t, "api", got)
	})

	t.Run("defaults to wildcard when nothing is provided", func(t *testing.T) {
		got := resolvePkgSearchPattern("", nil)
		assert.Equal(t, "*", got)
	})
}

func TestBuildPkgSearchReport_Good(t *testing.T) {
	repos := []ghRepo{
		{
			FullName:       "host-uk/core-api",
			Name:           "core-api",
			Description:    "REST API framework",
			Visibility:     "public",
			UpdatedAt:      "2026-03-30T12:00:00Z",
			StargazerCount: 42,
			PrimaryLanguage: ghLanguage{
				Name: "Go",
			},
		},
	}

	report := buildPkgSearchReport("host-uk", "core-*", "api", 50, true, repos)

	assert.Equal(t, "json", report.Format)
	assert.Equal(t, "host-uk", report.Org)
	assert.Equal(t, "core-*", report.Pattern)
	assert.Equal(t, "api", report.Type)
	assert.Equal(t, 50, report.Limit)
	assert.True(t, report.Cached)
	assert.Equal(t, 1, report.Count)
	requireRepo := report.Repos
	if assert.Len(t, requireRepo, 1) {
		assert.Equal(t, "core-api", requireRepo[0].Name)
		assert.Equal(t, "host-uk/core-api", requireRepo[0].FullName)
		assert.Equal(t, "REST API framework", requireRepo[0].Description)
		assert.Equal(t, "public", requireRepo[0].Visibility)
		assert.Equal(t, 42, requireRepo[0].StargazerCount)
		assert.Equal(t, "Go", requireRepo[0].PrimaryLanguage)
		assert.Equal(t, "2026-03-30T12:00:00Z", requireRepo[0].UpdatedAt)
		assert.NotEmpty(t, requireRepo[0].Updated)
	}

	out, err := json.Marshal(report)
	assert.NoError(t, err)
	assert.Contains(t, string(out), `"format":"json"`)
}
