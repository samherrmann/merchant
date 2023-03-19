package cmd

import (
	"strconv"

	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/exec"
	"github.com/spf13/cobra"
)

func newProductPullCommand() *cobra.Command {
	var skipCache *bool
	var openFile *bool
	cmd := &cobra.Command{
		Use:   "pull <id>|inventory",
		Short: "Fetch product and its metadata from the store",
		Long:  "Fetch a single product or the entire product inventory from the store. The product metafields are included.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]
			if arg == "inventory" {
				products, err := shopClient.GetInventory(*skipCache)
				if err != nil {
					return err
				}
				csv.WriteInventoryFile(products)
				return nil
			}
			productID, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return err
			}
			product, err := shopClient.GetProduct(productID, *skipCache)
			if err != nil {
				return err
			}
			return csv.WriteProductFile(product)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if *openFile {
				return exec.RunSpreadsheetEditor(args[0] + ".csv")
			}
			return nil
		},
	}
	openFile = cmd.Flags().Bool("open", false, "Open product file after pulling")
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}
