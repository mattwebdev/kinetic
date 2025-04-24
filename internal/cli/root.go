package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kinetic",
	Short: "Kinetic - Avalanche development environment",
	Long: `Kinetic is a powerful development environment for Avalanche and Subnet development.
It provides tools for managing local nodes, creating and deploying subnets,
and working with smart contracts.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/kinetic/config.yaml)")

	// Add commands
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(subnetCmd)
	rootCmd.AddCommand(contractCmd)
}
