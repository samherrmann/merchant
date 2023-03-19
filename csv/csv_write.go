package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/collection"
)

func WriteProductFile(product *goshopify.Product) error {
	rows, err := makeRowsFromProducts([]goshopify.Product{*product})
	if err != nil {
		return err
	}
	return writeFile(fmt.Sprintf("%v.csv", product.ID), rows)
}

func WriteInventoryFile(products []goshopify.Product) error {
	rows, err := makeRowsFromProducts(products)
	if err != nil {
		return err
	}
	return writeFile("inventory.csv", rows)
}

func makeRowsFromProducts(products []goshopify.Product) ([][]string, error) {
	colIndexes := make(map[string]int)
	colIndexes[keyProductID] = len(colIndexes)
	colIndexes[keyVariantID] = len(colIndexes)
	colIndexes[keySKU] = len(colIndexes)
	colIndexes[keyBarcode] = len(colIndexes)
	colIndexes[keyTitle] = len(colIndexes)
	colIndexes[keyVendor] = len(colIndexes)
	colIndexes[keyProductType] = len(colIndexes)
	colIndexes[keyWeight] = len(colIndexes)
	colIndexes[keyWeightUnit] = len(colIndexes)
	colIndexes[keyPrice] = len(colIndexes)
	colIndexes[keyOption1Name] = len(colIndexes)
	colIndexes[keyOption1Value] = len(colIndexes)
	colIndexes[keyOption2Name] = len(colIndexes)
	colIndexes[keyOption2Value] = len(colIndexes)
	colIndexes[keyOption3Name] = len(colIndexes)
	colIndexes[keyOption3Value] = len(colIndexes)

	// Initialize rows with one row for the heading. We will come back at the end
	// to populate it with all the columns.
	rows := [][]string{{}}
	for _, p := range products {

		for _, v := range p.Variants {
			row := make([]string, len(colIndexes))
			weight, _ := v.Weight.Float64()
			price, _ := v.Price.Float64()
			row[colIndexes[keyProductID]] = fmt.Sprintf("%v", p.ID)
			row[colIndexes[keyVariantID]] = fmt.Sprintf("%v", v.ID)
			row[colIndexes[keySKU]] = v.Sku
			row[colIndexes[keyBarcode]] = v.Barcode
			row[colIndexes[keyTitle]] = p.Title
			row[colIndexes[keyVendor]] = p.Vendor
			row[colIndexes[keyProductType]] = p.ProductType
			row[colIndexes[keyWeight]] = fmt.Sprintf("%v", weight)
			row[colIndexes[keyWeightUnit]] = v.WeightUnit
			row[colIndexes[keyPrice]] = fmt.Sprintf("%v", price)

			if len(p.Options) > 0 {
				if p.Options[0].Name != "Title" {
					row[colIndexes[keyOption1Name]] = p.Options[0].Name
				}
				if v.Option1 != "Default Title" {
					row[colIndexes[keyOption1Value]] = v.Option1
				}
			}
			if len(p.Options) > 1 {
				row[colIndexes[keyOption2Name]] = p.Options[1].Name
				row[colIndexes[keyOption2Value]] = v.Option2
			}
			if len(p.Options) > 2 {
				row[colIndexes[keyOption3Name]] = p.Options[2].Name
				row[colIndexes[keyOption3Value]] = v.Option3
			}

			attachMetafield := func(metaType string, m goshopify.Metafield) {
				key := fmt.Sprintf("%s.metafields.%s.%s", metaType, m.Namespace, m.Key)
				index, exists := colIndexes[key]
				// If this is the first time encountering this metafield, then add it to
				// the colPositions map and grow the row slice.
				if !exists {
					index = len(colIndexes)
					colIndexes[key] = index
					row = append(row, "")
				}
				row[index] = fmt.Sprintf("%v", m.Value)
			}

			for _, m := range p.Metafields {
				attachMetafield("product", m)
			}
			for _, m := range v.Metafields {
				attachMetafield("variant", m)
			}
			rows = append(rows, row)
		}
	}
	// Populate the first row with all column names.
	rows[0] = make([]string, len(colIndexes))
	for k, v := range colIndexes {
		rows[0][v] = k
	}
	return padRows(rows), nil
}

func writeFile(filename string, rows [][]string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return csv.NewWriter(file).WriteAll(rows)
}

func padRows(rows [][]string) [][]string {
	if len(rows) == 0 {
		return rows
	}
	heading := rows[0]
	headingLength := len(heading)
	rowsLength := len(rows)
	for i := 1; i < rowsLength; i++ {
		rowLength := len(rows[i])
		if rowLength > headingLength {
			continue
		}
		rows[i] = collection.PadSliceRight(rows[i], headingLength)
	}
	return rows
}
