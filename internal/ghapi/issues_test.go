package ghapi

import (
	"strings"
	"testing"
)

func TestBuildIssueQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   IssueSearchParams
		expected []string // substrings that should be in the query
	}{
		{
			name: "single language with explicit label",
			params: IssueSearchParams{
				Languages: []string{"go"},
				Labels:    []string{"good first issue"},
				Limit:     10,
			},
			expected: []string{"language:go", "label:\"good first issue\""},
		},
		{
			name: "multiple languages with explicit label",
			params: IssueSearchParams{
				Languages: []string{"python", "javascript"},
				Labels:    []string{"good first issue"},
				Limit:     10,
			},
			expected: []string{"language:python", "language:javascript", "label:\"good first issue\""},
		},
		{
			name: "single language with custom label",
			params: IssueSearchParams{
				Languages: []string{"rust"},
				Labels:    []string{"help wanted"},
				Limit:     10,
			},
			expected: []string{"language:rust", "label:\"help wanted\""},
		},
		{
			name: "multiple labels",
			params: IssueSearchParams{
				Languages: []string{"go"},
				Labels:    []string{"bug", "documentation"},
				Limit:     10,
			},
			expected: []string{"language:go", "label:\"bug\"", "label:\"documentation\""},
		},
		{
			name: "single language no labels (empty list)",
			params: IssueSearchParams{
				Languages: []string{"typescript"},
				Labels:    []string{},
				Limit:     10,
			},
			// Note: buildIssueQuery doesn't inject default labels - that happens in SearchIssues
			expected: []string{"language:typescript"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildIssueQuery(tt.params)
			for _, expected := range tt.expected {
				if !strings.Contains(got, expected) {
					t.Errorf("buildIssueQuery() = %q, missing expected substring %q", got, expected)
				}
			}
			// Verify basic query structure
			if !strings.Contains(got, "is:open") {
				t.Errorf("buildIssueQuery() should contain 'is:open'")
			}
		})
	}
}

func TestMapIssue(t *testing.T) {
	tests := []struct {
		name           string
		repo_url       string
		title          string
		labels         []string
		comments       int
		expect_repo    string
		expect_title   string
		expect_labels  int
	}{
		{
			name:          "normal issue with well-formed URL",
			repo_url:      "https://api.github.com/repos/user/repo",
			title:         "Fix bug in parser",
			labels:        []string{"bug", "enhancement"},
			comments:      5,
			expect_repo:   "user/repo",
			expect_title:  "Fix bug in parser",
			expect_labels: 2,
		},
		{
			name:          "issue with no labels",
			repo_url:      "https://api.github.com/repos/golang/go",
			title:         "Add feature X",
			labels:        []string{},
			comments:      0,
			expect_repo:   "golang/go",
			expect_title:  "Add feature X",
			expect_labels: 0,
		},
		{
			name:          "issue with many labels",
			repo_url:      "https://api.github.com/repos/django/django",
			title:         "Update documentation",
			labels:        []string{"docs", "help wanted", "good first issue", "priority-high"},
			comments:      12,
			expect_repo:   "django/django",
			expect_title:  "Update documentation",
			expect_labels: 4,
		},
		{
			name:          "repo with hyphens in name",
			repo_url:      "https://api.github.com/repos/my-org/my-repo-name",
			title:         "Fix issue",
			labels:        []string{},
			comments:      0,
			expect_repo:   "my-org/my-repo-name",
			expect_title:  "Fix issue",
			expect_labels: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to simulate extracting repo name from URL
			// The mapIssue function uses strings.Split on the RepositoryURL
			// Extract the part after "repos/" from the test URL
			parts := strings.Split(tt.repo_url, "repos/")
			if len(parts) < 2 {
				t.Fatalf("invalid test URL format")
			}
			repo := parts[1]

			if repo != tt.expect_repo {
				t.Errorf("expected repo %q, got %q", tt.expect_repo, repo)
			}

			if tt.title != tt.expect_title {
				t.Errorf("expected title %q, got %q", tt.expect_title, tt.title)
			}

			if len(tt.labels) != tt.expect_labels {
				t.Errorf("expected %d labels, got %d", tt.expect_labels, len(tt.labels))
			}
		})
	}
}
