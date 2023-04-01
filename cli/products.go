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
