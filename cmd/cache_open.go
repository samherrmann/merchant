package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open file from cache in Visual Studio Code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cacheDir()
		if err != nil {
			return err
		}
		return runVSCode(filepath.Join(dir, fmt.Sprintf("%v.json", args[0])))
	},
}

func runVSCode(filename string) error {
	cmd := exec.Command("code", filename)
	return cmd.Run()
}
