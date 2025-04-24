package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Check default values
	if cfg.Node.Port != 9650 {
		t.Errorf("expected default node port 9650, got %d", cfg.Node.Port)
	}
	if cfg.Node.APIPort != 9651 {
		t.Errorf("expected default API port 9651, got %d", cfg.Node.APIPort)
	}
	if cfg.Node.NetworkID != 12345 {
		t.Errorf("expected default network ID 12345, got %d", cfg.Node.NetworkID)
	}
	if cfg.Docker.ImageTag != "avaplatform/avalanchego:latest" {
		t.Errorf("expected default image tag 'avaplatform/avalanchego:latest', got %s", cfg.Docker.ImageTag)
	}
	if cfg.Docker.ContainerName != "kinetic-node" {
		t.Errorf("expected default container name 'kinetic-node', got %s", cfg.Docker.ContainerName)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test loading non-existent config
	cfg, err := Load("")
	if err != nil {
		t.Errorf("expected no error loading default config, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}

	// Create test config file
	configContent := `{
		"node": {
			"port": 9660,
			"api_port": 9661,
			"network_id": 54321,
			"db_dir": "custom/data",
			"log_dir": "custom/logs",
			"staking_dir": "custom/staking"
		},
		"docker": {
			"image_tag": "custom/avalanchego:v1.2.3",
			"container_name": "custom-node"
		}
	}`
	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loading custom config
	cfg, err = Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify custom values
	if cfg.Node.Port != 9660 {
		t.Errorf("expected custom node port 9660, got %d", cfg.Node.Port)
	}
	if cfg.Node.APIPort != 9661 {
		t.Errorf("expected custom API port 9661, got %d", cfg.Node.APIPort)
	}
	if cfg.Node.NetworkID != 54321 {
		t.Errorf("expected custom network ID 54321, got %d", cfg.Node.NetworkID)
	}
	if cfg.Docker.ImageTag != "custom/avalanchego:v1.2.3" {
		t.Errorf("expected custom image tag 'custom/avalanchego:v1.2.3', got %s", cfg.Docker.ImageTag)
	}
	if cfg.Docker.ContainerName != "custom-node" {
		t.Errorf("expected custom container name 'custom-node', got %s", cfg.Docker.ContainerName)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test config
	cfg := &Config{
		Node: struct {
			Port       int    `mapstructure:"port"`
			APIPort    int    `mapstructure:"api_port"`
			NetworkID  int    `mapstructure:"network_id"`
			DBDir      string `mapstructure:"db_dir"`
			LogDir     string `mapstructure:"log_dir"`
			StakingDir string `mapstructure:"staking_dir"`
		}{
			Port:       9670,
			APIPort:    9671,
			NetworkID:  98765,
			DBDir:      "test/data",
			LogDir:     "test/logs",
			StakingDir: "test/staking",
		},
		Docker: struct {
			ImageTag      string `mapstructure:"image_tag"`
			ContainerName string `mapstructure:"container_name"`
		}{
			ImageTag:      "test/avalanchego:latest",
			ContainerName: "test-node",
		},
	}

	// Save the config
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the saved config
	savedCfg, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Verify saved values
	if savedCfg.Node.Port != cfg.Node.Port {
		t.Errorf("expected saved node port %d, got %d", cfg.Node.Port, savedCfg.Node.Port)
	}
	if savedCfg.Node.APIPort != cfg.Node.APIPort {
		t.Errorf("expected saved API port %d, got %d", cfg.Node.APIPort, savedCfg.Node.APIPort)
	}
	if savedCfg.Docker.ImageTag != cfg.Docker.ImageTag {
		t.Errorf("expected saved image tag %s, got %s", cfg.Docker.ImageTag, savedCfg.Docker.ImageTag)
	}
}
