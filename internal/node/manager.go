package node

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/kinetic-dev/kinetic/internal/config"
	"github.com/kinetic-dev/kinetic/internal/system"
)

// Manager handles Avalanche node operations
type Manager struct {
	cfg    *config.Config
	docker *system.DockerClient
}

// NewManager creates a new node manager
func NewManager(cfg *config.Config) (*Manager, error) {
	docker, err := system.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Manager{
		cfg:    cfg,
		docker: docker,
	}, nil
}

// Start starts the Avalanche node
func (m *Manager) Start(ctx context.Context) error {
	// Check if node is already running
	running, err := m.docker.IsRunning(ctx, m.cfg.Docker.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to check node status: %w", err)
	}
	if running {
		return fmt.Errorf("node is already running")
	}

	// Ensure directories exist
	if err := system.EnsureDir(m.cfg.Node.DBDir); err != nil {
		return fmt.Errorf("failed to create DB directory: %w", err)
	}
	if err := system.EnsureDir(m.cfg.Node.LogDir); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	if err := system.EnsureDir(m.cfg.Node.StakingDir); err != nil {
		return fmt.Errorf("failed to create staking directory: %w", err)
	}

	// Pull the latest image
	if err := m.docker.PullImage(ctx, m.cfg.Docker.ImageTag); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Create container configuration
	containerConfig := &container.Config{
		Image: m.cfg.Docker.ImageTag,
		Cmd: []string{
			"--network-id=" + fmt.Sprint(m.cfg.Node.NetworkID),
			"--http-host=0.0.0.0",
			"--public-ip=127.0.0.1",
			"--db-dir=/root/.avalanchego/db",
			"--log-dir=/root/.avalanchego/logs",
			"--staking-enabled=false",
		},
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%d/tcp", m.cfg.Node.Port)):    struct{}{},
			nat.Port(fmt.Sprintf("%d/tcp", m.cfg.Node.APIPort)): struct{}{},
		},
	}

	// Create host configuration
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", m.cfg.Node.Port)): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", m.cfg.Node.Port)},
			},
			nat.Port(fmt.Sprintf("%d/tcp", m.cfg.Node.APIPort)): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", m.cfg.Node.APIPort)},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: m.cfg.Node.DBDir,
				Target: "/root/.avalanchego/db",
			},
			{
				Type:   mount.TypeBind,
				Source: m.cfg.Node.LogDir,
				Target: "/root/.avalanchego/logs",
			},
			{
				Type:   mount.TypeBind,
				Source: m.cfg.Node.StakingDir,
				Target: "/root/.avalanchego/staking",
			},
		},
	}

	// Create and start the container
	if err := m.docker.CreateContainer(ctx, containerConfig, hostConfig, m.cfg.Docker.ContainerName); err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := m.docker.StartContainer(ctx, m.cfg.Docker.ContainerName); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// Stop stops the Avalanche node
func (m *Manager) Stop(ctx context.Context) error {
	running, err := m.docker.IsRunning(ctx, m.cfg.Docker.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to check node status: %w", err)
	}
	if !running {
		return fmt.Errorf("node is not running")
	}

	if err := m.docker.StopContainer(ctx, m.cfg.Docker.ContainerName); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// Status returns the current node status
func (m *Manager) Status(ctx context.Context) (bool, error) {
	return m.docker.IsRunning(ctx, m.cfg.Docker.ContainerName)
}

// Close cleans up resources
func (m *Manager) Close() error {
	return m.docker.Close()
}
