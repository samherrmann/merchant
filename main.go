package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	client := newClient(config)

	products, err := getProductsWithMetafields(client)
	if err != nil {
		log.Fatalln(err)
	}

	if err := writeJSONFile(products); err != nil {
		log.Fatalln(err)
	}

	log.Println("Export was successful! :)")
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

func loadConfig() (*Config, error) {
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

func writeJSONFile(products []goshopify.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("products.json", bytes, 0644); err != nil {
		return err
	}
	return nil
}

func writeCSVFile(products []goshopify.Product) {
	productsFilename := "products.csv"

	productsFile, err := os.OpenFile(productsFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening %v: %v", productsFilename, err)
	}
	defer productsFile.Close()

	csvWriter := csv.NewWriter(productsFile)

	header := []string{
		"id",
		// "handle",
		// "title",
		// "vendor",
		// "body_html",
		// "product_type",
		// "tags",
		// Options                        []ProductOption `json:"options,omitempty"`,
		// Variants                       []Variant       `json:"variants,omitempty"`,
		// Image                          Image           `json:"image,omitempty"`,
		// Images                         []Image         `json:"images,omitempty"`,
		// TemplateSuffix                 string          `json:"template_suffix,omitempty"`,
		// MetafieldsGlobalTitleTag       string          `json:"metafields_global_title_tag,omitempty"`,
		// MetafieldsGlobalDescriptionTag string          `json:"metafields_global_description_tag,omitempty"`,
		"metafields",
	}
	csvWriter.Write(header)

	for _, p := range products {
		row := []string{
			fmt.Sprintf("%q", strconv.FormatInt(p.ID, 10)),
			// fmt.Sprintf("%q", p.Handle),
			// fmt.Sprintf("%q", p.Title),
			// fmt.Sprintf("%q", p.Vendor),
			// fmt.Sprintf("%q", p.BodyHTML),
			// fmt.Sprintf("%q", p.ProductType),
			// fmt.Sprintf("%q", p.Tags),
			fmt.Sprintf("%q", p.Metafields),
		}
		csvWriter.Write(row)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatalf("CSV write error: %v", err)
	}
}

type Config struct {
	ShopName string `json:"shopName"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}

type Row struct {
	ProductID int64
	variantID int64
}
