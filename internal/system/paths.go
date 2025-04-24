package system

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetUserConfigDir returns the path to the user's configuration directory
func GetUserConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	kineticDir := filepath.Join(configDir, "kinetic")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(kineticDir, 0755); err != nil {
		return "", err
	}

	return kineticDir, nil
}

// GetDataDir returns the path to the data directory
func GetDataDir() (string, error) {
	var baseDir string
	var err error

	switch runtime.GOOS {
	case "windows":
		baseDir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(baseDir, "AppData", "Local", "Kinetic")
	case "darwin":
		baseDir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(baseDir, "Library", "Application Support", "Kinetic")
	default: // linux and others
		baseDir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(baseDir, ".local", "share", "kinetic")
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	return baseDir, nil
}

// GetNodeDataDir returns the path to the Avalanche node data directory
func GetNodeDataDir() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	nodeDir := filepath.Join(dataDir, "node")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(nodeDir, 0755); err != nil {
		return "", err
	}

	return nodeDir, nil
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
