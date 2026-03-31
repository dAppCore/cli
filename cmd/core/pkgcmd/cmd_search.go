package pkgcmd

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"forge.lthn.ai/core/cli/pkg/cli"
	"forge.lthn.ai/core/go-cache"
	"forge.lthn.ai/core/go-i18n"
	coreio "forge.lthn.ai/core/go-io"
	"forge.lthn.ai/core/go-scm/repos"
	"github.com/spf13/cobra"
)

var (
	searchOrg     string
	searchPattern string
	searchType    string
	searchLimit   int
	searchRefresh bool
	searchFormat  string
)

// addPkgSearchCommand adds the 'pkg search' command.
func addPkgSearchCommand(parent *cobra.Command) {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: i18n.T("cmd.pkg.search.short"),
		Long:  i18n.T("cmd.pkg.search.long"),
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			org := searchOrg
			pattern := resolvePkgSearchPattern(searchPattern, args)
			limit := searchLimit
			if org == "" {
				org = "host-uk"
			}
			if limit == 0 {
				limit = 50
			}
			return runPkgSearch(org, pattern, searchType, limit, searchRefresh, searchFormat)
		},
	}

	searchCmd.Flags().StringVar(&searchOrg, "org", "", i18n.T("cmd.pkg.search.flag.org"))
	searchCmd.Flags().StringVar(&searchPattern, "pattern", "", i18n.T("cmd.pkg.search.flag.pattern"))
	searchCmd.Flags().StringVar(&searchType, "type", "", i18n.T("cmd.pkg.search.flag.type"))
	searchCmd.Flags().IntVar(&searchLimit, "limit", 0, i18n.T("cmd.pkg.search.flag.limit"))
	searchCmd.Flags().BoolVar(&searchRefresh, "refresh", false, i18n.T("cmd.pkg.search.flag.refresh"))
	searchCmd.Flags().StringVar(&searchFormat, "format", "table", "Output format: table or json")

	parent.AddCommand(searchCmd)
}

type ghRepo struct {
	FullName        string     `json:"fullName"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Visibility      string     `json:"visibility"`
	UpdatedAt       string     `json:"updatedAt"`
	StargazerCount  int        `json:"stargazerCount"`
	PrimaryLanguage ghLanguage `json:"primaryLanguage"`
}

type ghLanguage struct {
	Name string `json:"name"`
}

type pkgSearchEntry struct {
	Name            string `json:"name"`
	FullName        string `json:"fullName,omitempty"`
	Description     string `json:"description,omitempty"`
	Visibility      string `json:"visibility,omitempty"`
	StargazerCount  int    `json:"stargazerCount,omitempty"`
	PrimaryLanguage string `json:"primaryLanguage,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
	Updated         string `json:"updated,omitempty"`
}

type pkgSearchReport struct {
	Format  string           `json:"format"`
	Org     string           `json:"org"`
	Pattern string           `json:"pattern"`
	Type    string           `json:"type,omitempty"`
	Limit   int              `json:"limit"`
	Cached  bool             `json:"cached"`
	Count   int              `json:"count"`
	Repos   []pkgSearchEntry `json:"repos"`
}

