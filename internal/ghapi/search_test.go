package ghapi

import (
	"strings"
	"testing"
	"time"
)

func TestBuildRepoQuery(t *testing.T) {
	tests := []struct {
		name     string
		params   RepoSearchParams
		expected []string // substrings that should be in the query
	}{
		{
			name: "single language",
			params: RepoSearchParams{
				Languages: []string{"go"},
				Limit:     10,
			},
			expected: []string{"language:go"},
		},
		{
			name: "multiple languages",
			params: RepoSearchParams{
				Languages: []string{"go", "rust"},
				Limit:     10,
			},
			expected: []string{"language:go", "language:rust"},
		},
		{
			name: "with topics",
			params: RepoSearchParams{
				Languages: []string{"python"},
				Topics:    []string{"web", "cli"},
				Limit:     10,
			},
			expected: []string{"language:python", "topic:web", "topic:cli"},
		},
		{
			name: "with min-stars",
			params: RepoSearchParams{
				Languages: []string{"javascript"},
				MinStars:  100,
				Limit:     10,
			},
			expected: []string{"language:javascript", "stars:>=100"},
		},
		{
			name: "all filters combined",
			params: RepoSearchParams{
				Languages: []string{"typescript"},
				Topics:    []string{"react"},
				MinStars:  50,
				Limit:     10,
			},
			expected: []string{"language:typescript", "topic:react", "stars:>=50"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRepoQuery(tt.params)
			for _, expected := range tt.expected {
				if !strings.Contains(got, expected) {
					t.Errorf("buildRepoQuery() = %q, missing expected substring %q", got, expected)
				}
			}
		})
	}
}

func TestMapRepo(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		full_name         string
		description       string
		language          string
		stargazers_count  int
		forks_count       int
		open_issues_count int
		pushed_at         time.Time
		topics            []string
		expect_full_name  string
		expect_desc_empty bool
		expect_lang_empty bool
	}{
		{
			name:              "normal repo with all fields",
			full_name:         "user/repo",
			description:       "A test repository",
			language:          "Go",
			stargazers_count:  100,
			forks_count:       10,
			open_issues_count: 5,
			pushed_at:         now,
			topics:            []string{"go", "cli"},
			expect_full_name:  "user/repo",
			expect_desc_empty: false,
			expect_lang_empty: false,
		},
		{
			name:              "repo with empty description",
			full_name:         "user/repo",
			description:       "",
			language:          "Rust",
			stargazers_count:  50,
			forks_count:       5,
			open_issues_count: 2,
			pushed_at:         now,
			topics:            nil,
			expect_full_name:  "user/repo",
			expect_desc_empty: true,
			expect_lang_empty: false,
		},
		{
			name:              "repo with empty language",
			full_name:         "user/repo",
			description:       "Description",
			language:          "",
			stargazers_count:  20,
			forks_count:       3,
			open_issues_count: 1,
			pushed_at:         now,
			topics:            []string{},
			expect_full_name:  "user/repo",
			expect_desc_empty: false,
			expect_lang_empty: true,
		},
		{
			name:              "repo with zero pushed_at",
			full_name:         "user/repo",
			description:       "Test",
			language:          "Python",
			stargazers_count:  75,
			forks_count:       8,
			open_issues_count: 3,
			pushed_at:         time.Time{},
			topics:            []string{"python"},
			expect_full_name:  "user/repo",
			expect_desc_empty: false,
			expect_lang_empty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := RepoResult{
				FullName:        tt.full_name,
				Description:     tt.description,
				Language:        tt.language,
				Stars:           tt.stargazers_count,
				Forks:           tt.forks_count,
				OpenIssues:      tt.open_issues_count,
				LastPushedAt:    tt.pushed_at,
				Topics:          tt.topics,
				GoodFirstIssues: 0,
				HelpWantedCount: 0,
				HasContribGuide: false,
			}

			if repo.FullName != tt.expect_full_name {
				t.Errorf("expected FullName %q, got %q", tt.expect_full_name, repo.FullName)
			}

			desc_empty := repo.Description == ""
			if desc_empty != tt.expect_desc_empty {
				t.Errorf("description empty mismatch: got %v, want %v", desc_empty, tt.expect_desc_empty)
			}

			lang_empty := repo.Language == ""
			if lang_empty != tt.expect_lang_empty {
				t.Errorf("language empty mismatch: got %v, want %v", lang_empty, tt.expect_lang_empty)
			}
		})
	}
}
