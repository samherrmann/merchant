package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCountVariantsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "variants",
		Short: "Count total number of variants in the store",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := shopClient.GetVariantCount()
			if err != nil {
				return err
			}
			fmt.Println(count)
			return nil
		},
	}
	return cmd
}
