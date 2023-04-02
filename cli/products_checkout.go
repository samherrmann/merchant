package cli

import (
	"github.com/samherrmann/merchant/cache"
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/editor"
	"github.com/spf13/cobra"
)

func newProductsCheckoutCommand() *cobra.Command {
	var openFile *bool

	cmd := &cobra.Command{
		Use:   "checkout",
		Short: "Creates a CSV file of all the products in the cache",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			c, err := cache.New()
			if err != nil {
				return err
			}

			products, err := c.Products().List()
			if err != nil {
				return err
			}

			if err := csv.WriteProductsFile(products); err != nil {
				return err
			}

			if *openFile {
				editor := newSpreadsheetEditor(cfg.SpreadsheetEditor...)
				if err := editor.Open(csv.ProductsFilename); err != nil {
					return err
				}
			}
			return nil
		},
	}
	openFile = cmd.Flags().Bool("open", false, "Open product file after pulling")
	return cmd
}

func newSpreadsheetEditor(cmd ...string) editor.Editor {
	if len(cmd) == 0 {
		cmd = config.DefaultSpreadsheetEditor
	}
	return editor.New(cmd...)
}
