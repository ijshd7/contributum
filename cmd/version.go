package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of contributum",
	Run: func(cmd *cobra.Command, args []string) {
		style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
		fmt.Printf("%s %s\n", style.Render("contributum"), Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
