package cmd

import (
	"encoding/json"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull [<productid>]",
	Short: "Fetch products and their metadata from the store",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get all products if no product ID is provided.
		if len(args) == 0 {
			products, err := getProductsWithMetafields()
			if err != nil {
				return err
			}
			return writeProductsFile(products)
		}
		// Get product for ID given as argument.
		productID, err := parseID(args[0])
		if err != nil {
			return err
		}
		product, err := getProductWithMetafields(productID)
		if err != nil {
			return err
		}
		return writeProductFile(product)
	},
}

func getProductWithMetafields(id int64) (*goshopify.Product, error) {
	product, err := shopClient.Product.Get(id, nil)
	if err != nil {
		return nil, err
	}
	if err := attachMetafields(product); err != nil {
		return nil, err
	}
	return product, nil
}

func getProductsWithMetafields() ([]goshopify.Product, error) {
	products := []goshopify.Product{}
	options := &goshopify.ListOptions{
		// 250 is the maximum limit
		// https://shopify.dev/api/admin/rest/reference/products/product?api%5Bversion%5D=2020-10#endpoints-2020-10
		Limit: 250,
	}
	for {
		productsPacket, pagination, err := shopClient.Product.ListWithPagination(options)
		if err != nil {
			return nil, fmt.Errorf("failed to get packet of products: %w", err)
		}

		for i, product := range productsPacket {
			fmt.Printf("Getting metafields for product %v\n", product.ID)
			if err := attachMetafields(&product); err != nil {
				return nil, err
			}
			// TODO check if this is necessary.
			productsPacket[i] = product
		}

		products = append(products, productsPacket...)
		if pagination.NextPageOptions == nil {
			break
		}
		options = pagination.NextPageOptions
	}
	return products, nil
}

// attachMetafields fetches and attaches all metafields for the given product and its variants.
func attachMetafields(product *goshopify.Product) error {
	metafields, err := shopClient.Product.ListMetafields(product.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to get metafields for product %v: %w", product.ID, err)
	}
	product.Metafields = metafields

	for j, variant := range product.Variants {
		metafields, err := shopClient.Variant.ListMetafields(variant.ID, nil)
		if err != nil {
			return fmt.Errorf("failed to get metafields for variant %v: %w", variant.ID, err)
		}
		product.Variants[j].Metafields = metafields
	}
	return nil
}

func writeProductFile(product *goshopify.Product) error {
	bytes, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		return err
	}
	return writeCacheFile(fmt.Sprintf("%v.json", product.ID), bytes)
}

func writeProductsFile(products []goshopify.Product) error {
	bytes, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return writeCacheFile(defaultCacheFilename, bytes)
}
