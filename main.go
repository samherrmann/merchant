package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/jszwec/csvutil"
)

const (
	cacheFilename = "products.json"
	csvFilename   = "metafields.csv"
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	client := newClient(config)

	products, err := readCache()
	if err != nil {
		// Get products from store if failed to read from cache.
		products, err = getProductsWithMetafields(client)
		if err != nil {
			log.Fatalln(err)
		}
		if err := writeCache(products); err != nil {
			log.Fatalln(err)
		}
	}

	rows := convertProductsToCSVRows(products)

	if err := writeCSVFile(rows); err != nil {
		log.Fatalln(err)
	}

	log.Println("Export was successful! :)")
}

// TODO remove
func sampleProductMetafieldCreate(client *goshopify.Client) (*goshopify.Metafield, error) {
	return client.Product.CreateMetafield(6573170753578, goshopify.Metafield{
		Key:       "box_per_carton",
		Value:     123,
		ValueType: "integer",
		Namespace: "common",
	})
}

func newClient(c *Config) *goshopify.Client {
	return goshopify.NewClient(
		goshopify.App{
			ApiKey:   c.APIKey,
			Password: c.Password,
		},
		c.ShopName,
		"",
		goshopify.WithRetry(3),
	)
}

func readConfig() (*Config, error) {
	configFilename := "goshopctl.json"

	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %v: %w", configFilename, err)
	}

	config := &Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse %v: %w", configFilename, err)
	}
	return config, nil
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

func convertProductsToCSVRows(products []goshopify.Product) []CSVRow {
	rows := []CSVRow{}
	for _, p := range products {
		for _, m := range p.Metafields {
			rows = append(rows, convertMetafieldToCSVRow(p.ID, 0, m))
		}

		for _, v := range p.Variants {
			for _, m := range p.Metafields {
				rows = append(rows, convertMetafieldToCSVRow(p.ID, v.ID, m))
			}
		}
	}
	return rows
}

func convertMetafieldToCSVRow(productID int64, variantID int64, metafield goshopify.Metafield) CSVRow {
	return CSVRow{
		ProductID:      productID,
		VariantID:      variantID,
		MetafiledID:    metafield.ID,
		MetafieldKey:   metafield.Key,
		MetafieldValue: metafield.Value,
	}
}

type Config struct {
	ShopName string `json:"shopName"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}

type CSVRow struct {
	ProductID      int64       `csv:"product_id"`
	VariantID      int64       `csv:"variant_id,omitempty"`
	MetafiledID    int64       `csv:"metafiled_id"`
	MetafieldKey   string      `csv:"metafield_key"`
	MetafieldValue interface{} `csv:"metafield_value"`
}
