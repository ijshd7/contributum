package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contributum",
	Short: "Find open-source repos to contribute to",
	Long: `Contributum helps you discover open-source repositories that match your
skills and interests. Search by language, topic, and skill level to find
projects that are actively looking for contributors.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
