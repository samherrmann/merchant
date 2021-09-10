package cache

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/tabwriter"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/goshopctl/config"
	"github.com/samherrmann/goshopctl/utils"
)

const (
	productsFilename = "products.json"
)

func WriteFile(filename string, data []byte) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	// We first join the filename with the cache directory and then call
	// filepath.Dir so that if filename includes a directory that doen't exist
	// yet then we can create it before writing the file.
	path := filepath.Join(dir, filename)
	dir = filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ReadFile(filename string) ([]byte, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, filename)
	return os.ReadFile(path)
}

func ReadDir() ([]fs.DirEntry, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	return os.ReadDir(dir)
}

func Dir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, config.AppName), nil
}

func ReadProductFile(id int64) (*goshopify.Product, error) {
	bytes, err := ReadFile(fmt.Sprintf("%v.json", id))
	if err != nil {
		return nil, err
	}
	products := &goshopify.Product{}
	if err = json.Unmarshal(bytes, products); err != nil {
		return nil, err
	}
	return products, nil
}

func ReadProductsFile() ([]goshopify.Product, error) {
	bytes, err := ReadFile(productsFilename)
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
	return WriteFile(fmt.Sprintf("%v.json", product.ID), bytes)
}

func WriteProductsFile(products []goshopify.Product) error {
	bytes, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return WriteFile(productsFilename, bytes)
}

func PrintEntries() error {
	entries, err := ReadDir()
	if err != nil {
		return err
	}
	// Print entries table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\n", "FILE", "MODIFIED")
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%v\t%v\n", utils.RemoveExt(entry.Name()), info.ModTime())
	}
	w.Flush()
	return nil
}
