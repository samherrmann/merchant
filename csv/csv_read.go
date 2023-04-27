package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/collection"
	"github.com/samherrmann/merchant/csv/metafields"
	"github.com/shopspring/decimal"
)

func ReadProducts(filename string) ([]goshopify.Product, error) {
	rows, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	return groupVariants(rows)
}

func readFile(filename string) ([][]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return csv.NewReader(file).ReadAll()
}

// groupVariants groups variants that have the same title into the same product.
// The first row is expected to be the header.
func groupVariants(rows [][]string) ([]goshopify.Product, error) {
	products := collection.NewOrderedMap[string, goshopify.Product]()

	if len(rows) < 2 {
		return products.Slice(), nil
	}

	header := rows[0]
	titleColIndex := collection.IndexOf(header, keyTitle)
	if titleColIndex < 0 {
		return nil, fmt.Errorf("no %q column found", keyTitle)
	}

	rowsLength := len(rows)
	for i := 1; i < rowsLength; i++ {
		row := collection.PadSliceRight(rows[i], len(header))
		title := row[titleColIndex]
		if title == "" {
			return nil, fmt.Errorf("title in row %v can not be empty", i)
		}
		product, exists := products.Get(title)
		if !exists {
			product = goshopify.Product{}
		}
		if err := attachRecordToProduct(&product, header, row); err != nil {
			return nil, fmt.Errorf("row %v: %w", i, err)
		}
		products.Set(product.Title, product)
	}
	return products.Slice(), nil
}

// attachRecordToProduct attaches a record (row) from the CSV file to the given
// product.
func attachRecordToProduct(product *goshopify.Product, header []string, record []string) error {
	variant := &goshopify.Variant{}
	for i, v := range record {
		colName := header[i]
		switch colName {
		case keyProductID:
			id, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return colError(colName, err)
			}
			variant.ProductID = id
			product.ID = id
		case keyVariantID:
			id, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return colError(colName, err)
			}
			variant.ID = id
		case keySKU:
			variant.Sku = v
		case keyBarcode:
			variant.Barcode = v
		case keyTitle:
			product.Title = v
		case keyVendor:
			product.Vendor = v
		case keyProductType:
			product.ProductType = v
		case keyWeight:
			dec, err := parseDecimal(v)
			if err != nil {
				return colError(colName, err)
			}
			variant.Weight = dec
		case keyWeightUnit:
			variant.WeightUnit = v
		case keyPrice:
			dec, err := parseDecimal(v)
			if err != nil {
				return colError(colName, err)
			}
			variant.Price = dec
		case keyOption1Name:
			attachOptionToProduct(product, 0, v)
		case keyOption2Name:
			attachOptionToProduct(product, 1, v)
		case keyOption3Name:
			attachOptionToProduct(product, 2, v)
		case keyOption1Value:
			variant.Option1 = v
		case keyOption2Value:
			variant.Option2 = v
		case keyOption3Value:
			variant.Option3 = v
		default:
			if strings.HasPrefix(colName, "product.metafields") {
				m, err := metafields.New(colName)
				if err != nil {
					return err
				}
				m.Value = v
				product.Metafields = metafields.Attach(product.Metafields, m)
				break
			}
			if strings.HasPrefix(colName, "variant.metafields") {
				m, err := metafields.New(colName)
				if err != nil {
					return err
				}
				m.Value = v
				variant.Metafields = metafields.Attach(variant.Metafields, m)
				break
			}
			return fmt.Errorf("unknown column %q", colName)
		}
	}
	product.Variants = append(product.Variants, *variant)
	return nil
}

func attachOptionToProduct(p *goshopify.Product, index int, name string) {
	if name != "" {
		p.Options = collection.PadSliceRight(p.Options, index+1)
		p.Options[index].Name = name
	}
}

func colError(colName string, err error) error {
	return fmt.Errorf("column %q: %w", colName, err)
}

func parseDecimal(s string) (*decimal.Decimal, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	d := decimal.NewFromFloat(v)
	return &d, nil
}
