package cmd

import (
	"github.com/samherrmann/goshopctl/config"
	"github.com/samherrmann/goshopctl/shop"
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
	productCmd := newProductCommand()
	productCmd.AddCommand(
		newProductPullCommand(shopClient, &c.MetafieldDefinitions),
		newProductPushCommand(shopClient, &c.MetafieldDefinitions),
	)
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		cacheCmd,
		configCmd,
		productCmd,
	)
	return rootCmd.Execute()
}
