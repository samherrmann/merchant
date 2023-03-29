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

// addMetafields adds the "metafields" flag to the the given command.
func addMetafields(cmd *cobra.Command) *bool {
	return cmd.Flags().Bool(
		"metafields",
		false,
		"Pull product and variant metafields",
	)
}
