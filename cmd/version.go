package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(appName, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the %v version number", appName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
