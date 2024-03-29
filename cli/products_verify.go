package cli

import (
	"fmt"
	"io"

	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/memdb"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsVerifyCommand(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verifies the integrity of products and variants",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)

			inventory, err := store.GetProducts()
			if err != nil {
				return err
			}
			_, err = memdb.New(inventory)
			if err == nil {
				fmt.Fprintln(w, "Everything looks good!")
			}
			return err
		},
	}
	return cmd
}
