package cli

import "github.com/spf13/cobra"

func NewCLI() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(NewVisuazlize())

	return rootCmd
}
