package cmd

import (
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/goshopctl/cache"
	"github.com/samherrmann/goshopctl/config"
	"github.com/samherrmann/goshopctl/csv"
	"github.com/samherrmann/goshopctl/shop"
	"github.com/spf13/cobra"
)

func newPullProductsCommand(shopClient *shop.Client, metafieldDefs *config.MetafieldDefinitions) *cobra.Command {
	var skipCache *bool
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Fetch products and their metadata from store",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var products []goshopify.Product
			var err error
			if *skipCache {
				// Ignore cache that may exist and just pull from store.
				products, err = pullProducts(shopClient)
				if err != nil {
					return err
				}
			} else {
				// Try reading from cache.
				products, err = cache.ReadProductsFile()
				if os.IsNotExist(err) {
					// Cache does not exist, therefore pull from store.
					products, err = pullProducts(shopClient)
					if err != nil {
						return err
					}
				} else if err != nil {
					return err
				}
			}
			return csv.WriteProductsFile(products, metafieldDefs)
		},
	}
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}

func pullProducts(shopClient *shop.Client) ([]goshopify.Product, error) {
	products, err := shopClient.GetProductsWithMetafields()
	if err != nil {
		return nil, err
	}
	if err := cache.WriteProductsFile(products); err != nil {
		return nil, err
	}
	return products, nil
}
