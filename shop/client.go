package shop

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/config"
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

func (c *Client) GetInventoryWithMetafields() ([]goshopify.Product, error) {
	products, err := c.ListProducts(nil)
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
	products, err := c.ListProducts(options)
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

func (c *Client) ListProducts(options *goshopify.ListOptions) ([]goshopify.Product, error) {
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
