package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Ports   string
	State   string
	Created int64
	CPUPerc float64
	MemPerc float64
}

type StatsResult struct {
	CPUPercent float64
	MemUsage   uint64
	MemLimit   uint64
	MemPercent float64
	NetRx      uint64
	NetTx      uint64
}

type DockerClient struct {
	cli *client.Client
	ctx context.Context
	cancel context.CancelFunc
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &DockerClient{cli: cli, ctx: ctx, cancel: cancel}, nil
}

func (d *DockerClient) Close() {
	d.cancel()
	d.cli.Close()
}

func (d *DockerClient) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	if ctx == nil {
		ctx = d.ctx
	}
	raw, err := d.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}
	out := make([]ContainerInfo, 0, len(raw))
	for _, c := range raw {
		name := "unknown"
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		out = append(out, ContainerInfo{
			ID:      c.ID[:12],
			Name:    name,
			Image:   c.Image,
			Status:  c.Status,
			State:   c.State,
			Created: c.Created,
		})
	}
	return out, nil
}

func (d *DockerClient) StartContainer(ctx context.Context, id string) error {
	return d.cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

func (d *DockerClient) StopContainer(ctx context.Context, id string) error {
	timeout := 10 * time.Second
	return d.cli.ContainerStop(ctx, id, &timeout)
}

func (d *DockerClient) GetLogs(ctx context.Context, nameOrID string, tail int, follow bool) (io.ReadCloser, error) {
	tailStr := "all"
	if tail > 0 {
		tailStr = strconv.Itoa(tail)
	}
	return d.cli.ContainerLogs(ctx, nameOrID, types.ContainerLogsOptions{
		ShowStdout: true, ShowStderr: true,
		Follow: follow, Tail: tailStr,
	})
}

func (d *DockerClient) GetStats(ctx context.Context, nameOrID string) (*StatsResult, error) {
	resp, err := d.cli.ContainerStats(ctx, nameOrID, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw types.StatsJSON
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage - raw.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPUStats.SystemUsage - raw.PreCPUStats.SystemUsage)
	numCPU := float64(len(raw.CPUStats.CPUUsage.PercpuUsage))
	cpuPct := 0.0
	if sysDelta > 0 && cpuDelta > 0 {
		cpuPct = (cpuDelta / sysDelta) * numCPU * 100.0
	}
	memUsage := raw.MemoryStats.Usage
	memLimit := raw.MemoryStats.Limit
	memPct := 0.0
	if memLimit > 0 {
		memPct = float64(memUsage) / float64(memLimit) * 100.0
	}
	return &StatsResult{CPUPercent: cpuPct, MemUsage: memUsage, MemLimit: memLimit, MemPercent: memPct}, nil
}
