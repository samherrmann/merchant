package cmd

import (
	"github.com/spf13/cobra"
)

func newCountCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "count",
		Short: "Count number of items for a resource",
	}
}
