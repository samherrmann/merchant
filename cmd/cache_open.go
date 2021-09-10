package cmd

import (
	"path/filepath"

	"github.com/samherrmann/goshopctl/cache"
	"github.com/samherrmann/goshopctl/utils"
	"github.com/spf13/cobra"
)

func newCacheOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open file from cache in Visual Studio Code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := cache.Dir()
			if err != nil {
				return err
			}
			filename := filepath.Join(dir, args[0])
			return utils.RunVSCode(filename + ".json")
		},
	}
}
