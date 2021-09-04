package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files in cache.",
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
		fmt.Fprintf(w, "%v\t%v\n", "FILENAME", "MODIFIED")
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%v\t%v\n", entry.Name(), info.ModTime())
		}
		w.Flush()
		return nil
	},
}
