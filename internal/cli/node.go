package cli

import (
	"fmt"
	"time"

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
	RunE:  runNodeStart,
}

var nodeStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop local Avalanche node",
	RunE:  runNodeStop,
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check local node status",
	RunE:  runNodeStatus,
}

func runNodeStart(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cfg := config.Get()

	// Override config with flags if provided
	if apiPort, _ := cmd.Flags().GetInt("api-port"); apiPort != 0 {
		cfg.Node.APIPort = apiPort
	}
	if nodePort, _ := cmd.Flags().GetInt("node-port"); nodePort != 0 {
		cfg.Node.Port = nodePort
	}

	manager, err := node.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create node manager: %w", err)
	}
	defer manager.Close()

	if err := manager.Start(ctx, cfg); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	fmt.Println("Starting Avalanche node...")
	fmt.Printf("API endpoint: http://localhost:%d\n", cfg.Node.APIPort)

	// Wait for node to become healthy with a 2-minute timeout
	fmt.Println("Waiting for node to become healthy...")
	if err := manager.WaitForHealthy(ctx, 2*time.Minute); err != nil {
		return fmt.Errorf("node failed to become healthy: %w", err)
	}

	fmt.Println("Node is healthy and ready!")
	return nil
}

func runNodeStop(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cfg := config.Get()

	manager, err := node.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create node manager: %w", err)
	}
	defer manager.Close()

	if err := manager.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop node: %w", err)
	}

	fmt.Println("Node stopped successfully")
	return nil
}

func runNodeStatus(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	cfg := config.Get()

	manager, err := node.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create node manager: %w", err)
	}
	defer manager.Close()

	status, err := manager.CheckHealth(ctx)
	if err != nil {
		return fmt.Errorf("failed to check node status: %w", err)
	}

	fmt.Printf("Node Status:\n")
	fmt.Printf("  Running: %v\n", status.IsRunning)
	if status.IsRunning {
		fmt.Printf("  Healthy: %v\n", status.IsHealthy)
		fmt.Printf("  Bootstrapped: %v\n", status.IsBootstrapped)
		fmt.Printf("  Version: %s\n", status.Version)
		fmt.Printf("  Network ID: %d\n", status.NetworkID)
		fmt.Printf("  API Endpoint: http://localhost:%d\n", cfg.Node.APIPort)
		fmt.Printf("  Last Checked: %s\n", status.LastChecked.Format(time.RFC3339))
	}
	if status.Error != "" {
		fmt.Printf("  Error: %s\n", status.Error)
	}

	return nil
}

func init() {
	nodeCmd.AddCommand(nodeStartCmd)
	nodeCmd.AddCommand(nodeStopCmd)
	nodeCmd.AddCommand(nodeStatusCmd)

	// Add flags
	nodeStartCmd.Flags().IntP("node-port", "p", 9650, "Node port")
	nodeStartCmd.Flags().IntP("api-port", "a", 9651, "API port")
}
