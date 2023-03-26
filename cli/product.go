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
