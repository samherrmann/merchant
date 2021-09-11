package cmd

import (
	"github.com/samherrmann/shopctl/cache"
	"github.com/spf13/cobra"
)

func newCacheListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List files in cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cache.PrintEntries()
		},
	}
}