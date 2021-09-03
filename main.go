package main

import (
	"fmt"
	"os"

	"github.com/samherrmann/goshopctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// TODO remove
// func sampleProductMetafieldCreate(client *goshopify.Client) (*goshopify.Metafield, error) {
// 	return client.Product.CreateMetafield(6573170753578, goshopify.Metafield{
// 		Key:       "box_per_carton",
// 		Value:     123,
// 		ValueType: "integer",
// 		Namespace: "common",
// 	})
// }
