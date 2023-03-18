// Package memdb is a basic in-memory database that indexes products and
// variants by several properties.
package memdb

import (
	"errors"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

var (
	ErrNotExist = errors.New("does not exist")
)

// New returns a new in-memory database.
func New(products []goshopify.Product) (*MemoryDB, error) {
	pdb := NewProductDB()
	db := &MemoryDB{
		products: pdb,
		variants: NewVariantDB(pdb),
	}

	errs := []error{}
	for i := range products {
		p := products[i]
		if err := db.products.Add(&p); err != nil {
			errs = append(errs, err)
		}
		for j := range p.Variants {
			v := p.Variants[j]
			if err := db.variants.Add(&v); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return db, errors.Join(errs...)
}

type MemoryDB struct {
	products *ProductDB
	variants *VariantDB
}

func (db *MemoryDB) Variants() *VariantDB {
	return db.variants
}

func (db *MemoryDB) Products() *ProductDB {
	return db.products
}

// Operations groups the given changes by the type of operation needed to apply
// them to the database.
func (db *MemoryDB) Operations(changes []goshopify.Product) (*Operations, error) {
	operations := &Operations{}
outerLoop:
	for i := range changes {
		p := changes[i]
		err := db.Products().PatchID(&p)
		if err != nil && !errors.Is(err, ErrNotExist) {
			return nil, err
		}
		if p.ID != 0 {
			operations.UpdateProduct(p)
		}
		for i := range p.Variants {
			v := p.Variants[i]
			err := db.Variants().PatchID(&v)
			if err != nil && !errors.Is(err, ErrNotExist) {
				return nil, err
			}
			if v.ProductID == 0 {
				operations.CreateProduct(p)
				continue outerLoop
			}
			if v.ID == 0 {
				operations.CreateVariant(v)
				continue
			}
			operations.UpdateVariant(v)
		}
	}
	return operations, nil
}

func newInMemoryDBError(format string, a ...any) error {
	return fmt.Errorf("in-memory database: %v", fmt.Errorf(format, a...))
}

func encodeOptions(productID int64, option1 string, option2 string, option3 string) string {
	return fmt.Sprintf("%v/%v/%v/%v", productID, option1, option2, option3)
}
