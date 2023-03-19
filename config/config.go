// Package config manages the application configurations.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samherrmann/merchant/exec"
)

// Build-time variables set by -ldflags.
const (
	AppName = "merchant"
	Version = "dev"
)

type Config struct {
	// List of Shopify stores.
	Stores StoreConfigs `json:"stores"`
	// MetafieldDefinitions contains metafield definitions.
	MetafieldDefinitions MetafieldDefinitions `json:"metafieldDefinitions"`
	// TextEditorCmd is the command that launches the text editor.
	TextEditor []string `json:"textEditor"`
	// SpreadsheetEditor is the command that launches the spreadsheet editor.
	SpreadsheetEditor []string `json:"spreadsheetEditor"`
}

type StoreConfig struct {
	// Name is the name of the Shopify store as shown in
	// <store-name>.myshopify.com.
	Name string `json:"name"`
	// APIKey is the API key for the Shopify store.
	APIKey string `json:"apiKey"`
	// Password is the password associated with the API key.
	Password string `json:"password"`
}

type StoreConfigs []StoreConfig

// Get returns the configuration for the given name. The first configuration is
// returned if name is an empty string. Nil is returned if no configuration can
// be found for the given name.
func (configs StoreConfigs) Get(name string) *StoreConfig {
	if name == "" {
		return &configs[0]
	}
	for _, c := range configs {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// MetafieldDefinitions define product and variant metafields.
//
// At the time of writing, metafield definitions are not available via the REST
// Admin API. In the meantime, the user must define the metafield definitions in
// the merchant.json file.
// https://shopify.dev/apps/metafields/definitions#structure-of-a-metafield-definition
type MetafieldDefinitions struct {
	Product []MetafieldDefinition `json:"product"`
	Variant []MetafieldDefinition `json:"variant"`
}

type MetafieldDefinition struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	Namespace string `json:"namespace"`
}

func Load() (*Config, error) {
	filename, err := Filename()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		if err := writeSampleFile(filename); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(
			"%v not found, but a sample file was created for you. Run '%v config open' to edit the file",
			filename,
			AppName,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %v: %w", filename, err)
	}

	cFile := &Config{}
	if err := json.Unmarshal(bytes, cFile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %v: %w", filename, err)
	}
	return cFile, nil
}

func Filename() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppName) + ".json", nil
}

func Dir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppName), nil
}

func FindMetafieldDefinition(defs []MetafieldDefinition, namespace string, key string) *MetafieldDefinition {
	for _, def := range defs {
		if namespace == def.Namespace && key == def.Key {
			return &def
		}
	}
	return nil
}

// OpenInTextEditor opens the configuration file in a text editor.
func OpenInTextEditor() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, AppName) + ".json"
	return exec.RunTextEditor(filename)
}

// writeSampleFile write a default configuration file.
func writeSampleFile(filename string) error {
	// https://github.com/golang/go/issues/27589
	c := &Config{
		MetafieldDefinitions: MetafieldDefinitions{
			Product: []MetafieldDefinition{},
			Variant: []MetafieldDefinition{},
		},
		SpreadsheetEditor: []string{},
		TextEditor:        []string{},
	}
	bytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return nil
	}
	return os.WriteFile(filename, bytes, 0644)
}
