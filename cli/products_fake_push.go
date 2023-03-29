package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/memdb"
	"github.com/samherrmann/merchant/shopify"
	"github.com/spf13/cobra"
)

func newProductsFakePushCommand(output io.Writer, filename string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fake-push <filename>",
		Short: "Print the data that the push command would send to the store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true
			filename := args[0]

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			store := shopify.NewClient(&cfg.Store)
			incoming, err := csv.ReadProducts(filename)
			if err != nil {
				return err
			}
			inventory, err := store.GetProducts(false)
			if err != nil {
				return err
			}
			db, err := memdb.New(inventory)
			if err != nil {
				return err
			}
			operations, err := db.Operations(incoming)
			if err != nil {
				return err
			}
			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()
			if err := operations.PrintJSON(file); err != nil {
				return err
			}
			if err := operations.PrintSummary(output); err != nil {
				return err
			}
			_, err = fmt.Fprintf(output, "\nSee file %q for details\n", filename)
			return err
		},
	}
	return cmd
}
