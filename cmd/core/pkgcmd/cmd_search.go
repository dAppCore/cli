package pkgcmd

import (
	"cmp"
	"context"
	"slices"
	"time"

	"dappco.re/go"
	"dappco.re/go/cache"
	"dappco.re/go/cli/pkg/cli"
	"dappco.re/go/cli/pkg/i18n"
	coreio "dappco.re/go/io"
	"dappco.re/go/scm/repos"
)

func pkgSearchAction(opts core.Options) core.Result {
	org := opts.String("org")
	pattern := opts.String("pattern")
	repoType := opts.String("type")
	limit := opts.Int("limit")
	refresh := opts.Bool("refresh")

	if org == "" {
		org = "host-uk"
	}
	if pattern == "" {
		pattern = "*"
	}
	if limit == 0 {
		limit = 50
	}

	if err := runPkgSearch(org, pattern, repoType, limit, refresh); err != nil {
		return core.Result{Value: err, OK: false}
	}
	return core.Result{OK: true}
}

type ghRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	UpdatedAt   string `json:"updated_at"`
	Language    string `json:"language"`
}

func runPkgSearch(org, pattern, repoType string, limit int, refresh bool) error {
	// Initialise cache in workspace .core/ directory.
	var cacheDirectory string
	if registryPath, err := repos.FindRegistry(coreio.Local); err == nil {
		cacheDirectory = core.Path(core.PathDir(registryPath), ".core", "cache")
	}

	cacheInstance, err := cache.New(coreio.Local, cacheDirectory, 0)
	if err != nil {
		cacheInstance = nil
	}

	cacheKey := cache.GitHubReposKey(org)
	var ghRepos []ghRepo
	var fromCache bool

	// Try cache first (unless refresh requested).
	if cacheInstance != nil && !refresh {
		if found, err := cacheInstance.Get(cacheKey, &ghRepos); found && err == nil {
			fromCache = true
			age := cacheInstance.Age(cacheKey)
			cli.Println("%s %s %s", dimStyle.Render(i18n.T("cmd.pkg.search.cache_label")), org, dimStyle.Render(cli.Sprintf("(%s ago)", age.Round(time.Second))))
		}
	}

	// Fetch from GitHub if not cached.
	if !fromCache {
		if !ghAuthenticated() {
			return cli.Err(i18n.T("cmd.pkg.error.gh_not_authenticated"))
		}

		if core.Env("GH_TOKEN") != "" {
			cli.Println("%s %s", dimStyle.Render(i18n.Label("note")), i18n.T("cmd.pkg.search.gh_token_warning"))
			cli.Println("%s %s\n", dimStyle.Render(""), i18n.T("cmd.pkg.search.gh_token_unset"))
		}

		cli.Print("%s %s... ", dimStyle.Render(i18n.T("cmd.pkg.search.fetching_label")), org)

		result := cli.Core().Process().Run(context.Background(), "gh",
			"repo", "list", org,
			"--json", "name,description,visibility,updatedAt,primaryLanguage",
			"--limit", cli.Sprintf("%d", limit))
		output, _ := result.Value.(string)

		if !result.OK {
			cli.Blank()
			errorOutput := core.Trim(output)
			if errorOutput == "" {
				if err, ok := result.Value.(error); ok {
					errorOutput = core.Trim(err.Error())
				}
			}
			if core.Contains(errorOutput, "401") || core.Contains(errorOutput, "Bad credentials") {
				return cli.Err(i18n.T("cmd.pkg.error.auth_failed"))
			}
			return cli.Err("%s: %s", i18n.T("cmd.pkg.error.search_failed"), errorOutput)
		}

		parseResult := core.JSONUnmarshal([]byte(output), &ghRepos)
		if !parseResult.OK {
			return cli.Wrap(parseResult.Value.(error), i18n.T("i18n.fail.parse", "results"))
		}

		if cacheInstance != nil {
			_ = cacheInstance.Set(cacheKey, ghRepos)
		}

		cli.Println("%s", successStyle.Render("ok"))
	}

	// Filter by glob pattern and type.
	var filtered []ghRepo
	for _, repo := range ghRepos {
		if !matchGlob(pattern, repo.Name) {
			continue
		}
		if repoType != "" && !core.Contains(repo.Name, repoType) {
			continue
		}
		filtered = append(filtered, repo)
	}

	if len(filtered) == 0 {
		cli.Println("%s", i18n.T("cmd.pkg.search.no_repos_found"))
		return nil
	}

	slices.SortFunc(filtered, func(a, b ghRepo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	cli.Print(i18n.T("cmd.pkg.search.found_repos", map[string]int{"Count": len(filtered)}) + "\n\n")

	for _, repo := range filtered {
		visibility := ""
		if repo.Visibility == "private" {
			visibility = dimStyle.Render(" " + i18n.T("cmd.pkg.search.private_label"))
		}

		description := repo.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		if description == "" {
			description = dimStyle.Render(i18n.T("cmd.pkg.no_description"))
		}

		cli.Println("  %s%s", repoNameStyle.Render(repo.Name), visibility)
		cli.Println("    %s", description)
	}

	cli.Blank()
	cli.Println("%s %s", i18n.T("common.hint.install_with"), dimStyle.Render(cli.Sprintf("core pkg install %s/<repo-name>", org)))

	return nil
}

// matchGlob does simple glob matching with * wildcards.
//
//	matchGlob("core-*", "core-php")   // true
//	matchGlob("*-mod", "core-php")    // false
func matchGlob(pattern, name string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}

	parts := core.Split(pattern, "*")
	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		// Find part in name starting from pos.
		remaining := name[pos:]
		idx := -1
		for j := 0; j <= len(remaining)-len(part); j++ {
			if remaining[j:j+len(part)] == part {
				idx = j
				break
			}
		}
		if idx == -1 {
			return false
		}
		if i == 0 && !core.HasPrefix(pattern, "*") && idx != 0 {
			return false
		}
		pos += idx + len(part)
	}
	if !core.HasSuffix(pattern, "*") && pos != len(name) {
		return false
	}
	return true
}
