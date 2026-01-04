package main

import (
	"os"

	"github.com/igneel64/iskndr/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
