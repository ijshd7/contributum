package cmd

import (
	"context"
	"fmt"

	"github.com/ijshd7/contributum/internal/ghapi"
	"github.com/ijshd7/contributum/internal/ui"
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Find open issues to contribute to",
	Long: `Search GitHub for open issues labeled "good first issue", "help wanted", or
custom labels across repositories matching your language preferences.`,
	Example: `  contributum issues --lang go
  contributum issues --lang python --label "help wanted" --limit 20
  contributum issues --lang rust,go --label "good first issue,bug" --json`,
	RunE: runIssues,
}

func init() {
	issuesCmd.Flags().StringP("lang", "l", "", "Languages to search for (comma-separated, required)")
	issuesCmd.Flags().String("label", "good first issue", "Issue labels to filter by (comma-separated)")
	issuesCmd.Flags().Int("limit", 10, "Maximum number of results")
	issuesCmd.Flags().Bool("json", false, "Output results as JSON")
	_ = issuesCmd.MarkFlagRequired("lang")
	rootCmd.AddCommand(issuesCmd)
}

func runIssues(cmd *cobra.Command, args []string) error {
	langStr, _ := cmd.Flags().GetString("lang")
	labelStr, _ := cmd.Flags().GetString("label")
	limit, _ := cmd.Flags().GetInt("limit")
	asJSON, _ := cmd.Flags().GetBool("json")

	languages := splitCSV(langStr)
	labels := splitCSV(labelStr)

	client := ghapi.NewClient(ghapi.TokenFromEnv())
	params := ghapi.IssueSearchParams{
		Languages: languages,
		Labels:    labels,
		Limit:     limit,
	}

	result, err := ui.RunWithSpinner("Searching for open issues...", func() (any, error) {
		return client.SearchIssues(context.Background(), params)
	})
	if err != nil {
		return fmt.Errorf("searching issues: %w", err)
	}

	issues := result.([]ghapi.IssueResult)

	return ui.RenderIssueResults(issues, asJSON)
}
