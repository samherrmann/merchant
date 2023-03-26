package cli

import (
	"github.com/spf13/cobra"
)

func newProductCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "product",
		Short: "Manage products",
	}
}

// addCacheFlag adds the "skip-cache" flag the the given command.
func addCacheFlag(cmd *cobra.Command) *bool {
	return cmd.Flags().Bool(
		"skip-cache",
		false,
		"Pull directly from store even if local copy exists in cache",
	)
}
