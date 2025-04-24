package cli

import (
	"context"
	"fmt"

	"github.com/kinetic-dev/kinetic/internal/config"
	"github.com/kinetic-dev/kinetic/internal/node"
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

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override config with flags if provided
		if port, _ := cmd.Flags().GetInt("port"); port != 0 {
			cfg.Node.Port = port
		}
		if apiPort, _ := cmd.Flags().GetInt("api-port"); apiPort != 0 {
			cfg.Node.APIPort = apiPort
		}

		// Create node manager
		manager, err := node.NewManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create node manager: %w", err)
		}
		defer manager.Close()

		// Start the node
		if err := manager.Start(context.Background()); err != nil {
			return fmt.Errorf("failed to start node: %w", err)
		}

		fmt.Println("Node started successfully!")
		fmt.Printf("API endpoint: http://localhost:%d\n", cfg.Node.APIPort)
		return nil
	},
}

var nodeStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop local Avalanche node",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Stopping local node...")

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create node manager
		manager, err := node.NewManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create node manager: %w", err)
		}
		defer manager.Close()

		// Stop the node
		if err := manager.Stop(context.Background()); err != nil {
			return fmt.Errorf("failed to stop node: %w", err)
		}

		fmt.Println("Node stopped successfully!")
		return nil
	},
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check local node status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking node status...")

		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create node manager
		manager, err := node.NewManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create node manager: %w", err)
		}
		defer manager.Close()

		// Get node status
		running, err := manager.Status(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get node status: %w", err)
		}

		if running {
			fmt.Printf("Node is running (API endpoint: http://localhost:%d)\n", cfg.Node.APIPort)
		} else {
			fmt.Println("Node is not running")
		}
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
