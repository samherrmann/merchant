package memdb

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

func NewProductDB() *ProductDB {
	return &ProductDB{
		ids:     make(map[int64]goshopify.Product),
		titles:  make(map[string]goshopify.Product),
		handles: make(map[string]goshopify.Product),
	}
}

type ProductDB struct {
	ids     map[int64]goshopify.Product
	titles  map[string]goshopify.Product
	handles map[string]goshopify.Product
}

func (db *ProductDB) Add(p *goshopify.Product) error {
	if _, exists := db.ids[p.ID]; exists {
		return newInMemoryDBError("product id %q already exists", p.ID)
	}
	db.ids[p.ID] = *p
	if _, exists := db.titles[p.Title]; exists {
		return newInMemoryDBError("product title %q already exists", p.Title)
	}
	db.titles[p.Title] = *p
	if _, exists := db.titles[p.Handle]; exists {
		return newInMemoryDBError("product handle %q already exists", p.Handle)
	}
	db.handles[p.Handle] = *p
	return nil
}

func (db *ProductDB) Get(p *goshopify.Product) (*goshopify.Product, bool) {
	if p.ID != 0 {
		return db.GetByID(p.ID)
	}
	if p.Handle != "" {
		return db.GetByHandle(p.Title)
	}
	if p.Title != "" {
		return db.GetByTitle(p.Title)
	}
	return nil, false
}

func (db *ProductDB) GetByID(id int64) (*goshopify.Product, bool) {
	v, exists := db.ids[id]
	return &v, exists
}

func (db *ProductDB) GetByTitle(title string) (*goshopify.Product, bool) {
	v, exists := db.titles[title]
	return &v, exists
}

func (db *ProductDB) GetByHandle(handle string) (*goshopify.Product, bool) {
	v, exists := db.handles[handle]
	return &v, exists
}

// PatchID sets the ID on p if a matching product can be found in the database.
// ErrNotExists is returned if no match is found. The incoming product is not
// expected to be a complete product, but a partial update product.
func (db *ProductDB) PatchID(p *goshopify.Product) error {
	current, exists := db.Get(p)
	if !exists && p.ID != 0 {
		return fmt.Errorf("product ID %q does not exist", p.ID)
	}
	patchProductID(p, current.ID)
	return nil
}

// patchProductID sets the given ID on p and its variants.
func patchProductID(p *goshopify.Product, id int64) {
	p.ID = id
	for i := range p.Variants {
		p.Variants[i].ProductID = id
	}
}
