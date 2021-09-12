package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time with -ldflags.
	Version = "dev"
)

func newVersionCommand(appName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print version number of %v", appName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
}
