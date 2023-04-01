package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samherrmann/merchant/cache"
	"github.com/spf13/cobra"
)

func newCacheDumpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dump",
		Short: "Exports the cache database to a file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			c, err := cache.New()
			if err != nil {
				return fmt.Errorf("cache: %w", err)
			}

			products, err := c.Products().List()
			if err != nil {
				return fmt.Errorf("getting products from cache: %w", err)
			}

			b, err := json.MarshalIndent(products, "", "  ")
			if err != nil {
				return fmt.Errorf("json marshal products: %w", err)
			}

			filename := fmt.Sprintf("%s.cache.products.json", cache.AppName)

			if err := os.WriteFile(filename, b, 0644); err != nil {
				return fmt.Errorf("writing cache to file: %w", err)
			}
			return nil
		},
	}
}
