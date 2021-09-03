package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates files per sub-command.",
}
