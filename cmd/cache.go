package cmd

import "github.com/spf13/cobra"

func newCacheCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cache",
		Short: "Manage cache",
	}
}
