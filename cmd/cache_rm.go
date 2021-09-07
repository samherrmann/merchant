package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove file from cache",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cacheDir()
		if err != nil {
			return err
		}
		return os.Remove(filepath.Join(dir, fmt.Sprintf("%v.json", args[0])))
	},
}
