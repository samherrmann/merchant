// Package shop provides a client to communicate with a shopify store.
package shop

import (
	"errors"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/cache"
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/memdb"
)

var (
	ErrNotExist = errors.New("does not exist")
)

type Client struct {
	*goshopify.Client
}

func NewClient(c *config.StoreConfig) *Client {
	return &Client{
		Client: goshopify.NewClient(
			goshopify.App{
				ApiKey:   c.APIKey,
				Password: c.Password,
			},
			c.Name,
			"",
			goshopify.WithRetry(3),
		),
	}
}

func (c *Client) GetProduct(id int64, skipCache bool) (*goshopify.Product, error) {
	p, err := c.getProduct(id, skipCache)
	// Fall back to pulling product from store if no cache exists.
	if !skipCache && os.IsNotExist(err) {
		return c.getProduct(id, true)
	}
	return p, err
}

func (c *Client) GetProductWithMetafields(id int64) (*goshopify.Product, error) {
	product, err := c.Product.Get(id, nil)
	if err != nil {
		return nil, err
	}
	if err := c.attachMetafields(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (c *Client) GetVariantBySKU(sku string) (*goshopify.Variant, error) {
	return c.searchVariant(func(v *goshopify.Variant) bool {
		return v.Sku == sku
	})
}

func (c *Client) GetVariantByBarcode(barcode string) (*goshopify.Variant, error) {
	return c.searchVariant(func(v *goshopify.Variant) bool {
		return v.Barcode == barcode
	})
}

func (c *Client) GetInventory(skipCache bool) ([]goshopify.Product, error) {
	p, err := c.getInventory(skipCache)
	// Fall back to pulling product from store if no cache exists.
	if !skipCache && os.IsNotExist(err) {
		return c.getInventory(true)
	}
	return p, err
}

func (c *Client) GetInventoryWithMetafields() ([]goshopify.Product, error) {
	products, err := c.listProducts(nil)
	if err != nil {
		return nil, err
	}
	for i, product := range products {
		fmt.Printf("Getting metafields for product %v\n", product.ID)
		if err := c.attachMetafields(&product); err != nil {
			return nil, err
		}
		// TODO check if this is necessary.
		products[i] = product
	}
	return products, nil
}

// GetVariantCount returns the total number of variants for all products.
func (c *Client) GetVariantCount() (int, error) {
	options := &goshopify.ListOptions{
		Fields: "variants",
	}
	products, err := c.listProducts(options)
	if err != nil {
		return 0, err
	}
	variantIds := []int64{}
	for _, p := range products {
		for _, v := range p.Variants {
			variantIds = append(variantIds, v.ID)
		}
	}
	return len(variantIds), nil
}

// UpdateProducts updates the given products in the store.
func (c *Client) UpdateProducts(products []goshopify.Product) error {
	// Get latest inventory from live store so that we don't accidentally make
	// updates based on an outdated cache.
	inventory, err := c.GetInventory(true)
	if err != nil {
		return err
	}
	db, err := memdb.New(inventory)
	if err != nil {
		return err
	}
	operations, err := db.Operations(products)
	if err != nil {
		return err
	}
	errs := []error{}
	for _, p := range operations.NewProducts {
		if _, err := c.Product.Create(p); err != nil {
			errs = append(errs, err)
		}
	}
	for _, p := range operations.ProductUpdates {
		if _, err := c.Product.Update(p); err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range operations.NewVariants {
		if _, err := c.Variant.Create(v.ProductID, v); err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range operations.VariantUpdates {
		if _, err := c.Variant.Update(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (c *Client) getProduct(id int64, skipCache bool) (*goshopify.Product, error) {
	if skipCache {
		p, err := c.Product.Get(id, nil)
		if err != nil {
			return nil, err
		}
		if err := cache.WriteProductFile(p); err != nil {
			return nil, err
		}
	}
	return cache.ReadProductFile(id)
}

// getInventory gets all products from the store and stores them in the cache
// file.
func (c *Client) getInventory(skipCache bool) ([]goshopify.Product, error) {
	if skipCache {
		p, err := c.listProducts(nil)
		if err != nil {
			return nil, err
		}
		if err := cache.WriteInventoryFile(p); err != nil {
			return nil, err
		}
	}
	return cache.ReadInventoryFile()
}

func (c *Client) listProducts(options *goshopify.ListOptions) ([]goshopify.Product, error) {
	products := []goshopify.Product{}
	defaultOptions := &goshopify.ListOptions{
		// 250 is the maximum limit
		// https://shopify.dev/api/admin/rest/reference/products/product?api%5Bversion%5D=2020-10#endpoints-2020-10
		Limit: 250,
	}
	if options == nil {
		options = defaultOptions
	}
	if options.Limit == 0 {
		options.Limit = defaultOptions.Limit
	}
	for {
		productsPacket, pagination, err := c.Product.ListWithPagination(options)
		if err != nil {
			return nil, fmt.Errorf("failed to get packet of products: %w", err)
		}
		products = append(products, productsPacket...)
		if pagination.NextPageOptions == nil {
			break
		}
		options = pagination.NextPageOptions
	}
	return products, nil
}

// attachMetafields fetches and attaches all metafields for the given product and its variants.
func (c *Client) attachMetafields(product *goshopify.Product) error {
	metafields, err := c.Product.ListMetafields(product.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to get metafields for product %v: %w", product.ID, err)
	}
	product.Metafields = metafields

	for j, variant := range product.Variants {
		metafields, err := c.Variant.ListMetafields(variant.ID, nil)
		if err != nil {
			return fmt.Errorf("failed to get metafields for variant %v: %w", variant.ID, err)
		}
		product.Variants[j].Metafields = metafields
	}
	return nil
}

func (c *Client) searchVariant(fn func(v *goshopify.Variant) bool) (*goshopify.Variant, error) {
	inventory, err := c.getInventory(false)
	if err != nil {
		return nil, err
	}
	for _, p := range inventory {
		for _, v := range p.Variants {
			if fn(&v) {
				return &v, nil
			}
		}
	}
	return nil, ErrNotExist
}
