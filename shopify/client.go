// Package shopify provides a client to communicate with a shopify store.
package shopify

import (
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

var (
	ErrNotExist = errors.New("does not exist")
)

type Client struct {
	*goshopify.Client
}

func NewClient(c *Configuration) *Client {
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

func (c *Client) GetVariantBySKU(sku string) (*goshopify.Variant, error) {
	return searchVariant(
		c.Product,
		c.Variant,
		func(v *goshopify.Variant) bool {
			return v.Sku == sku
		},
	)
}

func (c *Client) GetVariantByBarcode(barcode string) (*goshopify.Variant, error) {
	return searchVariant(
		c.Product,
		c.Variant,
		func(v *goshopify.Variant) bool {
			return v.Barcode == barcode
		},
	)
}

func (c *Client) GetProducts() ([]goshopify.Product, error) {
	return getProducts(c.Product, c.Variant)
}

// GetVariantCount returns the total number of variants for all products.
func (c *Client) GetVariantCount() (int, error) {
	return getVariantCount(c.Product)
}

// UpdateProducts updates the given products in the store.
func (c *Client) UpdateProducts(products []goshopify.Product) error {
	return updateProducts(c.Product, c.Variant, products)
}
