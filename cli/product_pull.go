package cli

import (
	"strconv"

	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/editor"
	"github.com/samherrmann/merchant/shopify"
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
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)
			arg := args[0]
			if arg == "inventory" {
				products, err := store.GetInventory(*skipCache)
				if err != nil {
					return err
				}
				csv.WriteInventoryFile(products)
			} else {
				productID, err := strconv.ParseInt(arg, 10, 64)
				if err != nil {
					return err
				}
				product, err := store.GetProduct(productID, *skipCache)
				if err != nil {
					return err
				}
				if err := csv.WriteProductFile(product); err != nil {
					return err
				}
			}
			if *openFile {
				editor := newSpreadsheetEditor(cfg.SpreadsheetEditor...)
				if err := editor.Open(arg + ".csv"); err != nil {
					return err
				}
			}
			return nil
		},
	}
	openFile = cmd.Flags().Bool("open", false, "Open product file after pulling")
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}

func newSpreadsheetEditor(cmd ...string) editor.Editor {
	if len(cmd) == 0 {
		cmd = config.DefaultSpreadsheetEditor
	}
	return editor.New(cmd...)
}
