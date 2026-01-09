package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  "Print the version information for iskndr",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("iskndr version %s\n", Version)
		},
	}
}
