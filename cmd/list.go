package cmd

import (
    "fmt"
    "os"
    "text/tabwriter"

    "docker-manager/internal/docker"

    "github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List Docker containers",
    Long:  `List all Docker containers in a static table format.`,
    Run: func(cmd *cobra.Command, args []string) {
        dockerClient, err := docker.NewDockerClient()
        if err != nil {
            fmt.Printf("Error connecting to Docker: %v\n", err)
            os.Exit(1)
        }

        containers, err := dockerClient.ListContainers(listAll)
        if err != nil {
            fmt.Printf("Error listing containers: %v\n", err)
            os.Exit(1)
        }

        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
        fmt.Fprintln(w, "ID\tNAME\tIMAGE\tSTATUS\tPORTS\tCPU%\tMEMORY%\tNETWORK\tUPTIME")

        for _, c := range containers {
            uptime := time.Since(c.Created).Truncate(time.Second).String()
            status := c.Status
            if strings.Contains(status, "Up") {
                status = fmt.Sprintf("\033[32m%s\033[0m", status) // Green
            } else if strings.Contains(status, "Exited") {
                status = fmt.Sprintf("\033[31m%s\033[0m", status) // Red
            }

            fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%.1f\t%.1f\t%s\t%s\n",
                c.ID, c.Name, c.Image, status, c.Ports, c.CPU, c.Memory, c.Network, uptime)
        }
        w.Flush()
    },
}

func init() {
    listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all containers (default shows just running)")
}
