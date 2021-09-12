package cmd

import (
	"github.com/samherrmann/shopctl/config"
	"github.com/samherrmann/shopctl/shop"
	"github.com/spf13/cobra"
)

func Execute() error {
	c, err := config.Load()
	if err != nil {
		return err
	}

	shopClient := shop.NewClient(c)

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
	pullCmd := newPullCommand()
	pullCmd.AddCommand(
		newPullProductCommand(shopClient, &c.MetafieldDefinitions),
		newPullProductsCommand(shopClient, &c.MetafieldDefinitions),
	)
	pushCmd := newPushCommand()
	pushCmd.AddCommand(
		newPushProductsCommand(shopClient, &c.MetafieldDefinitions),
	)
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		cacheCmd,
		configCmd,
		pullCmd,
		pushCmd,
		newVersionCommand(config.AppName),
	)
	return rootCmd.Execute()
}
