package cli

import (
	"fmt"
	"io"

	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsCountCommand(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Args:  cobra.NoArgs,
		Short: "Count total number of products and variants in the store",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)
			productCount, err := store.Product.Count(nil)
			if err != nil {
				return err
			}
			variantCount, err := store.GetVariantCount()
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "Products: %v\n", productCount)
			fmt.Fprintf(w, "Variants: %v\n", variantCount)

			return nil
		},
	}
	return cmd
}
