package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"docker-manager/internal/docker"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [container-name]",
	Short: "Show live CPU/memory stats for a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dc, err := docker.NewDockerClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ctx := context.Background()
		stats, err := dc.GetStats(ctx, args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("Container : %s\n", args[0])
		fmt.Printf("CPU %%     : %.2f\n", stats.CPUPercent)
		fmt.Printf("Memory    : %.1f MiB / %.1f MiB (%.1f%%)\n",
			float64(stats.MemUsage)/1024/1024,
			float64(stats.MemLimit)/1024/1024,
			stats.MemPercent,
		)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
