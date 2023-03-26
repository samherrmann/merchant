package cmd

import (
	"fmt"
	"io"

	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/memdb"
	"github.com/samherrmann/merchant/shop"
	"github.com/spf13/cobra"
)

func newProductVerifyCommand(w io.Writer) *cobra.Command {
	var skipCache *bool
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
			store := shop.NewClient(&cfg.Store)

			inventory, err := store.GetInventory(*skipCache)
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
	skipCache = cmd.Flags().Bool("skip-cache", false, "Pull directly from store even if a local copy exists in the cache")
	return cmd
}
