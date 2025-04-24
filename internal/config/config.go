package config

import (
	"fmt"
	"path/filepath"

	"github.com/kinetic-dev/kinetic/internal/system"
	"github.com/spf13/viper"
)

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

// Load loads the configuration from file
func Load() (*Config, error) {
	configDir, err := system.GetUserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Set up Viper
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	// Load defaults
	cfg := DefaultConfig()

	// Set up data directories
	dataDir, err := system.GetDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	cfg.Node.DBDir = filepath.Join(dataDir, "db")
	cfg.Node.LogDir = filepath.Join(dataDir, "logs")
	cfg.Node.StakingDir = filepath.Join(dataDir, "staking")

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, create it with defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			configFile := filepath.Join(configDir, "config.yaml")
			if err := v.SafeWriteConfigAs(configFile); err != nil {
				return nil, fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// Save saves the current configuration to file
func (c *Config) Save() error {
	configDir, err := system.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join(configDir, "config.yaml"))

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
