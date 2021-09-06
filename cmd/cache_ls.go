package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files in cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cacheDir()
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\n", "FILE", "MODIFIED")
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%v\t%v\n", removeExtension(entry.Name()), info.ModTime())
		}
		w.Flush()
		return nil
	},
}

func removeExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
