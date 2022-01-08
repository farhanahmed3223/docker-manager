package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7c3aed")).
			Padding(0, 1)

	statusRunning = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e"))
	statusStopped = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))
	statusPaused  = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Padding(1, 0)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#374151"))
)

func colorStatus(status string) string {
	switch {
	case len(status) > 2 && status[:2] == "Up":
		return statusRunning.Render(status)
	case status == "paused":
		return statusPaused.Render(status)
	default:
		return statusStopped.Render(status)
	}
}
