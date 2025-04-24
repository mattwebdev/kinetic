package node

import (
	"context"
	"testing"
	"time"

	"github.com/kinetic-dev/kinetic/internal/config"
)

func TestNewManager(t *testing.T) {
	cfg := config.DefaultConfig()
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create node manager: %v", err)
	}
	defer manager.Close()

	if manager == nil {
		t.Error("Expected non-nil manager")
	}
}

func TestNodeLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping node lifecycle test in short mode")
	}

	cfg := config.DefaultConfig()
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create node manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Test initial status
	status, err := manager.Status(ctx)
	if err != nil {
		t.Fatalf("Failed to get initial status: %v", err)
	}
	if status.IsRunning {
		t.Error("Node should not be running initially")
	}

	// Test starting node
	if err := manager.Start(ctx, cfg); err != nil {
		t.Fatalf("Failed to start node: %v", err)
	}

	// Wait for node to start
	time.Sleep(5 * time.Second)

	// Test status after starting
	status, err = manager.Status(ctx)
	if err != nil {
		t.Fatalf("Failed to get status after start: %v", err)
	}
	if !status.IsRunning {
		t.Error("Node should be running after start")
	}

	// Test health check
	health, err := manager.CheckHealth(ctx)
	if err != nil {
		t.Fatalf("Failed to check health: %v", err)
	}
	if !health.IsRunning {
		t.Error("Health check should show node as running")
	}

	// Test stopping node
	if err := manager.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop node: %v", err)
	}

	// Wait for node to stop
	time.Sleep(2 * time.Second)

	// Test status after stopping
	status, err = manager.Status(ctx)
	if err != nil {
		t.Fatalf("Failed to get status after stop: %v", err)
	}
	if status.IsRunning {
		t.Error("Node should not be running after stop")
	}
}

func TestWaitForHealthy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping wait for healthy test in short mode")
	}

	cfg := config.DefaultConfig()
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create node manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Start the node
	if err := manager.Start(ctx, cfg); err != nil {
		t.Fatalf("Failed to start node: %v", err)
	}
	defer manager.Stop(ctx)

	// Test waiting for healthy with timeout
	timeout := 2 * time.Minute
	if err := manager.WaitForHealthy(ctx, timeout); err != nil {
		t.Fatalf("Failed to wait for healthy: %v", err)
	}

	// Test health status after waiting
	health, err := manager.CheckHealth(ctx)
	if err != nil {
		t.Fatalf("Failed to check health after waiting: %v", err)
	}
	if !health.IsHealthy {
		t.Error("Node should be healthy after waiting")
	}
}

func TestHealthStatus(t *testing.T) {
	cfg := config.DefaultConfig()
	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("Failed to create node manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Test health status when node is not running
	health, err := manager.CheckHealth(ctx)
	if err != nil {
		t.Fatalf("Failed to check health: %v", err)
	}
	if health.IsRunning {
		t.Error("Health status should show node as not running")
	}
	if health.Error == "" {
		t.Error("Health status should contain error message when node is not running")
	}
}
