package cmd

import (
	"github.com/spf13/cobra"
)

func newProductCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "product",
		Short: "Operate on a product",
	}
}
