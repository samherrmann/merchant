package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/memdb"
	"github.com/spf13/cobra"
)

func newProductFakePushCommand(output io.Writer, filename string) *cobra.Command {
	var skipCache *bool
	cmd := &cobra.Command{
		Use:   "fake-push <filename>",
		Short: "Print the data that the push command would send to the store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			incoming, err := csv.ReadProducts(args[0])
			if err != nil {
				return err
			}
			inventory, err := shopClient.GetInventory(*skipCache)
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
	skipCache = cmd.Flags().Bool(
		"skip-cache",
		false,
		"Pull from store even if a local copy exists in the cache",
	)
	return cmd
}
