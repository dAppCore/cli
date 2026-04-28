package pkgcmd

import (
	. "dappco.re/go"
	"encoding/json"
)

func TestResolvePkgSearchPattern_Good(t *T) {
	t.Run("uses flag pattern when set", func(t *T) {
		got := resolvePkgSearchPattern("core-*", []string{"api"})
		AssertEqual(t, "core-*", got)
	})

	t.Run("uses positional pattern when flag is empty", func(t *T) {
		got := resolvePkgSearchPattern("", []string{"api"})
		AssertEqual(t, "api", got)
	})

	t.Run("defaults to wildcard when nothing is provided", func(t *T) {
		got := resolvePkgSearchPattern("", nil)
		AssertEqual(t, "*", got)
	})
}

func TestBuildPkgSearchReport_Good(t *T) {
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
	AssertEqual(t, "json", report.Format)
	AssertEqual(t, "host-uk", report.Org)
	AssertEqual(t, "core-*", report.Pattern)
	AssertEqual(t, "api", report.Type)
	AssertEqual(t, 50, report.Limit)
	AssertTrue(t, report.Cached)
	AssertEqual(t, 1, report.Count)
	requireRepo := report.Repos
	AssertLen(t, requireRepo, 1)
	if len(requireRepo) == 1 {
		AssertEqual(t, "core-api", requireRepo[0].Name)
		AssertEqual(t, "host-uk/core-api", requireRepo[0].FullName)
		AssertEqual(t, "REST API framework", requireRepo[0].Description)
		AssertEqual(t, "public", requireRepo[0].Visibility)
		AssertEqual(t, 42, requireRepo[0].StargazerCount)
		AssertEqual(t, "Go", requireRepo[0].PrimaryLanguage)
		AssertEqual(t, "2026-03-30T12:00:00Z", requireRepo[0].UpdatedAt)
		AssertNotEmpty(t, requireRepo[0].Updated)
	}

	out, err := json.Marshal(report)
	AssertNoError(t, err)
	AssertContains(t, string(out), `"format":"json"`)
}
