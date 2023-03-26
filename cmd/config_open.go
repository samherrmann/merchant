package cmd

import (
	"os"

	"github.com/samherrmann/merchant/config"
	"github.com/spf13/cobra"
)

func newConfigOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open configuration file in Visual Studio Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			cfg, err := config.Load()
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
				cfg, err = config.InitFile()
				if err != nil {
					return err
				}
			}
			return cfg.OpenInTextEditor()
		},
	}
}
