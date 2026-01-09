package commands

import "github.com/spf13/cobra"

func Execute() error {

	rootCmd := &cobra.Command{
		Use:   "iskndr",
		Short: "Iskandar is a CLI tool for exposing a tunnel to your local application.",
		Long:  "A lightweight, self-hosted HTTP tunnel service for exposing local applications to the internet in Go.",
		Args:  cobra.NoArgs,
	}

	rootCmd.AddCommand(newTunnelCommand())
	rootCmd.AddCommand(newVersionCommand())

	return rootCmd.Execute()
}
