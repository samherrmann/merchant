package cli

import (
	"os"

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

			dir, err := cache.Dir()
			if err != nil {
				return err
			}
			return os.RemoveAll(dir)
		},
	}
}
