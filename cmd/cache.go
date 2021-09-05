package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(cacheCmd)
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Operate on cache",
}
