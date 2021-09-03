package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/spf13/cobra"
)

var (
	cacheFilename = "products.json"
	csvFilename   = "products.csv"
	rootCmd       = &cobra.Command{}
	shopClient    *goshopify.Client
)

func Execute() error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	shopClient = newClient(config)
	return rootCmd.Execute()
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

func parseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

type Config struct {
	ShopName string `json:"shopName"`
	APIKey   string `json:"apiKey"`
	Password string `json:"password"`
}
