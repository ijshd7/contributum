package ghapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
	"golang.org/x/sync/errgroup"
)

type RepoSearchParams struct {
	Languages []string
	Topics    []string
	MinStars  int
	Limit     int
}

type RepoResult struct {
	FullName        string
	Description     string
	URL             string
	Language        string
	Stars           int
	Forks           int
	OpenIssues      int
	Topics          []string
	LastPushedAt    time.Time
	GoodFirstIssues int
	HelpWantedCount int
	HasContribGuide bool
}

func (c *Client) SearchRepos(ctx context.Context, params RepoSearchParams) ([]RepoResult, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	query := buildRepoQuery(params)
	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: params.Limit,
		},
	}

	result, _, err := c.gh.Search.Repositories(ctx, query, opts)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GitHub API rate limit exceeded, resets at %v. Set GITHUB_TOKEN in .env for higher limits",
				rateLimitErr.Rate.Reset.Time.Format(time.Kitchen))
		}
		return nil, fmt.Errorf("searching repositories: %w", err)
	}

	repos := make([]RepoResult, 0, len(result.Repositories))
	for _, r := range result.Repositories {
		repos = append(repos, mapRepo(r))
	}

	if err := c.enrichRepos(ctx, repos); err != nil {
		// Non-fatal: we still have basic results
		fmt.Printf("Warning: could not fetch all repo details: %v\n", err)
	}

	return repos, nil
}

func buildRepoQuery(params RepoSearchParams) string {
	var parts []string
	for _, lang := range params.Languages {
		parts = append(parts, "language:"+lang)
	}
	for _, topic := range params.Topics {
		parts = append(parts, "topic:"+topic)
	}
	if params.MinStars > 0 {
		parts = append(parts, fmt.Sprintf("stars:>=%d", params.MinStars))
	}
	return strings.Join(parts, " ")
}

func mapRepo(r *github.Repository) RepoResult {
	result := RepoResult{
		FullName: r.GetFullName(),
		URL:      r.GetHTMLURL(),
		Language: r.GetLanguage(),
		Stars:    r.GetStargazersCount(),
		Forks:    r.GetForksCount(),
		Topics:   r.Topics,
	}
	if r.Description != nil {
		result.Description = *r.Description
	}
	if r.OpenIssuesCount != nil {
		result.OpenIssues = *r.OpenIssuesCount
	}
	if r.PushedAt != nil {
		result.LastPushedAt = r.PushedAt.Time
	}
	return result
}

func (c *Client) enrichRepos(ctx context.Context, repos []RepoResult) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for i := range repos {
		i := i
		g.Go(func() error {
			return c.enrichRepo(ctx, &repos[i])
		})
	}

	return g.Wait()
}

func (c *Client) enrichRepo(ctx context.Context, repo *RepoResult) error {
	parts := strings.SplitN(repo.FullName, "/", 2)
	if len(parts) != 2 {
		return nil
	}
	owner, name := parts[0], parts[1]

	// Count good first issues
	gfiQuery := fmt.Sprintf("repo:%s is:issue is:open label:\"good first issue\"", repo.FullName)
	gfiResult, _, err := c.gh.Search.Issues(ctx, gfiQuery, &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}})
	if err == nil && gfiResult.Total != nil {
		repo.GoodFirstIssues = *gfiResult.Total
	}

	// Count help wanted issues
	hwQuery := fmt.Sprintf("repo:%s is:issue is:open label:\"help wanted\"", repo.FullName)
	hwResult, _, err := c.gh.Search.Issues(ctx, hwQuery, &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}})
	if err == nil && hwResult.Total != nil {
		repo.HelpWantedCount = *hwResult.Total
	}

	// Check for CONTRIBUTING.md
	_, _, resp, err := c.gh.Repositories.GetContents(ctx, owner, name, "CONTRIBUTING.md", nil)
	if err == nil {
		repo.HasContribGuide = true
	} else if resp != nil && resp.StatusCode == 404 {
		repo.HasContribGuide = false
	}

	return nil
}
