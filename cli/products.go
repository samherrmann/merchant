package cli

import (
	"github.com/spf13/cobra"
)

func newProductsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "products",
		Short: "Manage products",
	}
}

// addCacheFlag adds the "skip-cache" flag to the the given command.
func addCacheFlag(cmd *cobra.Command) *bool {
	return cmd.Flags().Bool(
		"skip-cache",
		false,
		"Pull directly from store even if local copy exists in cache",
	)
}

// addMetafields adds the "metafields" flag to the the given command.
func addMetafields(cmd *cobra.Command) *bool {
	return cmd.Flags().Bool(
		"metafields",
		false,
		"Pull product and variant metafields",
	)
}
