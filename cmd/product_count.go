package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newProductCountCommand(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count total number of products and variants in the store",
		RunE: func(cmd *cobra.Command, args []string) error {
			productCount, err := shopClient.Product.Count(nil)
			if err != nil {
				return err
			}
			variantCount, err := shopClient.GetVariantCount()
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "Products: %v\n", productCount)
			fmt.Fprintf(w, "Variants: %v\n", variantCount)

			return nil
		},
	}
	return cmd
}
