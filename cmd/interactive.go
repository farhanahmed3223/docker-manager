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
			fmt.Fprintln(os.Stderr, "Is Docker running? Try: docker ps")
			os.Exit(1)
		}
		p := tea.NewProgram(
			ui.NewModel(dc, compactMode),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error running TUI:", err)
			os.Exit(1)
		}
	},
}
