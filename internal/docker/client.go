package docker

import (
    "context"
    "fmt"
    "io"
    "time"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/client"
    "github.com/shirou/gopsutil/v3/process"
)

type DockerClient struct {
    cli *client.Client
}

type ContainerInfo struct {
    ID      string
    Name    string
    Image   string
    Status  string
    State   string
    Ports   string
    Created time.Time
    CPU     float64
    Memory  float64
    Network string
}

func NewDockerClient() (*DockerClient, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, fmt.Errorf("failed to create Docker client: %w", err)
    }
    return &DockerClient{cli: cli}, nil
}

func (d *DockerClient) ListContainers(all bool) ([]ContainerInfo, error) {
    ctx := context.Background()
    containers, err := d.cli.ContainerList(ctx, types.ContainerListOptions{
        All: all,
    })
    if err != nil {
        return nil, err
    }

    var result []ContainerInfo
    for _, c := range containers {
        info := ContainerInfo{
            ID:      c.ID[:12],
            Name:    c.Names[0][1:], // Remove leading slash
            Image:   c.Image,
            Status:  c.Status,
            State:   c.State,
            Ports:   formatPorts(c.Ports),
            Created: time.Unix(c.Created, 0),
        }

        // Get detailed stats
        stats, err := d.getContainerStats(c.ID)
        if err == nil {
            info.CPU = stats.CPU
            info.Memory = stats.Memory
            info.Network = stats.Network
        }

        result = append(result, info)
    }
    return result, nil
}

func (d *DockerClient) GetContainerStats(containerID string) (*ContainerStats, error) {
    return d.getContainerStats(containerID)
}

type ContainerStats struct {
    CPU     float64
    Memory  float64
    Network string
}

func (d *DockerClient) getContainerStats(containerID string) (*ContainerStats, error) {
    ctx := context.Background()
    stats, err := d.cli.ContainerStats(ctx, containerID, false)
    if err != nil {
        return nil, err
    }
    defer stats.Body.Close()

    var v types.StatsJSON
    if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
        return nil, err
    }

    // Calculate CPU percentage
    cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
    systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
    cpuPercent := 0.0
    if systemDelta > 0 {
        cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100
    }

    // Calculate memory usage
    memUsage := float64(v.MemoryStats.Usage) / 1024 / 1024 // MB
    memLimit := float64(v.MemoryStats.Limit) / 1024 / 1024 // MB
    memPercent := 0.0
    if memLimit > 0 {
        memPercent = (memUsage / memLimit) * 100
    }

    // Network stats
    networkRx := float64(v.Networks["eth0"].RxBytes) / 1024 / 1024
    networkTx := float64(v.Networks["eth0"].TxBytes) / 1024 / 1024
    network := fmt.Sprintf("↓%.1fM/↑%.1fM", networkRx, networkTx)

    return &ContainerStats{
        CPU:     cpuPercent,
        Memory:  memPercent,
        Network: network,
    }, nil
}

func (d *DockerClient) StartContainer(containerID string) error {
    ctx := context.Background()
    return d.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (d *DockerClient) StopContainer(containerID string) error {
    ctx := context.Background()
    return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (d *DockerClient) RestartContainer(containerID string) error {
    ctx := context.Background()
    return d.cli.ContainerRestart(ctx, containerID, container.StopOptions{})
}

func (d *DockerClient) RemoveContainer(containerID string) error {
    ctx := context.Background()
    return d.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
}

func (d *DockerClient) GetContainerLogs(containerID string) (string, error) {
    ctx := context.Background()
    out, err := d.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Tail:       "100",
        Follow:     false,
    })
    if err != nil {
        return "", err
    }
    defer out.Close()

    logBytes, err := io.ReadAll(out)
    if err != nil {
        return "", err
    }

    return string(logBytes), nil
}

func (d *DockerClient) StreamLogs(containerID string) (io.ReadCloser, error) {
    ctx := context.Background()
    return d.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Follow:     true,
        Since:      time.Now().Format(time.RFC3339),
    })
}

func formatPorts(ports []types.Port) string {
    if len(ports) == 0 {
        return ""
    }
    var result string
    for _, p := range ports {
        if p.PublicPort > 0 {
            result += fmt.Sprintf("%d->%d/%s ", p.PublicPort, p.PrivatePort, p.Type)
        }
    }
    return result
}
