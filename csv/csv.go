// Package csv enables reading and writing product data to and from a CSV file.
package csv

import (
	"encoding/json"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/jszwec/csvutil"
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/shop"
)

type Row struct {
	ProductID          int64       `csv:"product_id"`
	ProductTitle       string      `csv:"product_title"`
	VariantID          int64       `csv:"variant_id,omitempty"`
	MetafiledID        int64       `csv:"metafiled_id,omitempty"`
	MetafieldKey       string      `csv:"metafield_key"`
	MetafieldNamespace string      `csv:"metafield_namespace"`
	MetafieldValue     interface{} `csv:"metafield_value"`
	MetafieldUnit      string      `csv:"metafield_unit,omitempty"`
}

type MetafieldKey = string

type RowsMap = map[MetafieldKey]Row

func ReadFile(filename string) ([]Row, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	rows := []Row{}
	if err := csvutil.Unmarshal(bytes, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func WriteFile(filename string, rows []Row) error {
	bytes, err := csvutil.Marshal(rows)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}

func WriteProductFile(product *goshopify.Product, metafieldDefs *config.MetafieldDefinitions) error {
	rows, err := makeRowsFromProducts([]goshopify.Product{*product}, metafieldDefs)
	if err != nil {
		return err
	}
	return WriteFile(fmt.Sprintf("%v.csv", product.ID), rows)
}

func WriteInventoryFile(products []goshopify.Product, metafieldDefs *config.MetafieldDefinitions) error {
	rows, err := makeRowsFromProducts(products, metafieldDefs)
	if err != nil {
		return err
	}
	return WriteFile("inventory.csv", rows)
}

func ParseMetafieldValue(row *Row) (interface{}, error) {
	if row.MetafieldUnit == "" {
		return row.MetafieldValue, nil
	}
	num, ok := row.MetafieldValue.(float64)
	if !ok {
		return nil, fmt.Errorf("metafield value %q cannot be converted to float64", row.MetafieldValue)
	}
	return shop.Measurement{Value: num, Unit: row.MetafieldUnit}, nil
}

func makeRowsFromProducts(products []goshopify.Product, metafieldDefs *config.MetafieldDefinitions) ([]Row, error) {
	rows := []Row{}
	for _, p := range products {
		productRows := makeRowsFromDefinitions(p.ID, p.Title, 0, metafieldDefs.Product)
		for _, m := range p.Metafields {
			row, err := makeRowFromMetafield(p.ID, p.Title, 0, &m)
			if err != nil {
				return nil, err
			}
			productRows[row.MetafieldKey] = *row
		}
		for _, row := range productRows {
			rows = append(rows, row)
		}

		for _, v := range p.Variants {
			variantRows := makeRowsFromDefinitions(p.ID, p.Title, v.ID, metafieldDefs.Variant)
			for _, m := range v.Metafields {
				row, err := makeRowFromMetafield(p.ID, p.Title, v.ID, &m)
				if err != nil {
					return nil, err
				}
				variantRows[row.MetafieldKey] = *row
			}
			for _, row := range variantRows {
				rows = append(rows, row)
			}
		}
	}
	return rows, nil
}

func makeRowsFromDefinitions(productID int64, productTitle string, variantID int64, definitions []config.MetafieldDefinition) RowsMap {
	rowsMap := RowsMap{}
	for _, def := range definitions {
		row := Row{
			ProductID:          productID,
			ProductTitle:       productTitle,
			VariantID:          variantID,
			MetafieldKey:       def.Key,
			MetafieldNamespace: def.Namespace,
		}
		rowsMap[def.Key] = row
	}
	return rowsMap
}

func makeRowFromMetafield(productID int64, productTitle string, variantID int64, metafield *goshopify.Metafield) (*Row, error) {
	row := &Row{
		ProductID:          productID,
		ProductTitle:       productTitle,
		VariantID:          variantID,
		MetafiledID:        metafield.ID,
		MetafieldKey:       metafield.Key,
		MetafieldNamespace: metafield.Namespace,
		MetafieldValue:     metafield.Value,
	}

	if metafield.ValueType == "json_string" {
		measurement := &shop.Measurement{}
		bytes := []byte(fmt.Sprint(metafield.Value))
		if err := json.Unmarshal(bytes, measurement); err != nil {
			return nil, err
		}
		row.MetafieldValue = measurement.Value
		row.MetafieldUnit = measurement.Unit
	}
	return row, nil
}
