package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

type DockerClient struct {
	cli *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %w", err)
	}
	return &DockerClient{cli: cli}, nil
}

func (d *DockerClient) ListContainers(ctx context.Context) ([]types.Container, error) {
	return d.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
}
