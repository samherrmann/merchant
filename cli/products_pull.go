package cli

import (
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/editor"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsPullCommand() *cobra.Command {
	var skipCache *bool
	var openFile *bool
	var metafields *bool

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Fetch products and their metadata from the store",
		Long:  "Fetch all products from the store. Metafields may optionally be included",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)

			products, err := store.GetInventory(*skipCache, *metafields)
			if err != nil {
				return err
			}

			if err := csv.WriteInventoryFile(products); err != nil {
				return err
			}

			if *openFile {
				editor := newSpreadsheetEditor(cfg.SpreadsheetEditor...)
				if err := editor.Open(csv.InventoryFilename); err != nil {
					return err
				}
			}
			return nil
		},
	}
	openFile = cmd.Flags().Bool("open", false, "Open product file after pulling")
	skipCache = addCacheFlag(cmd)
	metafields = addMetafields(cmd)
	return cmd
}

func newSpreadsheetEditor(cmd ...string) editor.Editor {
	if len(cmd) == 0 {
		cmd = config.DefaultSpreadsheetEditor
	}
	return editor.New(cmd...)
}