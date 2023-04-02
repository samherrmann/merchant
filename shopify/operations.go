package shopify

import (
	"errors"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/memdb"
)

type ProductService = goshopify.ProductService
type VariantService = goshopify.VariantService
type Product = goshopify.Product
type Variant = goshopify.Variant
type ListOptions = goshopify.ListOptions

// getProducts gets all products from the store and stores them in the cache
// file.
func getProducts(pService ProductService, vService VariantService) ([]Product, error) {
	products, err := listProducts(pService, nil)
	if err != nil {
		return nil, err
	}
	count := len(products)
	for i, p := range products {
		fmt.Printf("Cloning metafields for product %v [%v/%v]\n", p.ID, i+1, count)
		if err := attachMetafields(pService, vService, &p); err != nil {
			return nil, err
		}
	}
	return products, err
}

func listProducts(service ProductService, options *ListOptions) ([]Product, error) {
	products := []Product{}
	defaultOptions := &ListOptions{
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
		productsPacket, pagination, err := service.ListWithPagination(options)
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

func updateProducts(
	pService ProductService,
	vService VariantService,
	products []Product,
) error {
	// Get latest inventory from live store so that we don't accidentally make
	// updates based on an outdated cache.
	inventory, err := getProducts(pService, vService)
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
		if _, err := pService.Create(p); err != nil {
			errs = append(errs, err)
		}
	}
	for _, p := range operations.ProductUpdates {
		if _, err := pService.Update(p); err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range operations.NewVariants {
		if _, err := vService.Create(v.ProductID, v); err != nil {
			errs = append(errs, err)
		}
	}
	for _, v := range operations.VariantUpdates {
		if _, err := vService.Update(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// attachMetafields fetches and attaches all metafields for the given product and its variants.
func attachMetafields(pService ProductService, vService VariantService, product *Product) error {
	metafields, err := pService.ListMetafields(product.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to get metafields for product %v: %w", product.ID, err)
	}
	product.Metafields = metafields

	for j, variant := range product.Variants {
		metafields, err := vService.ListMetafields(variant.ID, nil)
		if err != nil {
			return fmt.Errorf("failed to get metafields for variant %v: %w", variant.ID, err)
		}
		product.Variants[j].Metafields = metafields
	}
	return nil
}

func searchVariant(pService ProductService, vService VariantService, fn func(v *Variant) bool) (*Variant, error) {
	inventory, err := getProducts(pService, vService)
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

func getVariantCount(service ProductService) (int, error) {
	options := &ListOptions{
		Fields: "variants",
	}
	products, err := listProducts(service, options)
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
