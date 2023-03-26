package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(appName, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the %v version number", appName),
		Run: func(cmd *cobra.Command, args []string) {
			// Command usage is correct at this point.
			cmd.SilenceUsage = true

			fmt.Println(version)
		},
	}
}
