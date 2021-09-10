package cmd

import (
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/goshopctl/cache"
	"github.com/samherrmann/goshopctl/config"
	"github.com/samherrmann/goshopctl/csv"
	"github.com/samherrmann/goshopctl/shop"
	"github.com/samherrmann/goshopctl/utils"
	"github.com/spf13/cobra"
)

func newPullProductCommand(shopClient *shop.Client, metafieldDefs *config.MetafieldDefinitions) *cobra.Command {
	var skipCache *bool
	cmd := &cobra.Command{
		Use:   "product <id>",
		Short: "Fetch product and its metadata from store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get product for ID given as argument.
			productID, err := utils.ParseID(args[0])
			if err != nil {
				return err
			}
			var product *goshopify.Product
			if *skipCache {
				// Ignore cache that may exist and just pull from store.
				product, err = pullProduct(shopClient, productID)
				if err != nil {
					return err
				}
			} else {
				// Try reading from cache.
				product, err = cache.ReadProductFile(productID)
				if os.IsNotExist(err) {
					// Cache does not exist, therefore pull from store.
					product, err = pullProduct(shopClient, productID)
					if err != nil {
						return err
					}
				} else if err != nil {
					return err
				}
			}
			return csv.WriteProductFile(product, metafieldDefs)
		},
	}
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}

func pullProduct(shopClient *shop.Client, id int64) (*goshopify.Product, error) {
	product, err := shopClient.GetProductWithMetafields(id)
	if err != nil {
		return nil, err
	}
	if err := cache.WriteProductFile(product); err != nil {
		return nil, err
	}
	return product, nil
}
