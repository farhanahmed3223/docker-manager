package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	compactMode bool
)

var rootCmd = &cobra.Command{
	Use:   "docker-manager",
	Short: "A terminal-based Docker container manager",
	Long:  `A feature-rich TUI for managing Docker containers with real-time monitoring and controls.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&compactMode, "compact", "c", false, "Use compact display mode")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(interactiveCmd)
}
