package ghapi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
)

type IssueSearchParams struct {
	Languages []string
	Labels    []string
	Limit     int
}

type IssueResult struct {
	RepoFullName string
	Title        string
	URL          string
	Labels       []string
	CreatedAt    time.Time
	Comments     int
}

func (c *Client) SearchIssues(ctx context.Context, params IssueSearchParams) ([]IssueResult, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if len(params.Labels) == 0 {
		params.Labels = []string{"good first issue"}
	}

	query := buildIssueQuery(params)
	opts := &github.SearchOptions{
		Sort:  "created",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: params.Limit,
		},
	}

	result, _, err := c.gh.Search.Issues(ctx, query, opts)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GitHub API rate limit exceeded, resets at %v. Set GITHUB_TOKEN in .env for higher limits",
				rateLimitErr.Rate.Reset.Time.Format(time.Kitchen))
		}
		return nil, fmt.Errorf("searching issues: %w", err)
	}

	issues := make([]IssueResult, 0, len(result.Issues))
	for _, issue := range result.Issues {
		issues = append(issues, mapIssue(issue))
	}

	return issues, nil
}

func buildIssueQuery(params IssueSearchParams) string {
	parts := []string{"is:issue", "is:open", "state:open"}
	for _, lang := range params.Languages {
		parts = append(parts, "language:"+lang)
	}
	for _, label := range params.Labels {
		parts = append(parts, fmt.Sprintf("label:\"%s\"", label))
	}
	return strings.Join(parts, " ")
}

func mapIssue(issue *github.Issue) IssueResult {
	result := IssueResult{
		Title:    issue.GetTitle(),
		URL:      issue.GetHTMLURL(),
		Comments: issue.GetComments(),
	}
	if issue.CreatedAt != nil {
		result.CreatedAt = issue.CreatedAt.Time
	}
	for _, label := range issue.Labels {
		result.Labels = append(result.Labels, label.GetName())
	}
	// Extract repo name from the URL
	if issue.RepositoryURL != nil {
		parts := strings.Split(*issue.RepositoryURL, "/")
		if len(parts) >= 2 {
			result.RepoFullName = parts[len(parts)-2] + "/" + parts[len(parts)-1]
		}
	}
	return result
}
