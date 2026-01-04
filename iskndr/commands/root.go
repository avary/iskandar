package commands

import "github.com/spf13/cobra"

func Execute() error {

	rootCmd := &cobra.Command{
		Use:   "iskndr",
		Short: "Iskandar is a CLI tool for exposing a tunnel to your local application.",
		Long:  "Long description",
		Args:  cobra.NoArgs,
	}

	rootCmd.AddCommand(newTunnelCommand())

	return rootCmd.Execute()
}
