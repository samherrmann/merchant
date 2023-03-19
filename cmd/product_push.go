package cmd

import (
	"github.com/samherrmann/merchant/csv"
	"github.com/spf13/cobra"
)

func newProductPushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push <filename>",
		Short: "Update products in store with data from CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			products, err := csv.ReadProducts(args[0])
			if err != nil {
				return err
			}
			return shopClient.UpdateProducts(products)
		},
	}
}
