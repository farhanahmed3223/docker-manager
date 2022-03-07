package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"docker-manager/internal/docker"
	"github.com/spf13/cobra"
)

var (
	logsTail   int
	logsFollow bool
)

var logsCmd = &cobra.Command{
	Use:   "logs [container-name]",
	Short: "Show logs from a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dc, err := docker.NewDockerClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ctx := context.Background()
		reader, err := dc.GetLogs(ctx, args[0], logsTail, logsFollow)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer reader.Close()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	},
}

func init() {
	logsCmd.Flags().IntVarP(&logsTail, "tail", "n", 100, "Number of lines to show")
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	rootCmd.AddCommand(logsCmd)
}
