package cmd

import (
	"github.com/samherrmann/shopctl/config"
	"github.com/spf13/cobra"
)

func newConfigOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open configuration file in Visual Studio Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.OpenInTextEditor()
		},
	}
}
