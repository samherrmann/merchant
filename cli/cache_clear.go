package cli

import (
	"github.com/samherrmann/merchant/cache"
	"github.com/spf13/cobra"
)

func newCacheClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clears the entire cache",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true
			return cache.Clear()
		},
	}
}
