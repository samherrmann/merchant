// Package config manages the application configurations.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samherrmann/merchant/editor"
	"github.com/samherrmann/merchant/shopify"
)

// Build-time variables set by -ldflags.
const (
	AppName = "merchant"
	Version = "dev"
)

// New returns a new configuration.
func New() *Config {
	// We need to manually initialize slices to be able to marshal the config to
	// JSON without fields that are of type array set to null.
	// https://github.com/golang/go/issues/27589
	return &Config{
		Store: shopify.Configuration{},
		MetafieldDefinitions: MetafieldDefinitions{
			Product: []MetafieldDefinition{},
			Variant: []MetafieldDefinition{},
		},
		SpreadsheetEditor: DefaultSpreadsheetEditor,
		TextEditor:        DefaultTextEditor,
	}
}

type Config struct {
	// Store contains the Shopify store access information.
	Store shopify.Configuration `json:"store"`
	// MetafieldDefinitions contains metafield definitions.
	MetafieldDefinitions MetafieldDefinitions `json:"metafieldDefinitions"`
	// TextEditorCmd is the command that launches the text editor.
	TextEditor []string `json:"textEditor"`
	// SpreadsheetEditor is the command that launches the spreadsheet editor.
	SpreadsheetEditor []string `json:"spreadsheetEditor"`
}

// UnmarshalJSON implements the [encoding/json.Unmarshaler] interface.
func (c *Config) UnmarshalJSON(b []byte) error {
	type alias Config
	a := alias{}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	if len(a.SpreadsheetEditor) == 0 {
		a.SpreadsheetEditor = DefaultSpreadsheetEditor
	}
	if len(a.TextEditor) == 0 {
		a.TextEditor = DefaultTextEditor
	}
	*c = Config(a)
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

// InitFile write a default configuration file. os.ErrExist is returned if the
// file already exists.
func InitFile() (*Config, error) {
	dir, err := mkDefaultDir()
	if err != nil {
		return nil, err
	}
	return initFile(dir)
}

// Load loads the configuration from file.
func Load() (*Config, error) {
	dir, err := mkDefaultDir()
	if err != nil {
		return nil, err
	}
	return load(dir)
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
func (c *Config) OpenInTextEditor() error {
	dir, err := mkDefaultDir()
	if err != nil {
		return err
	}
	return c.newTextEditor().Open(joinFilename(dir))
}

// initFile write the default configuration file to the given directory.
// os.ErrExist is returned if the file already exists.
func initFile(dir string) (*Config, error) {
	filename := joinFilename(dir)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := New()
	bytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return nil, err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Load loads the configuration from file located in dir.
func load(dir string) (*Config, error) {
	filename := joinFilename(dir)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("file %v: %w", filename, err)
	}
	return cfg, nil
}

func (c *Config) newTextEditor() editor.Editor {
	cmd := c.TextEditor
	if len(cmd) == 0 {
		cmd = DefaultTextEditor
	}
	return editor.New(cmd...)

}

// mkDefaultDir creates and returns the app configuration directory, rooted in
// the default OS directory to use for user-specific configuration data.
func mkDefaultDir() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return mkdir(root)
}

// mkdir creates and returns the app configuration directory. No error is
// returned if the directory already exists.
func mkdir(root string) (string, error) {
	dir := filepath.Join(root, AppName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

// joinFilename joins the configuration filename with dir.
func joinFilename(dir string) string {
	return filepath.Join(dir, AppName) + ".json"
}
