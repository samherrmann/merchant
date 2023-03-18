package cmd

import (
	"fmt"

	"github.com/samherrmann/shopctl/config"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print version number of %v", config.AppName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Version)
		},
	}
}
