package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(csvCmd)
}

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Generates a CSV file for one or all products.",
	RunE: func(cmd *cobra.Command, args []string) error {
		products := []goshopify.Product{}
		if len(args) == 0 {
			var err error
			products, err = readProductsFile()
			if err != nil {
				return err
			}
		} else {
			productID, err := parseID(args[0])
			if err != nil {
				return err
			}
			product, err := readProductFile(productID)
			if err != nil {
				return err
			}
			products = append(products, *product)
		}
		csvRows, err := convertProductsToCSVRows(products)
		if err != nil {
			return err
		}
		return writeCSVFile(csvRows)
	},
}

func readProductFile(id int64) (*goshopify.Product, error) {
	bytes, err := readCacheFile(fmt.Sprintf("%v.json", id))
	if err != nil {
		return nil, err
	}
	products := &goshopify.Product{}
	if err = json.Unmarshal(bytes, products); err != nil {
		return nil, err
	}
	return products, nil
}

func readProductsFile() ([]goshopify.Product, error) {
	bytes, err := readCacheFile(cacheFilename)
	if err != nil {
		return nil, err
	}
	products := []goshopify.Product{}
	if err = json.Unmarshal(bytes, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func convertProductsToCSVRows(products []goshopify.Product) ([]CSVRow, error) {
	rows := []CSVRow{}
	for _, p := range products {
		for _, m := range p.Metafields {
			row, err := convertMetafieldToCSVRow(p.ID, 0, m)
			if err != nil {
				return nil, err
			}
			rows = append(rows, *row)
		}

		for _, v := range p.Variants {
			for _, m := range p.Metafields {
				row, err := convertMetafieldToCSVRow(p.ID, v.ID, m)
				if err != nil {
					return nil, err
				}
				rows = append(rows, *row)
			}
		}
	}
	return rows, nil
}

func convertMetafieldToCSVRow(productID int64, variantID int64, metafield goshopify.Metafield) (*CSVRow, error) {
	row := &CSVRow{
		ProductID:      productID,
		VariantID:      variantID,
		MetafiledID:    metafield.ID,
		MetafieldKey:   metafield.Key,
		MetafieldValue: metafield.Value,
	}

	if metafield.ValueType == "json_string" {
		measurement := &Measurement{}
		if err := json.Unmarshal([]byte(fmt.Sprint(metafield.Value)), measurement); err != nil {
			return nil, fmt.Errorf("error unmarshaling metafield JSON string: %w", err)
		}
		row.MetafieldValue = measurement.Value
		row.MetafieldUnit = measurement.Unit
	}
	return row, nil
}

func writeCSVFile(rows []CSVRow) error {
	bytes, err := csvutil.Marshal(rows)
	if err != nil {
		return fmt.Errorf("error encoding to CSV: %w", err)
	}
	return ioutil.WriteFile(csvFilename, bytes, 0644)
}

type CSVRow struct {
	ProductID      int64       `csv:"product_id"`
	VariantID      int64       `csv:"variant_id,omitempty"`
	MetafiledID    int64       `csv:"metafiled_id"`
	MetafieldKey   string      `csv:"metafield_key"`
	MetafieldValue interface{} `csv:"metafield_value"`
	MetafieldUnit  string      `csv:"metafield_unit,omitempty"`
}

type Measurement struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