func runPkgSearch(org, pattern, repoType string, limit int, refresh bool, format string) error {
	// Initialize cache in workspace .core/ directory
	var cacheDir string
	if regPath, err := repos.FindRegistry(coreio.Local); err == nil {
		cacheDir = filepath.Join(filepath.Dir(regPath), ".core", "cache")
	}

	c, err := cache.New(coreio.Local, cacheDir, 0)
	if err != nil {
		c = nil
	}

	cacheKey := cache.GitHubReposKey(org)
	var ghRepos []ghRepo
	var fromCache bool

	// Try cache first (unless refresh requested)
	if c != nil && !refresh {
		if found, err := c.Get(cacheKey, &ghRepos); found && err == nil {
			fromCache = true
		}
	}

	// Fetch from GitHub if not cached
	if !fromCache {
		if !ghAuthenticated() {
			return errors.New(i18n.T("cmd.pkg.error.gh_not_authenticated"))
		}

		if os.Getenv("GH_TOKEN") != "" && !strings.EqualFold(format, "json") {
			fmt.Printf("%s %s\n", dimStyle.Render(i18n.Label("note")), i18n.T("cmd.pkg.search.gh_token_warning"))
			fmt.Printf("%s %s\n\n", dimStyle.Render(""), i18n.T("cmd.pkg.search.gh_token_unset"))
		}

		if !strings.EqualFold(format, "json") {
			fmt.Printf("%s %s... ", dimStyle.Render(i18n.T("cmd.pkg.search.fetching_label")), org)
		}

		cmd := exec.Command("gh", "repo", "list", org,
			"--json", "fullName,name,description,visibility,updatedAt,stargazerCount,primaryLanguage",
			"--limit", fmt.Sprintf("%d", limit))
		output, err := cmd.CombinedOutput()

		if err != nil {
			if !strings.EqualFold(format, "json") {
				fmt.Println()
			}
			errStr := strings.TrimSpace(string(output))
			if strings.Contains(errStr, "401") || strings.Contains(errStr, "Bad credentials") {
				return errors.New(i18n.T("cmd.pkg.error.auth_failed"))
			}
			return fmt.Errorf("%s: %s", i18n.T("cmd.pkg.error.search_failed"), errStr)
		}

		if err := json.Unmarshal(output, &ghRepos); err != nil {
			return fmt.Errorf("%s: %w", i18n.T("i18n.fail.parse", "results"), err)
		}

		if c != nil {
			_ = c.Set(cacheKey, ghRepos)
		}

		if !strings.EqualFold(format, "json") {
			fmt.Printf("%s\n", successStyle.Render("✓"))
		}
	}

	// Filter by glob pattern and type
	var filtered []ghRepo
	for _, r := range ghRepos {
		if !matchGlob(pattern, r.Name) {
			continue
		}
		if repoType != "" && !strings.Contains(r.Name, repoType) {
			continue
		}
		filtered = append(filtered, r)
	}

	if len(filtered) == 0 {
		if strings.EqualFold(format, "json") {
			report := buildPkgSearchReport(org, pattern, repoType, limit, fromCache, filtered)
			return printPkgSearchJSON(report)
		}
		fmt.Println(i18n.T("cmd.pkg.search.no_repos_found"))
		return nil
	}

	slices.SortFunc(filtered, func(a, b ghRepo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	if strings.EqualFold(format, "json") {
		report := buildPkgSearchReport(org, pattern, repoType, limit, fromCache, filtered)
		return printPkgSearchJSON(report)
	}

	if fromCache && !strings.EqualFold(format, "json") {
		age := c.Age(cacheKey)
		fmt.Printf("%s %s %s\n", dimStyle.Render(i18n.T("cmd.pkg.search.cache_label")), org, dimStyle.Render(fmt.Sprintf("(%s ago)", age.Round(time.Second))))
	}
	renderPkgSearchResults(filtered)

	fmt.Println()
	fmt.Printf("%s %s\n", i18n.T("common.hint.install_with"), dimStyle.Render(fmt.Sprintf("core pkg install %s/<repo-name>", org)))

	return nil
}

func renderPkgSearchResults(repos []ghRepo) {
	fmt.Print(i18n.T("cmd.pkg.search.found_repos", map[string]int{"Count": len(repos)}) + "\n\n")

	for _, r := range repos {
		visibility := ""
		if r.Visibility == "private" {
			visibility = dimStyle.Render(" " + i18n.T("cmd.pkg.search.private_label"))
		}

		desc := r.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		if desc == "" {
			desc = dimStyle.Render(i18n.T("cmd.pkg.no_description"))
		}

		fmt.Printf("  %s%s\n", repoNameStyle.Render(r.Name), visibility)
		fmt.Printf("    %s\n", desc)

		if meta := formatPkgSearchMetadata(r); meta != "" {
			fmt.Printf("    %s\n", dimStyle.Render(meta))
		}
	}
}

func buildPkgSearchReport(org, pattern, repoType string, limit int, cached bool, repos []ghRepo) pkgSearchReport {
	report := pkgSearchReport{
		Format:  "json",
		Org:     org,
		Pattern: pattern,
		Type:    repoType,
		Limit:   limit,
		Cached:  cached,
		Count:   len(repos),
		Repos:   make([]pkgSearchEntry, 0, len(repos)),
	}

	for _, r := range repos {
		report.Repos = append(report.Repos, pkgSearchEntry{
			Name:            r.Name,
			FullName:        r.FullName,
			Description:     r.Description,
			Visibility:      r.Visibility,
			StargazerCount:  r.StargazerCount,
			PrimaryLanguage: strings.TrimSpace(r.PrimaryLanguage.Name),
			UpdatedAt:       r.UpdatedAt,
			Updated:         formatPkgSearchUpdatedAt(r.UpdatedAt),
		})
	}

	return report
}

func printPkgSearchJSON(report pkgSearchReport) error {
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("i18n.fail.format", "search results"), err)
	}

	fmt.Println(string(out))
	return nil
}

func formatPkgSearchMetadata(r ghRepo) string {
	var parts []string

	if r.StargazerCount > 0 {
		parts = append(parts, fmt.Sprintf("%d stars", r.StargazerCount))
	}

	if lang := strings.TrimSpace(r.PrimaryLanguage.Name); lang != "" {
		parts = append(parts, lang)
	}

	if updated := formatPkgSearchUpdatedAt(r.UpdatedAt); updated != "" {
		parts = append(parts, "updated "+updated)
	}

	return strings.Join(parts, "  ")
}

func formatPkgSearchUpdatedAt(raw string) string {
	if raw == "" {
		return ""
	}

	updatedAt, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}

	return cli.FormatAge(updatedAt)
}

func resolvePkgSearchPattern(flagPattern string, args []string) string {
	if flagPattern != "" {
		return flagPattern
	}
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		return args[0]
	}
	return "*"
}

// matchGlob does simple glob matching with * wildcards
func matchGlob(pattern, name string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}

	parts := strings.Split(pattern, "*")
	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(name[pos:], part)
		if idx == -1 {
			return false
		}
		if i == 0 && !strings.HasPrefix(pattern, "*") && idx != 0 {
			return false
		}
		pos += idx + len(part)
	}
	if !strings.HasSuffix(pattern, "*") && pos != len(name) {
		return false
	}
	return true
}
