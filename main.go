package main

import (
	"os"

	"github.com/samherrmann/merchant/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// No need to print the error because Cobra already does that for us.
		os.Exit(1)
	}
}
