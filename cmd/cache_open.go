package cmd

import (
	"github.com/samherrmann/merchant/cache"
	"github.com/spf13/cobra"
)

func newCacheOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open file from cache in Visual Studio Code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cache.OpenFileInTextEditor(args[0] + ".json")
		},
	}
}
