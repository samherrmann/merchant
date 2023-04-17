package cli

import (
	"fmt"
	"io"
	"os"

	cachepkg "github.com/samherrmann/merchant/cache"
	"github.com/samherrmann/merchant/csv"
	"github.com/samherrmann/merchant/memdb"
	"github.com/spf13/cobra"
)

func newProductsFakePushCommand(output io.Writer, outputFilename string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fake-push <filename>",
		Short: "Print the data that the push command would send to the store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true
			inputFilename := args[0]

			incoming, err := csv.ReadProducts(inputFilename)
			if err != nil {
				return err
			}
			cache, err := cachepkg.New()
			if err != nil {
				return err
			}
			inventory, err := cache.Products().List()
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
			file, err := os.Create(outputFilename)
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
			_, err = fmt.Fprintf(output, "\nSee file %q for details\n", outputFilename)
			return err
		},
	}
	return cmd
}
