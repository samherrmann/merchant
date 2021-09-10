package cmd

import "github.com/spf13/cobra"

func newConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage configuration file",
	}
}
