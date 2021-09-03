package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmdPull)
}

var cmdPull = &cobra.Command{
	Use:   "pull",
	Short: "Fetch products and their metadata from the store",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		products, err := readCache()
		if err != nil {
			// Get products from store if failed to read from cache.
			products, err = getProductsWithMetafields(shopClient)
			if err != nil {
				return err
			}
			if err := writeCache(products); err != nil {
				return err
			}
		}

		csvRows, err := convertProductsToCSVRows(products)
		if err != nil {
			return err
		}

		if err := writeCSVFile(csvRows); err != nil {
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

func readCache() ([]goshopify.Product, error) {
	bytes, err := ioutil.ReadFile(cacheFilename)
	if err != nil {
		return nil, err
	}
	products := []goshopify.Product{}
	if err = json.Unmarshal(bytes, &products); err != nil {
		return nil, err
	}
	return products, err
}

func writeCache(products []goshopify.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(cacheFilename, bytes, 0644); err != nil {
		return err
	}
	return nil
}

func writeCSVFile(rows []CSVRow) error {
	bytes, err := csvutil.Marshal(rows)
	if err != nil {
		return fmt.Errorf("error encoding to CSV: %w", err)
	}

	return ioutil.WriteFile(csvFilename, bytes, 0644)
}

func convertProductsToCSVRows(products []goshopify.Product) ([]CSVRow, error) {
	rows := []CSVRow{}
	for _, p := range products {
		for _, m := range p.Metafields {
			row, err := convertMetafieldToCSVRow(p.ID, 0, m)
			if err != nil {
				return nil, err
			}
			rows = append(rows, *row)
		}

		for _, v := range p.Variants {
			for _, m := range p.Metafields {
				row, err := convertMetafieldToCSVRow(p.ID, v.ID, m)
				if err != nil {
					return nil, err
				}
				rows = append(rows, *row)
			}
		}
	}
	return rows, nil
}

func convertMetafieldToCSVRow(productID int64, variantID int64, metafield goshopify.Metafield) (*CSVRow, error) {
	row := &CSVRow{
		ProductID:      productID,
		VariantID:      variantID,
		MetafiledID:    metafield.ID,
		MetafieldKey:   metafield.Key,
		MetafieldValue: metafield.Value,
	}

	if metafield.ValueType == "json_string" {
		measurement := &Measurement{}
		if err := json.Unmarshal([]byte(fmt.Sprint(metafield.Value)), measurement); err != nil {
			return nil, fmt.Errorf("error unmarshaling metafield JSON string: %w", err)
		}
		row.MetafieldValue = measurement.Value
		row.MetafieldUnit = measurement.Unit
	}

	return row, nil
}

type CSVRow struct {
	ProductID      int64       `csv:"product_id"`
	VariantID      int64       `csv:"variant_id,omitempty"`
	MetafiledID    int64       `csv:"metafiled_id"`
	MetafieldKey   string      `csv:"metafield_key"`
	MetafieldValue interface{} `csv:"metafield_value"`
	MetafieldUnit  string      `csv:"metafield_unit,omitempty"`
}

type Measurement struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
