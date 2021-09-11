package shop

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/config"
)

type Client struct {
	*goshopify.Client
}

func NewClient(c *config.Config) *Client {
	return &Client{
		Client: goshopify.NewClient(
			goshopify.App{
				ApiKey:   c.APIKey,
				Password: c.Password,
			},
			c.ShopName,
			"",
			goshopify.WithRetry(3),
		),
	}
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

func (c *Client) GetProductsWithMetafields() ([]goshopify.Product, error) {
	products := []goshopify.Product{}
	options := &goshopify.ListOptions{
		// 250 is the maximum limit
		// https://shopify.dev/api/admin/rest/reference/products/product?api%5Bversion%5D=2020-10#endpoints-2020-10
		Limit: 250,
	}
	for {
		productsPacket, pagination, err := c.Product.ListWithPagination(options)
		if err != nil {
			return nil, fmt.Errorf("failed to get packet of products: %w", err)
		}

		for i, product := range productsPacket {
			fmt.Printf("Getting metafields for product %v\n", product.ID)
			if err := c.attachMetafields(&product); err != nil {
				return nil, err
			}
			// TODO check if this is necessary.
			productsPacket[i] = product
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
