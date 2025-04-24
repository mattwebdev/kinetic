package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetUserConfigDir(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get config directory
	configDir, err := GetUserConfigDir()
	if err != nil {
		t.Fatalf("Failed to get user config dir: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory was not created")
	}

	// Check if directory is in the correct location
	if !filepath.IsAbs(configDir) {
		t.Error("Config directory path is not absolute")
	}
}

func TestGetDataDir(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "data-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get data directory
	dataDir, err := GetDataDir()
	if err != nil {
		t.Fatalf("Failed to get data dir: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Error("Data directory was not created")
	}

	// Check if directory is in the correct location based on OS
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(dataDir, "AppData\\Local\\Kinetic") {
			t.Error("Windows data directory not in correct location")
		}
	case "darwin":
		if !strings.Contains(dataDir, "Library/Application Support/Kinetic") {
			t.Error("macOS data directory not in correct location")
		}
	default:
		if !strings.Contains(dataDir, ".local/share/kinetic") {
			t.Error("Linux data directory not in correct location")
		}
	}
}

func TestGetNodeDataDir(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "node-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get node data directory
	nodeDir, err := GetNodeDataDir()
	if err != nil {
		t.Fatalf("Failed to get node data dir: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(nodeDir); os.IsNotExist(err) {
		t.Error("Node data directory was not created")
	}

	// Check if directory is named correctly
	if !strings.HasSuffix(nodeDir, "node") {
		t.Error("Node directory does not have correct name")
	}
}

func TestEnsureDir(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ensure-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create new directory",
			path:    filepath.Join(tmpDir, "new-dir"),
			wantErr: false,
		},
		{
			name:    "create nested directories",
			path:    filepath.Join(tmpDir, "parent", "child", "grandchild"),
			wantErr: false,
		},
		{
			name:    "ensure existing directory",
			path:    tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if directory exists
			if _, err := os.Stat(tt.path); os.IsNotExist(err) {
				t.Error("Directory was not created")
			}
		})
	}
}

func TestDockerClient(t *testing.T) {
	// Create a new Docker client
	client, err := NewDockerClient()
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}
	defer client.Close()

	// Test if client is not nil
	if client == nil {
		t.Error("Docker client is nil")
	}

	// Test if client has valid connection
	if client.client == nil {
		t.Error("Docker client connection is nil")
	}
}
