package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage local Avalanche node",
	Long:  `Commands for managing your local Avalanche node instance.`,
}

var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local Avalanche node",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting local node...")
		// TODO: Implement node start logic
		return nil
	},
}

var nodeStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop local Avalanche node",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Stopping local node...")
		// TODO: Implement node stop logic
		return nil
	},
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check local node status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking node status...")
		// TODO: Implement status check logic
		return nil
	},
}

func init() {
	nodeCmd.AddCommand(nodeStartCmd)
	nodeCmd.AddCommand(nodeStopCmd)
	nodeCmd.AddCommand(nodeStatusCmd)

	// Add flags
	nodeStartCmd.Flags().IntP("port", "p", 9650, "Node port")
	nodeStartCmd.Flags().IntP("api-port", "a", 9651, "API port")
}
