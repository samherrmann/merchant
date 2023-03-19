// Package cache handles the local caching of store data.
package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/exec"
)

const (
	inventoryFilename = "inventory.json"
)

// Dir returns the path to the cache directory. If the directory does not exist,
// then Dir will create it.
func Dir() (string, error) {
	cacheRootDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(cacheRootDir, config.AppName)
	err = os.MkdirAll(cacheDir, os.ModePerm)
	return cacheDir, err
}

func ReadProductFile(id int64) (*goshopify.Product, error) {
	bytes, err := readFile(fmt.Sprintf("%v.json", id))
	if err != nil {
		return nil, err
	}
	products := &goshopify.Product{}
	if err = json.Unmarshal(bytes, products); err != nil {
		return nil, err
	}
	return products, nil
}

func ReadInventoryFile() ([]goshopify.Product, error) {
	bytes, err := readFile(inventoryFilename)
	if err != nil {
		return nil, err
	}
	products := []goshopify.Product{}
	if err = json.Unmarshal(bytes, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func WriteProductFile(product *goshopify.Product) error {
	bytes, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(fmt.Sprintf("%v.json", product.ID), bytes)
}

func WriteInventoryFile(products []goshopify.Product) error {
	bytes, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(inventoryFilename, bytes)
}

// PrintEntries writes all cache entries to the given writer.
func PrintEntries(w io.Writer) error {
	entries, err := readDir()
	if err != nil {
		return err
	}
	// Print entries table
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "FILE", "MODIFIED")
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%v\t%v\n", removeExt(entry.Name()), info.ModTime())
	}
	tw.Flush()
	return nil
}

// OpenFileInTextEditor opens a cache file in a text editor.
func OpenFileInTextEditor(filename string) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	filename = filepath.Join(dir, filename)
	return exec.RunTextEditor(filename)
}

func readDir() ([]fs.DirEntry, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	return os.ReadDir(dir)
}

func readFile(filename string) ([]byte, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, filename)
	return os.ReadFile(path)
}

func writeFile(filename string, data []byte) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	// We first join the filename with the cache directory and then call
	// filepath.Dir so that if filename includes a directory that doesn't exist
	// yet then we can create it before writing the file.
	path := filepath.Join(dir, filename)
	dir = filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// removeExt returns filename without the extension
func removeExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
