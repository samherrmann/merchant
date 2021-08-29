package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

	clientConfig := goshopify.App{
		ApiKey:   config.APIKey,
		Password: config.Password,
	}

	client := goshopify.NewClient(clientConfig, config.ShopName, "")

	count, err := client.Product.Count(nil)
	if err != nil {
		log.Fatalf("Error getting list of products: %v\n", err)
	}
	log.Printf("%v\n", count)
}

type Config struct {
	ShopName string `json:"shopName"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}
