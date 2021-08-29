package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func main() {
	envFilename := "goshopctl.json"

	bytes, err := ioutil.ReadFile(envFilename)
	if err != nil {
		log.Fatalf("Error reading %v: %v\n", envFilename, err)
	}

	config := &Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		log.Fatalf("Error parsing %v: %v", envFilename, err)
	}

	client := goshopify.NewClient(
		goshopify.App{
			ApiKey:   config.APIKey,
			Password: config.Password,
		},
		config.ShopName,
		"",
	)

	products := []goshopify.Product{}
	options := &goshopify.ListOptions{
		// 250 is the maximum limit
		// https://shopify.dev/api/admin/rest/reference/products/product?api%5Bversion%5D=2020-10#endpoints-2020-10
		Limit: 250,
	}
	for {
		productsPacket, pagination, err := client.Product.ListWithPagination(options)
		if err != nil {
			log.Fatalf("Error getting list of products: %v\n", err)
		}
		products = append(products, productsPacket...)
		if pagination.NextPageOptions == nil {
			break
		}
		options = pagination.NextPageOptions
	}

	productsFilename := "products.csv"

	productsFile, err := os.OpenFile(productsFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening %v: %v", productsFilename, err)
	}
	defer productsFile.Close()

	csvWriter := csv.NewWriter(productsFile)

	header := []string{
		"ID",
	}
	csvWriter.Write(header)

	for _, p := range products {
		row := []string{
			strconv.FormatInt(p.ID, 10),
		}
		csvWriter.Write(row)
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatalf("CSV write error: %v", err)
	}

	log.Println("Export was successful! :)")
}

type Config struct {
	ShopName string `json:"shopName"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}
