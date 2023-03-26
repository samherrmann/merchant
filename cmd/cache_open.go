package cmd

import (
	"path/filepath"

	"github.com/samherrmann/merchant/cache"
	"github.com/samherrmann/merchant/config"
	"github.com/samherrmann/merchant/editor"
	"github.com/spf13/cobra"
)

func newCacheOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open file from cache in text editor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				return err
			}
			dir, err := cache.Dir()
			if err != nil {
				return err
			}
			filename := filepath.Join(dir, args[0]+".json")
			return editor.New(cfg.TextEditor...).Open(filename)
		},
	}
}
