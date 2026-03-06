package ui

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	ScoreHigh = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	ScoreMedium = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	ScoreLow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	Label = lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	Warning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))

	Success = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	TableCell = lipgloss.NewStyle().
			Padding(0, 1)
)

func ScoreStyle(score float64) lipgloss.Style {
	switch {
	case score >= 70:
		return ScoreHigh
	case score >= 40:
		return ScoreMedium
	default:
		return ScoreLow
	}
}
