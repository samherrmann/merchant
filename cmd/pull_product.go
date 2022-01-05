package cmd

import (
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/cache"
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/csv"
	"github.com/samherrmann/shopctl/exec"
	"github.com/samherrmann/shopctl/shop"
	"github.com/samherrmann/shopctl/utils"
	"github.com/spf13/cobra"
)

func newPullProductCommand(metafieldDefs *config.MetafieldDefinitions) *cobra.Command {
	var skipCache *bool
	var openFile *bool
	cmd := &cobra.Command{
		Use:   "product <id>|inventory",
		Short: "Fetch product and its metadata from store",
		Long:  "Fetch a single product or the entire product inventory from the store. The product metafields are included.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]
			if arg == "inventory" {
				products, err := getInventory(shopClient, *skipCache)
				if err != nil {
					return err
				}
				csv.WriteInventoryFile(products, metafieldDefs)
				return nil
			}
			productID, err := utils.ParseID(arg)
			if err != nil {
				return err
			}
			product, err := getProduct(productID, shopClient, *skipCache)
			if err != nil {
				return err
			}
			return csv.WriteProductFile(product, metafieldDefs)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if *openFile {
				return exec.RunSpreadsheetEditor(args[0] + ".csv")
			}
			return nil
		},
	}
	openFile = cmd.Flags().Bool("open", false, "Open product file after pulling")
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}

// getInventory returns the entire product inventory from the cache or from the
// store.
func getInventory(shopClient *shop.Client, skipCache bool) ([]goshopify.Product, error) {
	if skipCache {
		return pullAndCacheInventory(shopClient)
	}
	products, err := cache.ReadInventoryFile()
	if os.IsNotExist(err) {
		return pullAndCacheInventory(shopClient)
	}
	return products, err
}

// getProduct returns the product specified by ID from the cache or from the
// store.
func getProduct(id int64, shopClient *shop.Client, skipCache bool) (*goshopify.Product, error) {
	if skipCache {
		return pullAndCacheProduct(id, shopClient)
	}
	product, err := cache.ReadProductFile(id)
	if os.IsNotExist(err) {
		return pullAndCacheProduct(id, shopClient)
	}
	return product, err
}

// pullAndCacheInventory pulls the entire product inventory from the store and
// writes it to a cache file.
func pullAndCacheInventory(shopClient *shop.Client) ([]goshopify.Product, error) {
	products, err := shopClient.GetInventoryWithMetafields()
	if err != nil {
		return nil, err
	}
	if err := cache.WriteInventoryFile(products); err != nil {
		return nil, err
	}
	return products, nil
}

// pullAndCacheProduct pulls the product specified by ID from the store and
// writes it to a cache file.
func pullAndCacheProduct(id int64, shopClient *shop.Client) (*goshopify.Product, error) {
	product, err := shopClient.GetProductWithMetafields(id)
	if err != nil {
		return nil, err
	}
	if err := cache.WriteProductFile(product); err != nil {
		return nil, err
	}
	return product, nil
}
