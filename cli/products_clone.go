package cli

import (
	"github.com/samherrmann/merchant/cache"
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone products and their metadata from the store into the cache",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)

			products, err := store.GetProducts()
			if err != nil {
				return err
			}

			c, err := cache.New()
			if err != nil {
				return err
			}

			return c.Products().Update(products...)
		},
	}
	return cmd
}
