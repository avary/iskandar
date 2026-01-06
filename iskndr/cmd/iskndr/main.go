package main

import (
	"os"

	"github.com/igneel64/iskandar/iskndr/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
