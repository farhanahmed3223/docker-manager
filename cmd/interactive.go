package cmd

import (
    "fmt"
    "os"

    "docker-manager/internal/docker"
    "docker-manager/internal/ui"

    "github.com/spf13/cobra"
    tea "github.com/charmbracelet/bubbletea"
)

var compactMode bool

var interactiveCmd = &cobra.Command{
    Use:   "interactive",
    Short: "Launch interactive TUI mode",
    Long:  `Launch the full interactive terminal UI for Docker container management.`,
    Run: func(cmd *cobra.Command, args []string) {
        dockerClient, err := docker.NewDockerClient()
        if err != nil {
            fmt.Printf("Error connecting to Docker: %v\n", err)
            os.Exit(1)
        }

        model := ui.NewModel(dockerClient, compactMode)
        p := tea.NewProgram(model, tea.WithAltScreen())

        if _, err := p.Run(); err != nil {
            fmt.Printf("Error running application: %v\n", err)
            os.Exit(1)
        }
    },
}

func init() {
    interactiveCmd.Flags().BoolVarP(&compactMode, "compact", "c", false, "Use compact view")
}
