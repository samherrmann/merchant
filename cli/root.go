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
		newCacheClearCommand(),
		newCacheDumpCommand(),
	)
	configCmd := newConfigCommand()
	configCmd.AddCommand(
		newConfigOpenCommand(),
	)
	productsCmd := newProductsCommand()
	productsCmd.AddCommand(
		newProductsCountCommand(os.Stdout),
		newProductsFakePushCommand(os.Stdout, config.AppName+".push.json"),
		newProductsPullCommand(),
		newProductsPushCommand(),
		newProductsVerifyCommand(os.Stdout),
	)
	rootCmd.AddCommand(
		cacheCmd,
		configCmd,
		productsCmd,
		newVersionCommand(config.AppName, config.Version),
	)
	return rootCmd.Execute()
}
