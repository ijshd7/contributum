package scoring

import (
	"testing"
	"time"

	"github.com/ijshd7/contributum/internal/ghapi"
)

func TestActivityScore(t *testing.T) {
	tests := []struct {
		name     string
		repo     ghapi.RepoResult
		wantMin  float64
		wantMax  float64
	}{
		{
			name:    "zero stars and forks, recent push",
			repo:    ghapi.RepoResult{Stars: 0, Forks: 0, LastPushedAt: time.Now()},
			wantMin: 45, // recency dominates
			wantMax: 55,
		},
		{
			name:    "high stars, old push",
			repo:    ghapi.RepoResult{Stars: 50000, Forks: 5000, LastPushedAt: time.Now().Add(-400 * 24 * time.Hour)},
			wantMin: 40,
			wantMax: 55,
		},
		{
			name:    "moderate repo, recent push",
			repo:    ghapi.RepoResult{Stars: 1000, Forks: 200, LastPushedAt: time.Now().Add(-3 * 24 * time.Hour)},
			wantMin: 60,
			wantMax: 85,
		},
		{
			name:    "zero everything",
			repo:    ghapi.RepoResult{},
			wantMin: 0,
			wantMax: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := activityScore(tt.repo)
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("activityScore() = %.1f, want between %.1f and %.1f", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestFriendlinessScore(t *testing.T) {
	tests := []struct {
		name    string
		repo    ghapi.RepoResult
		wantMin float64
		wantMax float64
	}{
		{
			name:    "no friendliness signals",
			repo:    ghapi.RepoResult{},
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "max friendliness",
			repo:    ghapi.RepoResult{GoodFirstIssues: 10, HelpWantedCount: 10, HasContribGuide: true},
			wantMin: 100,
			wantMax: 100,
		},
		{
			name:    "some good first issues only",
			repo:    ghapi.RepoResult{GoodFirstIssues: 3},
			wantMin: 20,
			wantMax: 30,
		},
		{
			name:    "contributing guide only",
			repo:    ghapi.RepoResult{HasContribGuide: true},
			wantMin: 28,
			wantMax: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := friendlinessScore(tt.repo)
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("friendlinessScore() = %.1f, want between %.1f and %.1f", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestRelevanceScore(t *testing.T) {
	tests := []struct {
		name      string
		repo      ghapi.RepoResult
		languages []string
		topics    []string
		wantMin   float64
		wantMax   float64
	}{
		{
			name:      "perfect match",
			repo:      ghapi.RepoResult{Language: "Go", Topics: []string{"cli", "tools"}},
			languages: []string{"go"},
			topics:    []string{"cli", "tools"},
			wantMin:   100,
			wantMax:   100,
		},
		{
			name:      "language match, no topic match",
			repo:      ghapi.RepoResult{Language: "Go", Topics: []string{"web"}},
			languages: []string{"go"},
			topics:    []string{"cli"},
			wantMin:   58,
			wantMax:   62,
		},
		{
			name:      "no match at all",
			repo:      ghapi.RepoResult{Language: "Python", Topics: []string{"web"}},
			languages: []string{"go"},
			topics:    []string{"cli"},
			wantMin:   0,
			wantMax:   0,
		},
		{
			name:      "no topics specified",
			repo:      ghapi.RepoResult{Language: "Go"},
			languages: []string{"go"},
			topics:    nil,
			wantMin:   78,
			wantMax:   82,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := relevanceScore(tt.repo, tt.languages, tt.topics)
			if score < tt.wantMin || score > tt.wantMax {
				t.Errorf("relevanceScore() = %.1f, want between %.1f and %.1f", score, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestScoreRepos(t *testing.T) {
	repos := []ghapi.RepoResult{
		{
			FullName:        "low/score",
			Language:        "Python",
			Stars:           10,
			LastPushedAt:    time.Now().Add(-500 * 24 * time.Hour),
			GoodFirstIssues: 0,
		},
		{
			FullName:        "high/score",
			Language:        "Go",
			Stars:           5000,
			Forks:           500,
			LastPushedAt:    time.Now().Add(-2 * 24 * time.Hour),
			GoodFirstIssues: 8,
			HelpWantedCount: 5,
			HasContribGuide: true,
			Topics:          []string{"cli"},
		},
	}

	scored := ScoreRepos(repos, []string{"go"}, []string{"cli"}, Intermediate)

	if len(scored) != 2 {
		t.Fatalf("expected 2 scored repos, got %d", len(scored))
	}
	if scored[0].FullName != "high/score" {
		t.Error("expected high/score to be ranked first")
	}
	if scored[0].Score <= scored[1].Score {
		t.Error("expected first repo to have higher score than second")
	}
}

func TestScoreReposBeginnerFilters(t *testing.T) {
	repos := []ghapi.RepoResult{
		{FullName: "no/gfi", Language: "Go", GoodFirstIssues: 0},
		{FullName: "has/gfi", Language: "Go", GoodFirstIssues: 3},
	}

	scored := ScoreRepos(repos, []string{"go"}, nil, Beginner)

	if len(scored) != 1 {
		t.Fatalf("expected 1 scored repo for beginner, got %d", len(scored))
	}
	if scored[0].FullName != "has/gfi" {
		t.Error("expected only repo with good-first-issues for beginner")
	}
}

func TestParseSkillLevel(t *testing.T) {
	tests := []struct {
		input string
		want  SkillLevel
		ok    bool
	}{
		{"beginner", Beginner, true},
		{"Intermediate", Intermediate, true},
		{"ADVANCED", Advanced, true},
		{"expert", "", false},
	}

	for _, tt := range tests {
		got, ok := ParseSkillLevel(tt.input)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseSkillLevel(%q) = (%q, %v), want (%q, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}

func TestScoreReposEmpty(t *testing.T) {
	scored := ScoreRepos([]ghapi.RepoResult{}, []string{"go"}, nil, Intermediate)

	if len(scored) != 0 {
		t.Errorf("expected empty result for empty input, got %d repos", len(scored))
	}
}

func TestScoreReposAdvanced(t *testing.T) {
	repos := []ghapi.RepoResult{
		{
			FullName:        "inactive/repo",
			Language:        "Go",
			Stars:           100,
			LastPushedAt:    time.Now().Add(-500 * 24 * time.Hour), // inactive
			GoodFirstIssues: 10,                                     // high friendliness
		},
		{
			FullName:        "active/repo",
			Language:        "Go",
			Stars:           5000,
			LastPushedAt:    time.Now().Add(-2 * 24 * time.Hour), // recently active
			GoodFirstIssues: 0,                                    // no good first issues
		},
	}

	scored := ScoreRepos(repos, []string{"go"}, nil, Advanced)

	if len(scored) != 2 {
		t.Fatalf("expected 2 scored repos, got %d", len(scored))
	}

	// Advanced skill should prioritize activity (55%), so active/repo should rank higher
	// despite having lower friendliness
	if scored[0].FullName != "active/repo" {
		t.Errorf("expected active/repo to rank first for Advanced skill, got %s", scored[0].FullName)
	}
}

func TestRecencyScore(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		lastPush time.Time
		expected float64
	}{
		{
			name:     "pushed today (0 hours ago)",
			lastPush: now.Add(-1 * time.Hour), // definitely <= 7 days
			expected: 100,
		},
		{
			name:     "pushed 3 days ago",
			lastPush: now.Add(-3 * 24 * time.Hour), // definitely <= 7 days
			expected: 100,
		},
		{
			name:     "pushed 8 days ago",
			lastPush: now.Add(-8 * 24 * time.Hour), // > 7 days but <= 30 days
			expected: 75,
		},
		{
			name:     "pushed 20 days ago",
			lastPush: now.Add(-20 * 24 * time.Hour), // <= 30 days
			expected: 75,
		},
		{
			name:     "pushed 31 days ago",
			lastPush: now.Add(-31 * 24 * time.Hour), // > 30 days but <= 90 days
			expected: 50,
		},
		{
			name:     "pushed 60 days ago",
			lastPush: now.Add(-60 * 24 * time.Hour), // <= 90 days
			expected: 50,
		},
		{
			name:     "pushed 91 days ago",
			lastPush: now.Add(-91 * 24 * time.Hour), // > 90 days but <= 365 days
			expected: 25,
		},
		{
			name:     "pushed 180 days ago",
			lastPush: now.Add(-180 * 24 * time.Hour), // <= 365 days
			expected: 25,
		},
		{
			name:     "pushed 400 days ago",
			lastPush: now.Add(-400 * 24 * time.Hour), // > 365 days
			expected: 0,
		},
		{
			name:     "zero time (not set)",
			lastPush: time.Time{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := recencyScore(tt.lastPush)
			if got != tt.expected {
				t.Errorf("recencyScore() = %.0f, want %.0f", got, tt.expected)
			}
		})
	}
}

func TestWeights(t *testing.T) {
	tests := []struct {
		name          string
		skill         SkillLevel
		expActivity   float64
		expFriendly   float64
		expRelevance  float64
	}{
		{
			name:         "beginner",
			skill:        Beginner,
			expActivity:  0.25,
			expFriendly:  0.50,
			expRelevance: 0.25,
		},
		{
			name:         "intermediate",
			skill:        Intermediate,
			expActivity:  0.40,
			expFriendly:  0.35,
			expRelevance: 0.25,
		},
		{
			name:         "advanced",
			skill:        Advanced,
			expActivity:  0.55,
			expFriendly:  0.20,
			expRelevance: 0.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, friendly, relevance := weights(tt.skill)

			if act != tt.expActivity {
				t.Errorf("activity weight = %.2f, want %.2f", act, tt.expActivity)
			}
			if friendly != tt.expFriendly {
				t.Errorf("friendliness weight = %.2f, want %.2f", friendly, tt.expFriendly)
			}
			if relevance != tt.expRelevance {
				t.Errorf("relevance weight = %.2f, want %.2f", relevance, tt.expRelevance)
			}
		})
	}
}
