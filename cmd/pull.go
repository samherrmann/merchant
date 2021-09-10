package cmd

import (
	"github.com/spf13/cobra"
)

func newPullCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Fetch resource data from the store",
	}
}
