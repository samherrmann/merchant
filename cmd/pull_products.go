package cmd

import (
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/cache"
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/csv"
	"github.com/samherrmann/shopctl/shop"
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
				products, err = pullInventory(shopClient)
				if err != nil {
					return err
				}
			} else {
				// Try reading from cache.
				products, err = cache.ReadInventoryFile()
				if os.IsNotExist(err) {
					// Cache does not exist, therefore pull from store.
					products, err = pullInventory(shopClient)
					if err != nil {
						return err
					}
				} else if err != nil {
					return err
				}
			}
			return csv.WriteInventoryFile(products, metafieldDefs)
		},
	}
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}

func pullInventory(shopClient *shop.Client) ([]goshopify.Product, error) {
	products, err := shopClient.GetInventoryWithMetafields()
	if err != nil {
		return nil, err
	}
	if err := cache.WriteInventoryFile(products); err != nil {
		return nil, err
	}
	return products, nil
}
