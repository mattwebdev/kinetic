package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kinetic",
	Short: "Kinetic - Avalanche development toolkit",
	Long: `Kinetic is a development toolkit for building applications on Avalanche.
It provides tools for managing local nodes, deploying contracts, and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/kinetic/config.yaml)")

	// Add commands
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(contractCmd)
}
