// Package cmd defines the CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/exec"
	"github.com/samherrmann/shopctl/shop"
	"github.com/spf13/cobra"
)

var (
	shopClient *shop.Client
)

func Execute() error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	rootCmd := &cobra.Command{Use: config.AppName}
	storeName := rootCmd.PersistentFlags().String("store", "", "Name of the store")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		storeConfig := c.Stores.Get(*storeName)
		if storeConfig == nil {
			return fmt.Errorf("no config for store %q", *storeName)
		}
		shopClient = shop.NewClient(storeConfig)
		return nil
	}

	if len(c.TextEditor) > 0 {
		exec.TextEditorCmd = c.TextEditor
	}

	if len(c.SpreadsheetEditor) > 0 {
		exec.SpreadsheetEditorCmd = c.SpreadsheetEditor
	}

	cacheCmd := newCacheCommand()
	cacheCmd.AddCommand(
		newCacheListCommand(),
		newCacheOpenCommand(),
		newCacheRemoveCommand(),
	)
	configCmd := newConfigCommand()
	configCmd.AddCommand(
		newConfigOpenCommand(),
	)
	productCmd := newProductCommand()
	productCmd.AddCommand(
		newProductCountCommand(os.Stdout),
		newProductFakePushCommand(os.Stdout, config.AppName+".push.json"),
		newProductPullCommand(),
		newProductPushCommand(),
		newProductVerifyCommand(os.Stdout),
	)
	rootCmd.AddCommand(
		cacheCmd,
		configCmd,
		productCmd,
		newVersionCommand(config.AppName, config.Version),
	)
	return rootCmd.Execute()
}
