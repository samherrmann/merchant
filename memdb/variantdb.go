package memdb

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
)

const (
	defaultOptionValue = "Default Title"
)

func NewVariantDB(pdb *ProductDB) *VariantDB {
	return &VariantDB{
		products: pdb,
		ids:      make(map[int64]goshopify.Variant),
		skus:     make(map[string]goshopify.Variant),
		barcodes: make(map[string]goshopify.Variant),
		options:  make(map[string]goshopify.Variant),
	}
}

type VariantDB struct {
	products *ProductDB
	ids      map[int64]goshopify.Variant
	options  map[string]goshopify.Variant
	skus     map[string]goshopify.Variant
	barcodes map[string]goshopify.Variant
}

func (db *VariantDB) Add(v *goshopify.Variant) error {
	if _, exists := db.ids[v.ID]; exists {
		return newInMemoryDBError("variant id %q already exists", v.ID)
	}
	db.ids[v.ID] = *v
	optionsKey := encodeOptions(v.ProductID, v.Option1, v.Option2, v.Option3)
	if _, exists := db.options[optionsKey]; exists {
		return newInMemoryDBError("sku %q already exists", v.Sku)
	}
	db.options[optionsKey] = *v
	if v.Sku != "" {
		if _, exists := db.skus[v.Sku]; exists {
			return newInMemoryDBError("sku %q already exists", v.Sku)
		}
		db.skus[v.Sku] = *v
	}
	if v.Barcode != "" {
		if _, exists := db.barcodes[v.Barcode]; exists {
			return newInMemoryDBError("barcode %q already exists", v.Barcode)
		}
		db.barcodes[v.Barcode] = *v
	}
	return nil
}

func (db *VariantDB) GetByID(id int64) (*goshopify.Variant, bool) {
	v, exists := db.ids[id]
	return &v, exists
}

func (db *VariantDB) GetBySku(sku string) (*goshopify.Variant, bool) {
	v, exists := db.skus[sku]
	return &v, exists
}

func (db *VariantDB) GetByBarcode(barcode string) (*goshopify.Variant, bool) {
	v, exists := db.barcodes[barcode]
	return &v, exists
}

func (db *VariantDB) GetByOptions(productID int64, option1 string, option2, option3 string) (*goshopify.Variant, bool) {
	v, exists := db.options[encodeOptions(productID, option1, option2, option3)]
	return &v, exists
}

func (db *VariantDB) Get(v *goshopify.Variant) (*goshopify.Variant, error) {
	// Find match by ID:
	if v.ID != 0 {
		dbv, exists := db.GetByID(v.ID)
		if !exists {
			return nil, fmt.Errorf("variant ID %q does not exist", v.ID)
		}
		// Aside: If we have a matching variant ID, then we allow the SKU and/or
		// barcode to be updated, i.e. we do not validate the SKU and barcode
		// against existing values here.
		if err := validateProductID(dbv, v); err != nil {
			return nil, err
		}
		return dbv, nil
	}

	// Find match by barcode:
	if v.Barcode != "" {
		if dbv, exists := db.GetByBarcode(v.Barcode); exists {
			if err := validateProductID(dbv, v); err != nil {
				return nil, err
			}
			if !equalNonEmptyStrings(dbv.Sku, v.Sku) {
				return nil, fmt.Errorf(
					"SKU mismatch for barcode %v: current = %q, incoming = %q",
					v.Barcode,
					dbv.Sku,
					v.Sku,
				)
			}
			return dbv, nil
		}
		// New barcode, no matching variant found.
	}

	// Find match by SKU:
	if v.Sku != "" {
		if dbv, exists := db.GetBySku(v.Sku); exists {
			if err := validateProductID(dbv, v); err != nil {
				return nil, err
			}
			// The following validation is redundant because from above we already
			// know that if "incoming" has a non-empty barcode that it's a new
			// barcode.
			if !equalNonEmptyStrings(dbv.Barcode, v.Barcode) {
				return nil, fmt.Errorf(
					"barcode mismatch for SKU %v: current = %q, incoming = %q",
					v.Sku,
					dbv.Barcode,
					v.Barcode,
				)
			}
			return dbv, nil
		}
		// New SKU, no matching variant found.
	}

	// Find match by options:
	//
	// Options are only unique within a product, but not globally. Therefore, to
	// find a variant by options it must contain a product ID.
	if v.ProductID != 0 {
		dbp, exists := db.products.GetByID(v.ProductID)
		if !exists {
			return nil, fmt.Errorf("product ID %q does not exist", v.ProductID)
		}
		for i := range dbp.Variants {
			dbv := dbp.Variants[i]
			if optionEqual(v.Option1, dbv.Option1) &&
				optionEqual(v.Option2, dbv.Option2) &&
				optionEqual(v.Option3, dbv.Option3) {
				return &dbv, nil
			}
		}
	}
	return nil, ErrNotExist
}

// PatchID sets the ID on v if a matching variant can be found in the database.
// ErrNotExists is returned if no match is found. The incoming variant is not
// expected to be a complete variant, but a partial update variant.
func (db *VariantDB) PatchID(v *goshopify.Variant) error {
	dbv, err := db.Get(v)
	if err != nil {
		return err
	}
	return copyVariantIDs(dbv, v)
}

func validateProductID(current *goshopify.Variant, incoming *goshopify.Variant) error {
	if current.ID != incoming.ID {
		return fmt.Errorf(
			"variant ID mismatch: current = %v, incoming = %v",
			current.ID,
			incoming.ID,
		)
	}
	// If a product ID is provided together with a variant ID, then it must
	// match the current product ID.
	if incoming.ProductID != 0 && incoming.ProductID != current.ProductID {
		return fmt.Errorf(
			"product ID mismatch for variant ID %q: current = %q, incoming = %q",
			incoming.ID,
			current.ProductID,
			incoming.ProductID,
		)
	}
	return nil
}

// equalNonEmptyStrings returns true if either s1 or s2 are empty or if they are
// equal to each other. Returns false otherwise.
func equalNonEmptyStrings(s1 string, s2 string) bool {
	return !(s2 != "" && s1 != "" && s2 != s1)
}

// optionEqual returns true of o1 is equal to o2. The default option value is
// considered equal to an empty option value.
func optionEqual(o1 string, o2 string) bool {
	if o1 == defaultOptionValue {
		o1 = ""
	}
	if o2 == defaultOptionValue {
		o2 = ""
	}
	return o1 == o2
}

// copyVariantIDs copies the variant ID and ProductID from src to dst.
func copyVariantIDs(src, dst *goshopify.Variant) error {
	if dst.ID != 0 && dst.ID != src.ID {
		return fmt.Errorf(
			"variant ID already set: dst = %v, src = %v",
			dst.ID,
			src.ID,
		)
	}
	if dst.ProductID != 0 && dst.ProductID != src.ProductID {
		return fmt.Errorf(
			"product ID on variant %v already set: dst = %v, src = %v",
			src.ID,
			dst.ProductID,
			src.ProductID,
		)
	}
	dst.ID = src.ID
	dst.ProductID = src.ProductID
	return nil
}
