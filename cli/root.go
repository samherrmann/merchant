// Package cli defines the CLI commands.
package cli

import (
	"os"

	"github.com/samherrmann/merchant/config"
	"github.com/spf13/cobra"
)

func Execute() error {

	rootCmd := &cobra.Command{Use: config.AppName}

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
