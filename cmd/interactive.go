package cmd

import (
	"fmt"
	"os"

	"docker-manager/internal/docker"
	"docker-manager/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"tui", "i"},
	Short:   "Launch interactive TUI",
	Run: func(cmd *cobra.Command, args []string) {
		dc, err := docker.NewDockerClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		p := tea.NewProgram(ui.NewModel(dc), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
