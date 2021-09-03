package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Fetch products and their metadata from the store",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		products, err := getProductsWithMetafields(shopClient)
		if err != nil {
			return err
		}
		if err := writeProductsFile(products); err != nil {
			return err
		}
		return nil
	},
}

func getProductsWithMetafields(client *goshopify.Client) ([]goshopify.Product, error) {
	products := []goshopify.Product{}
	options := &goshopify.ListOptions{
		// 250 is the maximum limit
		// https://shopify.dev/api/admin/rest/reference/products/product?api%5Bversion%5D=2020-10#endpoints-2020-10
		Limit: 250,
	}
	for {
		productsPacket, pagination, err := client.Product.ListWithPagination(options)
		if err != nil {
			return nil, fmt.Errorf("failed to get packet of products: %w", err)
		}

		for i, product := range productsPacket {
			log.Printf("Getting metafields for product %v\n", product.ID)
			metafields, err := client.Product.ListMetafields(product.ID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get metafields for product %v: %w", product.ID, err)
			}
			productsPacket[i].Metafields = metafields

			for j, variant := range product.Variants {
				log.Printf("Getting metafields for variant %v\n", variant.ID)
				metafields, err := client.Variant.ListMetafields(variant.ID, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to get metafields for variant %v: %w", variant.ID, err)
				}
				productsPacket[i].Variants[j].Metafields = metafields
			}
		}

		products = append(products, productsPacket...)
		if pagination.NextPageOptions == nil {
			break
		}
		options = pagination.NextPageOptions
	}

	return products, nil
}

func writeProductsFile(products []goshopify.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(productsFilename, bytes, 0644)
}
