package node

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthStatus represents the node's health status
type HealthStatus struct {
	IsRunning      bool      `json:"is_running"`
	IsHealthy      bool      `json:"is_healthy"`
	IsBootstrapped bool      `json:"is_bootstrapped"`
	Version        string    `json:"version"`
	NetworkID      int       `json:"network_id"`
	LastChecked    time.Time `json:"last_checked"`
	Error          string    `json:"error,omitempty"`
}

// CheckHealth performs a health check on the node
func (m *NodeManager) CheckHealth(ctx context.Context) (*HealthStatus, error) {
	status := &HealthStatus{
		LastChecked: time.Now(),
	}

	// First check if the container is running
	nodeStatus, err := m.Status(ctx)
	if err != nil {
		return status, fmt.Errorf("failed to check container status: %w", err)
	}
	status.IsRunning = nodeStatus.IsRunning
	status.NetworkID = nodeStatus.NetworkID
	status.Version = nodeStatus.Version

	if !status.IsRunning {
		status.Error = "node is not running"
		return status, nil
	}

	// Check node health endpoint
	healthURL := fmt.Sprintf("http://localhost:%d/ext/health", m.cfg.Node.APIPort)
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		status.Error = fmt.Sprintf("failed to create request: %v", err)
		return status, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		status.Error = fmt.Sprintf("failed to connect to node: %v", err)
		return status, nil
	}
	defer resp.Body.Close()

	var healthResp struct {
		Healthy   bool   `json:"healthy"`
		Error     string `json:"error,omitempty"`
		Timestamp string `json:"timestamp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		status.Error = fmt.Sprintf("failed to decode response: %v", err)
		return status, nil
	}

	status.IsHealthy = healthResp.Healthy

	// Check if chains are bootstrapped
	infoURL := fmt.Sprintf("http://localhost:%d/ext/info", m.cfg.Node.APIPort)
	req, err = http.NewRequestWithContext(ctx, "GET", infoURL, nil)
	if err != nil {
		status.Error = fmt.Sprintf("failed to create info request: %v", err)
		return status, nil
	}

	resp, err = client.Do(req)
	if err != nil {
		status.Error = fmt.Sprintf("failed to get node info: %v", err)
		return status, nil
	}
	defer resp.Body.Close()

	var infoResp struct {
		NetworkID      int    `json:"networkID"`
		NodeVersion    string `json:"nodeVersion"`
		IsBootstrapped bool   `json:"isBootstrapped"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&infoResp); err != nil {
		status.Error = fmt.Sprintf("failed to decode info response: %v", err)
		return status, nil
	}

	status.IsBootstrapped = infoResp.IsBootstrapped

	return status, nil
}

// WaitForHealthy waits for the node to become healthy with a timeout
func (m *NodeManager) WaitForHealthy(ctx context.Context, timeout time.Duration) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for node to become healthy")
		case <-ticker.C:
			status, err := m.CheckHealth(ctx)
			if err != nil {
				return err
			}
			if status.IsHealthy && status.IsBootstrapped {
				return nil
			}
		}
	}
}
