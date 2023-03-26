package main

import (
	"os"

	"github.com/samherrmann/merchant/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		// No need to print the error because Cobra already does that for us.
		os.Exit(1)
	}
}
