package cli

import (
	"fmt"
	"io"

	units "github.com/docker/go-units"
	"github.com/samherrmann/merchant/cache"
	"github.com/spf13/cobra"
)

func newCacheSizeCommand(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "size",
		Short: "Prints the size of the cache database",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			size, err := cache.Size()
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(out, units.HumanSize(float64(size)))
			return err
		},
	}
}
