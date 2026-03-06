package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/ijshd7/contributum/internal/ghapi"
	"github.com/ijshd7/contributum/internal/scoring"
)

func RenderRepoResults(repos []scoring.ScoredRepo, asJSON bool) error {
	if asJSON {
		return renderJSON(repos, os.Stdout)
	}
	renderRepoTable(repos, os.Stdout)
	return nil
}

func RenderIssueResults(issues []ghapi.IssueResult, asJSON bool) error {
	if asJSON {
		return renderJSON(issues, os.Stdout)
	}
	renderIssueTable(issues, os.Stdout)
	return nil
}

func renderRepoTable(repos []scoring.ScoredRepo, w io.Writer) {
	if len(repos) == 0 {
		fmt.Fprintln(w, Warning.Render("No repositories found matching your criteria."))
		return
	}

	fmt.Fprintln(w, Title.Render("Repository Search Results"))

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Repository", Width: 30},
		{Title: "Lang", Width: 12},
		{Title: "Stars", Width: 8},
		{Title: "Score", Width: 8},
		{Title: "GFI", Width: 5},
		{Title: "Description", Width: 50},
	}

	rows := make([]table.Row, len(repos))
	for i, r := range repos {
		rows[i] = table.Row{
			fmt.Sprintf("%d", i+1),
			r.FullName,
			r.Language,
			formatStars(r.Stars),
			fmt.Sprintf("%.0f", r.Score),
			fmt.Sprintf("%d", r.GoodFirstIssues),
			truncate(r.Description, 48),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("99"))
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	fmt.Fprintln(w, t.View())
	fmt.Fprintln(w, Muted.Render(fmt.Sprintf("\nFound %d repositories. GFI = Good First Issues", len(repos))))
}

func renderIssueTable(issues []ghapi.IssueResult, w io.Writer) {
	if len(issues) == 0 {
		fmt.Fprintln(w, Warning.Render("No issues found matching your criteria."))
		return
	}

	fmt.Fprintln(w, Title.Render("Open Issues"))

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Repository", Width: 25},
		{Title: "Title", Width: 40},
		{Title: "Labels", Width: 25},
		{Title: "Age", Width: 10},
		{Title: "Comments", Width: 9},
	}

	rows := make([]table.Row, len(issues))
	for i, issue := range issues {
		rows[i] = table.Row{
			fmt.Sprintf("%d", i+1),
			issue.RepoFullName,
			truncate(issue.Title, 38),
			truncate(strings.Join(issue.Labels, ", "), 23),
			formatAge(issue.CreatedAt),
			fmt.Sprintf("%d", issue.Comments),
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("99"))
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	fmt.Fprintln(w, t.View())
	fmt.Fprintln(w, Muted.Render(fmt.Sprintf("\nFound %d issues.", len(issues))))
}

func renderJSON(data any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

func formatStars(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	days := int(time.Since(t).Hours() / 24)
	switch {
	case days == 0:
		return "today"
	case days == 1:
		return "1 day"
	case days < 30:
		return fmt.Sprintf("%dd", days)
	case days < 365:
		return fmt.Sprintf("%dmo", days/30)
	default:
		return fmt.Sprintf("%dy", days/365)
	}
}
