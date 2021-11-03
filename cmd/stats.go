package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"
    "time"

    "docker-manager/internal/docker"

    "github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
    Use:   "stats",
    Short: "Show real-time container statistics",
    Long:  `Display real-time CPU, memory, and network statistics for all containers.`,
    Run: func(cmd *cobra.Command, args []string) {
        dockerClient, err := docker.NewDockerClient()
        if err != nil {
            fmt.Printf("Error connecting to Docker: %v\n", err)
            os.Exit(1)
        }

        ticker := time.NewTicker(2 * time.Second)
        defer ticker.Stop()

        for {
            containers, err := dockerClient.ListContainers(true)
            if err != nil {
                fmt.Printf("Error listing containers: %v\n", err)
                os.Exit(1)
            }

            // Clear screen and move cursor to top
            fmt.Print("\033[H\033[2J")

            w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
            fmt.Fprintln(w, "ID\tNAME\tCPU%\tMEMORY%\tNETWORK I/O\tSTATUS")

            totalCPU, totalMemory := 0.0, 0.0
            for _, c := range containers {
                cpuStyle, memStyle := "", ""
                if c.CPU > 80 {
                    cpuStyle = "\033[31m" // Red
                } else if c.CPU > 60 {
                    cpuStyle = "\033[33m" // Yellow
                } else {
                    cpuStyle = "\033[32m" // Green
                }

                if c.Memory > 80 {
                    memStyle = "\033[31m" // Red
                } else if c.Memory > 60 {
                    memStyle = "\033[33m" // Yellow
                } else {
                    memStyle = "\033[32m" // Green
                }

                fmt.Fprintf(w, "%s\t%s\t%s%.1f%%\033[0m\t%s%.1f%%\033[0m\t%s\t%s\n",
                    c.ID[:12], c.Name, cpuStyle, c.CPU, memStyle, c.Memory, c.Network, c.Status)

                totalCPU += c.CPU
                totalMemory += c.Memory
            }

            fmt.Fprintf(w, "\n%s\t%s\t%.1f%%\t%.1f%%\t%s\t%s\n",
                "TOTAL", "", totalCPU, totalMemory, "", fmt.Sprintf("%d containers", len(containers)))

            w.Flush()
            fmt.Printf("\nRefreshing every 2 seconds. Press Ctrl+C to stop...")

            <-ticker.C
        }
    },
}
