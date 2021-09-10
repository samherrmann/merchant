package cmd

import (
	"github.com/spf13/cobra"
)

func newPushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Update resource data in the store",
	}
}
