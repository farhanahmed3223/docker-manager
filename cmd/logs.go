package cmd

import (
    "fmt"
    "os"

    "docker-manager/internal/docker"

    "github.com/spf13/cobra"
)

var tailLines int

var logsCmd = &cobra.Command{
    Use:   "logs [container]",
    Short: "Show container logs",
    Long:  `Display logs for a specific container.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        containerID := args[0]

        dockerClient, err := docker.NewDockerClient()
        if err != nil {
            fmt.Printf("Error connecting to Docker: %v\n", err)
            os.Exit(1)
        }

        logs, err := dockerClient.GetContainerLogs(containerID)
        if err != nil {
            fmt.Printf("Error getting logs: %v\n", err)
            os.Exit(1)
        }

        fmt.Println(logs)
    },
}

func init() {
    logsCmd.Flags().IntVarP(&tailLines, "tail", "t", 100, "Number of lines to show from the end of the logs")
}
