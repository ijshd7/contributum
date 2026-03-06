package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ijshd7/contributum/internal/ghapi"
	"github.com/ijshd7/contributum/internal/scoring"
	"github.com/ijshd7/contributum/internal/ui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for repositories to contribute to",
	Long: `Search GitHub for open-source repositories matching your language preferences,
topics of interest, and skill level. Results are scored by activity,
contribution-friendliness, and relevance.`,
	Example: `  contributum search --lang go --topic cli --skill beginner
  contributum search --lang python,javascript --topic web --limit 20
  contributum search --lang rust --min-stars 100 --json`,
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().StringP("lang", "l", "", "Languages to search for (comma-separated, required)")
	searchCmd.Flags().StringP("topic", "t", "", "Topics to match (comma-separated)")
	searchCmd.Flags().StringP("skill", "s", "intermediate", "Skill level: beginner, intermediate, advanced")
	searchCmd.Flags().Int("min-stars", 0, "Minimum star count")
	searchCmd.Flags().Int("limit", 10, "Maximum number of results")
	searchCmd.Flags().Bool("json", false, "Output results as JSON")
	_ = searchCmd.MarkFlagRequired("lang")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	langStr, _ := cmd.Flags().GetString("lang")
	topicStr, _ := cmd.Flags().GetString("topic")
	skillStr, _ := cmd.Flags().GetString("skill")
	minStars, _ := cmd.Flags().GetInt("min-stars")
	limit, _ := cmd.Flags().GetInt("limit")
	asJSON, _ := cmd.Flags().GetBool("json")

	languages := splitCSV(langStr)
	topics := splitCSV(topicStr)

	skill, ok := scoring.ParseSkillLevel(skillStr)
	if !ok {
		return fmt.Errorf("invalid skill level %q: must be beginner, intermediate, or advanced", skillStr)
	}

	client := ghapi.NewClient(ghapi.TokenFromEnv())
	params := ghapi.RepoSearchParams{
		Languages: languages,
		Topics:    topics,
		MinStars:  minStars,
		Limit:     limit,
	}

	result, err := ui.RunWithSpinner("Searching GitHub repositories...", func() (any, error) {
		return client.SearchRepos(context.Background(), params)
	})
	if err != nil {
		return err
	}

	repos := result.([]ghapi.RepoResult)
	scored := scoring.ScoreRepos(repos, languages, topics, skill)

	return ui.RenderRepoResults(scored, asJSON)
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
