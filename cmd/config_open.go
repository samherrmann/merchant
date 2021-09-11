package cmd

import (
	"path/filepath"

	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/utils"
	"github.com/spf13/cobra"
)

func newConfigOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open configuration file in Visual Studio Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := config.Dir()
			if err != nil {
				return err
			}
			filename := filepath.Join(dir, config.AppName) + ".json"
			return utils.RunVSCode(filename)
		},
	}
}
