package system

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerClient wraps the Docker API client
type DockerClient struct {
	client *client.Client
}

// NewDockerClient creates a new Docker client
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &DockerClient{client: cli}, nil
}

// IsRunning checks if a container is running
func (d *DockerClient) IsRunning(ctx context.Context, containerName string) (bool, error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				return true, nil
			}
		}
	}
	return false, nil
}

// StartContainer starts a container by name
func (d *DockerClient) StartContainer(ctx context.Context, containerName string) error {
	err := d.client.ContainerStart(ctx, containerName, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container %s: %w", containerName, err)
	}
	return nil
}

// StopContainer stops a container by name
func (d *DockerClient) StopContainer(ctx context.Context, containerName string) error {
	err := d.client.ContainerStop(ctx, containerName, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}
	return nil
}

// PullImage pulls a Docker image
func (d *DockerClient) PullImage(ctx context.Context, image string) error {
	_, err := d.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", image, err)
	}
	return nil
}

// CreateContainer creates a new container
func (d *DockerClient) CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, name string) error {
	_, err := d.client.ContainerCreate(ctx, config, hostConfig, nil, nil, name)
	if err != nil {
		return fmt.Errorf("failed to create container %s: %w", name, err)
	}
	return nil
}

// Close closes the Docker client
func (d *DockerClient) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
