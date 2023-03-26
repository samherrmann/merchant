package cli

import (
	"os"

	"github.com/samherrmann/merchant/cache"
	"github.com/spf13/cobra"
)

func newCacheListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List files in cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			return cache.PrintEntries(os.Stdout)
		},
	}
}
