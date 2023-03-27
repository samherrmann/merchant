package shopify

import (
	"errors"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/cache"
	"github.com/samherrmann/merchant/memdb"
)

type ProductService = goshopify.ProductService
type VariantService = goshopify.VariantService
type Product = goshopify.Product
type Variant = goshopify.Variant
type ListOptions = goshopify.ListOptions

// getInventory gets all products from the store and stores them in the cache
// file.
func getInventory(service ProductService, skipCache bool) ([]Product, error) {
	if skipCache {
		p, err := listProducts(service, nil)
		if err != nil {
			return nil, err
		}
		if err := cache.WriteInventoryFile(p); err != nil {
			return nil, err
		}
	}
	products, err := cache.ReadInventoryFile()
	// Fall back to pulling product from store if no cache exists, but only if
	// skipCache is false to prevent an infinite loop.
	if os.IsNotExist(err) && !skipCache {
		return getInventory(service, true)
	}
	return products, err
}

func getInventoryWithMetafields(pService ProductService, vService VariantService, skipCache bool) ([]Product, error) {
	products, err := getInventory(pService, skipCache)
	if err != nil {
		return nil, err
	}
	for i, product := range products {
		fmt.Printf("Getting metafields for product %v\n", product.ID)
		if err := attachMetafields(pService, vService, &product); err != nil {
			return nil, err
		}
		// TODO check if this is necessary.
		products[i] = product
	}
	if err := cache.WriteInventoryFile(products); err != nil {
		return nil, err
	}
	return products, nil
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
	inventory, err := getInventory(pService, true)
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

func searchVariant(service ProductService, fn func(v *Variant) bool) (*Variant, error) {
	inventory, err := getInventory(service, false)
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
