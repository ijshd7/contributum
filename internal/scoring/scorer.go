package scoring

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ijshd7/contributum/internal/ghapi"
)

type SkillLevel string

const (
	Beginner     SkillLevel = "beginner"
	Intermediate SkillLevel = "intermediate"
	Advanced     SkillLevel = "advanced"
)

type ScoredRepo struct {
	ghapi.RepoResult
	Score             float64 `json:"score"`
	ActivityScore     float64 `json:"activity_score"`
	FriendlinessScore float64 `json:"friendliness_score"`
	RelevanceScore    float64 `json:"relevance_score"`
}

func ScoreRepos(repos []ghapi.RepoResult, languages []string, topics []string, skill SkillLevel) []ScoredRepo {
	actW, friW, relW := weights(skill)

	scored := make([]ScoredRepo, 0, len(repos))
	for _, r := range repos {
		act := activityScore(r)
		fri := friendlinessScore(r)
		rel := relevanceScore(r, languages, topics)

		if skill == Beginner && r.GoodFirstIssues == 0 {
			continue
		}

		s := ScoredRepo{
			RepoResult:        r,
			ActivityScore:     math.Round(act*10) / 10,
			FriendlinessScore: math.Round(fri*10) / 10,
			RelevanceScore:    math.Round(rel*10) / 10,
			Score:             math.Round((act*actW+fri*friW+rel*relW)*10) / 10,
		}
		scored = append(scored, s)
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

func weights(skill SkillLevel) (activity, friendliness, relevance float64) {
	switch skill {
	case Beginner:
		return 0.25, 0.50, 0.25
	case Advanced:
		return 0.55, 0.20, 0.25
	default:
		return 0.40, 0.35, 0.25
	}
}

func activityScore(r ghapi.RepoResult) float64 {
	starScore := math.Min(math.Log10(float64(r.Stars+1))/5.0*100, 100)
	forkScore := math.Min(math.Log10(float64(r.Forks+1))/4.0*100, 100)
	recency := recencyScore(r.LastPushedAt)

	return starScore*0.3 + forkScore*0.2 + recency*0.5
}

func recencyScore(lastPush time.Time) float64 {
	if lastPush.IsZero() {
		return 0
	}
	days := time.Since(lastPush).Hours() / 24
	switch {
	case days <= 7:
		return 100
	case days <= 30:
		return 75
	case days <= 90:
		return 50
	case days <= 365:
		return 25
	default:
		return 0
	}
}

func friendlinessScore(r ghapi.RepoResult) float64 {
	gfi := math.Min(float64(r.GoodFirstIssues)/5.0*100, 100)
	hw := math.Min(float64(r.HelpWantedCount)/5.0*100, 100)
	var contrib float64
	if r.HasContribGuide {
		contrib = 100
	}

	return gfi*0.4 + hw*0.3 + contrib*0.3
}

func relevanceScore(r ghapi.RepoResult, languages []string, topics []string) float64 {
	var langScore float64
	for _, lang := range languages {
		if strings.EqualFold(r.Language, lang) {
			langScore = 100
			break
		}
	}

	var topicScore float64
	if len(topics) > 0 {
		matched := 0
		repoTopics := make(map[string]bool)
		for _, t := range r.Topics {
			repoTopics[strings.ToLower(t)] = true
		}
		for _, t := range topics {
			if repoTopics[strings.ToLower(t)] {
				matched++
			}
		}
		topicScore = float64(matched) / float64(len(topics)) * 100
	} else {
		topicScore = 50 // neutral when no topics specified
	}

	return langScore*0.6 + topicScore*0.4
}

func ParseSkillLevel(s string) (SkillLevel, bool) {
	switch strings.ToLower(s) {
	case "beginner":
		return Beginner, true
	case "intermediate":
		return Intermediate, true
	case "advanced":
		return Advanced, true
	default:
		return "", false
	}
}
