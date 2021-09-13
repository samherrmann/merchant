package cmd

import (
	"errors"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/csv"
	"github.com/samherrmann/shopctl/shop"
	"github.com/spf13/cobra"
)

var (
	errEmpty = errors.New("value and unit are empty")
)

func newPushProductsCommand(shopClient *shop.Client, metafieldDefs *config.MetafieldDefinitions) *cobra.Command {
	return &cobra.Command{
		Use:   "products <filename>",
		Short: "Update products in store with data from CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rows, err := csv.ReadFile(args[0] + ".csv")
			if err != nil {
				return err
			}
			for _, row := range rows {
				isProductMetafield := row.VariantID == 0
				isNewMetafield := row.MetafiledID == 0
				if isProductMetafield {
					if isNewMetafield {
						_, err := createMetafield(shopClient.Product, metafieldDefs.Product, &row)
						if err == errEmpty {
							continue
						}
						if err != nil {
							return fmt.Errorf(
								"cannot create metafield %q for product %v: %w",
								fmt.Sprintf("%v.%v", row.MetafieldNamespace, row.MetafieldKey),
								row.ProductID,
								err,
							)
						}
						continue
					}
					if _, err := updateMetafield(shopClient.Product, &row); err != nil {
						return fmt.Errorf("cannot update metafield %v for product %v: %w", row.MetafiledID, row.ProductID, err)
					}
					continue
				}
				if isNewMetafield {
					_, err := createMetafield(shopClient.Variant, metafieldDefs.Variant, &row)
					if err == errEmpty {
						continue
					}
					if err != nil {
						return fmt.Errorf(
							"cannot create metafield %q for variant %v: %w",
							fmt.Sprintf("%v.%v", row.MetafieldNamespace, row.MetafieldKey),
							row.VariantID,
							err,
						)
					}
					continue
				}
				if _, err := updateMetafield(shopClient.Variant, &row); err != nil {
					return fmt.Errorf("cannot update metafield %v for variant %v: %w", row.MetafiledID, row.VariantID, err)
				}
				continue
			}
			return nil
		},
	}
}

func createMetafield(service goshopify.MetafieldsService, definitions []config.MetafieldDefinition, row *csv.Row) (*goshopify.Metafield, error) {
	// Do nothing if metafield value and unit are empty.
	if row.MetafieldValue == "" && row.MetafieldUnit == "" {
		return nil, errEmpty
	}
	value, err := csv.ParseMetafieldValue(row)
	if err != nil {
		return nil, err
	}
	definition := config.FindMetafieldDefinition(definitions, row.MetafieldNamespace, row.MetafieldKey)
	if definition == nil {
		return nil, fmt.Errorf(
			"cannot find definition for metafield key %q",
			fmt.Sprintf("%v.%v", row.MetafieldNamespace, row.MetafieldKey),
		)
	}
	// Select product/variant ID based on service type.
	var id int64
	switch service.(type) {
	case *goshopify.ProductServiceOp:
		id = row.ProductID
	case *goshopify.VariantServiceOp:
		id = row.VariantID
	default:
		return nil, fmt.Errorf("unknown metafield service")
	}
	metafield := goshopify.Metafield{
		Namespace: row.MetafieldNamespace,
		Key:       row.MetafieldKey,
		ValueType: definition.Type,
		Value:     value,
	}
	return service.CreateMetafield(id, metafield)
}

func updateMetafield(service goshopify.MetafieldsService, row *csv.Row) (*goshopify.Metafield, error) {
	value, err := csv.ParseMetafieldValue(row)
	if err != nil {
		return nil, err
	}
	// Select product/variant ID based on service type.
	var id int64
	switch service.(type) {
	case *goshopify.ProductServiceOp:
		id = row.ProductID
	case *goshopify.VariantServiceOp:
		id = row.VariantID
	default:
		return nil, fmt.Errorf("unknown metafield service")
	}
	metafield := goshopify.Metafield{
		ID:    row.MetafiledID,
		Key:   row.MetafieldKey,
		Value: value,
	}
	return service.UpdateMetafield(id, metafield)
}
