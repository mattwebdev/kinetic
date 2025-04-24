package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kinetic-dev/kinetic/internal/system"
	"github.com/spf13/viper"
)

var globalConfig *Config

// Config holds the application configuration
type Config struct {
	// Node configuration
	Node struct {
		Port       int    `mapstructure:"port"`
		APIPort    int    `mapstructure:"api_port"`
		NetworkID  int    `mapstructure:"network_id"`
		DBDir      string `mapstructure:"db_dir"`
		LogDir     string `mapstructure:"log_dir"`
		StakingDir string `mapstructure:"staking_dir"`
	} `mapstructure:"node"`

	// Docker configuration
	Docker struct {
		ImageTag      string `mapstructure:"image_tag"`
		ContainerName string `mapstructure:"container_name"`
	} `mapstructure:"docker"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{}

	// Node defaults
	cfg.Node.Port = 9650
	cfg.Node.APIPort = 9651
	cfg.Node.NetworkID = 12345 // Local network

	// Docker defaults
	cfg.Docker.ImageTag = "avaplatform/avalanchego:latest"
	cfg.Docker.ContainerName = "kinetic-node"

	return cfg
}

// Get returns the global config instance
func Get() *Config {
	if globalConfig == nil {
		// Return default config if not initialized
		return &Config{
			Node: struct {
				Port       int    `mapstructure:"port"`
				APIPort    int    `mapstructure:"api_port"`
				NetworkID  int    `mapstructure:"network_id"`
				DBDir      string `mapstructure:"db_dir"`
				LogDir     string `mapstructure:"log_dir"`
				StakingDir string `mapstructure:"staking_dir"`
			}{
				Port:       9650,
				APIPort:    9651,
				NetworkID:  12345,
				DBDir:      "data",
				LogDir:     "data",
				StakingDir: "data",
			},
		}
	}
	return globalConfig
}

// Load reads the config from the specified file
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.json"
	}

	// Create default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := &Config{
			Node: struct {
				Port       int    `mapstructure:"port"`
				APIPort    int    `mapstructure:"api_port"`
				NetworkID  int    `mapstructure:"network_id"`
				DBDir      string `mapstructure:"db_dir"`
				LogDir     string `mapstructure:"log_dir"`
				StakingDir string `mapstructure:"staking_dir"`
			}{
				Port:       9650,
				APIPort:    9651,
				NetworkID:  12345,
				DBDir:      "data",
				LogDir:     "data",
				StakingDir: "data",
			},
		}
		globalConfig = cfg
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure data directory is absolute
	if !filepath.IsAbs(cfg.Node.DBDir) {
		absPath, err := filepath.Abs(cfg.Node.DBDir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve absolute path for data directory: %w", err)
		}
		cfg.Node.DBDir = absPath
	}

	globalConfig = cfg
	return cfg, nil
}

// Save saves the current configuration to file
func (c *Config) Save() error {
	configDir, err := system.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join(configDir, "config.json"))

	// Convert config struct to map
	if err := v.MergeConfigMap(map[string]interface{}{
		"node": map[string]interface{}{
			"port":        c.Node.Port,
			"api_port":    c.Node.APIPort,
			"network_id":  c.Node.NetworkID,
			"db_dir":      c.Node.DBDir,
			"log_dir":     c.Node.LogDir,
			"staking_dir": c.Node.StakingDir,
		},
		"docker": map[string]interface{}{
			"image_tag":      c.Docker.ImageTag,
			"container_name": c.Docker.ContainerName,
		},
	}); err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
