package cmd

import (
	"context"
	"fmt"
	"os"

	"docker-manager/internal/docker"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all containers",
	Run: func(cmd *cobra.Command, args []string) {
		dc, err := docker.NewDockerClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		containers, err := dc.ListContainers(context.Background())
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		for _, c := range containers {
			fmt.Printf("%-20s %-15s %s\n", c.Names[0], c.Status, c.Image)
		}
	},
}
