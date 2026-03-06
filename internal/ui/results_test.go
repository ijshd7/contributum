package ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ijshd7/contributum/internal/ghapi"
	"github.com/ijshd7/contributum/internal/scoring"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},            // no truncation needed
		{"hello world", 8, "hello w…"},    // truncate: s[:7]+"…"
		{"hi", 2, "hi"},                   // no truncation: len == maxLen
		{"a", 1, "a"},                     // no truncation: len <= maxLen
		{"abc", 2, "a…"},                  // truncate: s[:1]+"…"
		{"hello world", 5, "hell…"},       // truncate: s[:4]+"…"
		{"", 5, ""},                       // empty string, no truncation
		{"x", 5, "x"},                     // single char well within limit
	}

	for _, tt := range tests {
		t.Run("truncate("+tt.input+")", func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
			}
		})
	}
}

func TestFormatStars(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{1000, "1.0k"},
		{1500, "1.5k"},
		{10000, "10.0k"},
		{25100, "25.1k"},
		{1000000, "1000.0k"},
	}

	for _, tt := range tests {
		t.Run("stars", func(t *testing.T) {
			got := formatStars(tt.input)
			if got != tt.expected {
				t.Errorf("formatStars(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"zero time", time.Time{}, "unknown"},
		{"today", now.Add(-1 * time.Hour), "today"}, // 0 days
		{"yesterday", now.AddDate(0, 0, -1), "1 day"},
		{"7 days ago", now.AddDate(0, 0, -7), "7d"},
		{"30 days ago", now.AddDate(0, 0, -30), "1mo"},
		{"90 days ago", now.AddDate(0, 0, -90), "3mo"},
		{"364 days ago", now.AddDate(0, 0, -364), "12mo"}, // < 365 days
		{"366 days ago", now.AddDate(0, 0, -366), "1y"},   // >= 365 days
		{"2 years ago", now.AddDate(-2, 0, 0), "2y"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAge(tt.input)
			if got != tt.expected {
				t.Errorf("formatAge() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRenderJSON(t *testing.T) {
	type testData struct {
		Name  string
		Count int
	}

	data := []testData{
		{"test1", 10},
		{"test2", 20},
	}

	var buf bytes.Buffer
	err := renderJSON(data, &buf)

	if err != nil {
		t.Fatalf("renderJSON() error = %v", err)
	}

	output := buf.String()

	// Verify it's valid JSON
	var result []testData
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("renderJSON() output is not valid JSON: %v", err)
	}

	// Verify content
	if len(result) != len(data) {
		t.Errorf("expected %d items, got %d", len(data), len(result))
	}

	// Verify it's indented (contains newlines)
	if !strings.Contains(output, "\n") {
		t.Errorf("renderJSON() output should be indented")
	}
}

func TestRenderRepoTable(t *testing.T) {
	tests := []struct {
		name          string
		repos         []scoring.ScoredRepo
		expectTitle   string
		expectWarning bool
		expectCount   string
	}{
		{
			name:          "empty repos",
			repos:         []scoring.ScoredRepo{},
			expectTitle:   "",
			expectWarning: true,
			expectCount:   "No repositories found",
		},
		{
			name: "single repo",
			repos: []scoring.ScoredRepo{
				{
					RepoResult: ghapi.RepoResult{
						FullName:    "user/repo",
						Description: "A test repo",
						Language:    "Go",
						Stars:       100,
					},
					Score: 85.5,
				},
			},
			expectTitle:   "Repository Search Results",
			expectWarning: false,
			expectCount:   "Found 1 repositories",
		},
		{
			name: "multiple repos",
			repos: []scoring.ScoredRepo{
				{
					RepoResult: ghapi.RepoResult{
						FullName:    "user/repo1",
						Description: "Repo 1",
						Language:    "Go",
						Stars:       50,
					},
					Score: 70.0,
				},
				{
					RepoResult: ghapi.RepoResult{
						FullName:    "user/repo2",
						Description: "Repo 2",
						Language:    "Rust",
						Stars:       200,
					},
					Score: 90.0,
				},
				{
					RepoResult: ghapi.RepoResult{
						FullName:    "user/repo3",
						Description: "Repo 3",
						Language:    "Python",
						Stars:       150,
					},
					Score: 80.5,
				},
			},
			expectTitle:   "Repository Search Results",
			expectWarning: false,
			expectCount:   "Found 3 repositories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderRepoTable(tt.repos, &buf)
			output := buf.String()

			if tt.expectWarning {
				if !strings.Contains(output, "No repositories found") {
					t.Errorf("expected warning for empty repos, got: %s", output)
				}
			} else {
				if tt.expectTitle != "" && !strings.Contains(output, tt.expectTitle) {
					t.Errorf("expected title %q in output", tt.expectTitle)
				}
			}

			if !strings.Contains(output, tt.expectCount) {
				t.Errorf("expected count %q in output, got: %s", tt.expectCount, output)
			}
		})
	}
}

func TestRenderIssueTable(t *testing.T) {
	tests := []struct {
		name          string
		issues        []ghapi.IssueResult
		expectTitle   string
		expectWarning bool
		expectCount   string
	}{
		{
			name:          "empty issues",
			issues:        []ghapi.IssueResult{},
			expectTitle:   "",
			expectWarning: true,
			expectCount:   "No issues found",
		},
		{
			name: "single issue",
			issues: []ghapi.IssueResult{
				{
					RepoFullName: "user/repo",
					Title:        "Fix bug",
					Labels:       []string{"bug"},
					CreatedAt:    time.Now().AddDate(0, 0, -5),
					Comments:     3,
				},
			},
			expectTitle:   "Open Issues",
			expectWarning: false,
			expectCount:   "Found 1 issues",
		},
		{
			name: "multiple issues",
			issues: []ghapi.IssueResult{
				{
					RepoFullName: "user/repo1",
					Title:        "Fix parsing error",
					Labels:       []string{"bug", "high-priority"},
					CreatedAt:    time.Now().AddDate(0, 0, -10),
					Comments:     5,
				},
				{
					RepoFullName: "user/repo2",
					Title:        "Add documentation",
					Labels:       []string{"documentation", "good first issue"},
					CreatedAt:    time.Now().AddDate(0, 0, -20),
					Comments:     0,
				},
				{
					RepoFullName: "user/repo3",
					Title:        "Optimize performance",
					Labels:       []string{"enhancement"},
					CreatedAt:    time.Now().AddDate(0, 0, -3),
					Comments:     8,
				},
			},
			expectTitle:   "Open Issues",
			expectWarning: false,
			expectCount:   "Found 3 issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			renderIssueTable(tt.issues, &buf)
			output := buf.String()

			if tt.expectWarning {
				if !strings.Contains(output, "No issues found") {
					t.Errorf("expected warning for empty issues, got: %s", output)
				}
			} else {
				if tt.expectTitle != "" && !strings.Contains(output, tt.expectTitle) {
					t.Errorf("expected title %q in output", tt.expectTitle)
				}
			}

			if !strings.Contains(output, tt.expectCount) {
				t.Errorf("expected count %q in output, got: %s", tt.expectCount, output)
			}
		})
	}
}
