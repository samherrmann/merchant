package cmd

import (
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/goshopctl/config"
	"github.com/samherrmann/goshopctl/csv"
	"github.com/samherrmann/goshopctl/shop"
	"github.com/spf13/cobra"
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
						if _, err := createMetafield(shopClient.Product, metafieldDefs.Product, &row); err != nil {
							return fmt.Errorf("cannot create metafield for product %v: %w", row.ProductID, err)
						}
						continue
					}
					if _, err := updateMetafield(shopClient.Product, &row); err != nil {
						return fmt.Errorf("cannot update metafield %v for product %v: %w", row.MetafiledID, row.ProductID, err)
					}
					continue
				}
				if isNewMetafield {
					if _, err := createMetafield(shopClient.Variant, metafieldDefs.Variant, &row); err != nil {
						return fmt.Errorf("cannot create metafield for variant %v: %w", row.VariantID, err)
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
	value, err := csv.ParseMetafieldValue(row)
	if err != nil {
		return nil, err
	}
	definition := config.FindMetafieldDefinition(definitions, row.MetafieldNamespace, row.MetafieldKey)
	if definition == nil {
		return nil, fmt.Errorf("cannot find definition for metafiled key %v", row.MetafieldKey)
	}
	return service.CreateMetafield(row.VariantID, goshopify.Metafield{
		Key:       row.MetafieldKey,
		Value:     value,
		ValueType: definition.Type,
		Namespace: definition.Namespace,
	})
}

func updateMetafield(service goshopify.MetafieldsService, row *csv.Row) (*goshopify.Metafield, error) {
	value, err := csv.ParseMetafieldValue(row)
	if err != nil {
		return nil, err
	}
	return service.UpdateMetafield(row.VariantID, goshopify.Metafield{
		ID:    row.MetafiledID,
		Key:   row.MetafieldKey,
		Value: value,
	})
}
