package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push <filename>",
	Short: "Update products in the store with data in CSV file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rows, err := readCSVFile(fmt.Sprintf("%v.csv", args[0]))
		if err != nil {
			return err
		}
		for _, row := range rows {
			isProductMetafield := row.VariantID == 0
			isNewMetafield := row.MetafiledID == 0
			if isProductMetafield {
				if isNewMetafield {
					if _, err := createMetafield(shopClient.Product, &row); err != nil {
						return fmt.Errorf("cannot create metafield for product %v: %w", row.ProductID, err)
					}
					break
				}
				if _, err := updateMetafield(shopClient.Product, &row); err != nil {
					return fmt.Errorf("cannot update metafield %v for product %v: %w", row.MetafiledID, row.ProductID, err)
				}
				break
			}
			if isNewMetafield {
				if _, err := createMetafield(shopClient.Variant, &row); err != nil {
					return fmt.Errorf("cannot create metafield for variant %v: %w", row.VariantID, err)
				}
				break
			}
			if _, err := updateMetafield(shopClient.Variant, &row); err != nil {
				return fmt.Errorf("cannot update metafield %v for variant %v: %w", row.MetafiledID, row.VariantID, err)
			}
			break
		}
		return nil
	},
}

func readCSVFile(filename string) ([]CSVRow, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	rows := []CSVRow{}
	if err := csvutil.Unmarshal(bytes, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func createProductMetafield(row *CSVRow) (*goshopify.Metafield, error) {
	value, err := marshalValue(row)
	if err != nil {
		return nil, err
	}
	return shopClient.Product.CreateMetafield(row.ProductID, goshopify.Metafield{
		Key:   row.MetafieldKey,
		Value: value,
	})
}

func createMetafield(service goshopify.MetafieldsService, row *CSVRow) (*goshopify.Metafield, error) {
	value, err := marshalValue(row)
	if err != nil {
		return nil, err
	}
	return service.UpdateMetafield(row.VariantID, goshopify.Metafield{
		Key:   row.MetafieldKey,
		Value: value,
	})
}

func updateMetafield(service goshopify.MetafieldsService, row *CSVRow) (*goshopify.Metafield, error) {
	value, err := marshalValue(row)
	if err != nil {
		return nil, err
	}
	return service.UpdateMetafield(row.VariantID, goshopify.Metafield{
		ID:    row.MetafiledID,
		Key:   row.MetafieldKey,
		Value: value,
	})
}

func marshalMeasurement(value float64, unit string) ([]byte, error) {
	return json.Marshal(&Measurement{Value: value, Unit: unit})
}

func marshalValue(row *CSVRow) (interface{}, error) {
	value := row.MetafieldValue
	if row.MetafieldUnit != "" {
		num, ok := row.MetafieldValue.(float64)
		if !ok {
			return nil, fmt.Errorf("metafield value %q cannot be converted to a number", row.MetafieldValue)
		}
		bytes, err := marshalMeasurement(num, row.MetafieldUnit)
		if err != nil {
			return nil, err
		}
		value = string(bytes)
	}
	return value, nil
}
