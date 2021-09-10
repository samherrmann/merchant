package cmd

import (
	"os"
	"path/filepath"

	"github.com/samherrmann/goshopctl/cache"
	"github.com/spf13/cobra"
)

func newCacheRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "Remove file from cache",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := cache.Dir()
			if err != nil {
				return err
			}
			return os.Remove(filepath.Join(dir, args[0]) + ".json")
		},
	}
}
