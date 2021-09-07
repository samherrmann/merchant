package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configOpenCmd)
}

var configOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open config file in Visual Studio Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := configDir()
		if err != nil {
			return err
		}
		return runVSCode(filepath.Join(dir, fmt.Sprintf("%v.json", appName)))
	},
}
