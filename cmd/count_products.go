package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCountProductsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Count total number of products in the store",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := shopClient.Product.Count(nil)
			if err != nil {
				return err
			}
			fmt.Println(count)
			return nil
		},
	}
	return cmd
}
