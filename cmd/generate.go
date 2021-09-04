package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generates files as defines by sub-command",
}
