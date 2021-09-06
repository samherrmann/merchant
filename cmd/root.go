package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/spf13/cobra"
)

var (
	appName              = "shopctl"
	defaultCacheFilename = "products.json"
	defaultCSVFilename   = "products.csv"
	rootCmd              = &cobra.Command{}
	shopClient           *goshopify.Client
	metafieldDefinitions *MetafieldDefinitions
)

func Execute() error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	metafieldDefinitions = &config.MetafieldDefinitions
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
	configFilename := "shopctl.json"

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

func writeCacheFile(filename string, data []byte) error {
	dir, err := cacheDir()
	if err != nil {
		return err
	}
	// We first join the filename with the cache directory and then call
	// filepath.Dir so that if filename includes a directory that doen't exist
	// yet that we can create it before writing the file.
	path := filepath.Join(dir, filename)
	dir = filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func readCacheFile(filename string) ([]byte, error) {
	dir, err := cacheDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, filename)
	return os.ReadFile(path)
}

func cacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appName), nil
}

type Config struct {
	ShopName             string               `json:"shopName"`
	APIKey               string               `json:"apiKey"`
	Password             string               `json:"password"`
	MetafieldDefinitions MetafieldDefinitions `json:"metafieldDefinitions"`
}

type MetafieldDefinitions struct {
	Product []MetafieldDefinition `json:"product"`
	Variant []MetafieldDefinition `json:"variant"`
}

type MetafieldDefinition struct {
	Key string `json:"key"`
}
