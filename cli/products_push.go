package cli

import (
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsPushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push <filename>",
		Short: "Update products in store with data from CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)

			products, err := csv.ReadProducts(args[0])
			if err != nil {
				return err
			}
			return store.UpdateProducts(products)
		},
	}
}
