package ui

import "github.com/charmbracelet/lipgloss"

var (
    // Colors
    PrimaryColor    = lipgloss.Color("69")
    SecondaryColor  = lipgloss.Color("99")
    SuccessColor    = lipgloss.Color("46")
    WarningColor    = lipgloss.Color("214")
    DangerColor     = lipgloss.Color("196")
    MutedColor      = lipgloss.Color("240")

    // Styles
    TitleStyle = lipgloss.NewStyle().
        Foreground(PrimaryColor).
        Bold(true).
        Padding(0, 1)

    StatusBarStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("15")).
        Background(MutedColor).
        Padding(0, 1)

    ContainerRunningStyle = lipgloss.NewStyle().Foreground(SuccessColor)
    ContainerStoppedStyle = lipgloss.NewStyle().Foreground(DangerColor)
    ContainerPausedStyle  = lipgloss.NewStyle().Foreground(WarningColor)

    // Table styles
    HeaderStyle = lipgloss.NewStyle().
        Foreground(PrimaryColor).
        Bold(true).
        Padding(0, 1)

    CellStyle = lipgloss.NewStyle().Padding(0, 1)

    SelectedStyle = lipgloss.NewStyle().
        Background(SecondaryColor).
        Foreground(lipgloss.Color("15")).
        Padding(0, 1)

    // Stats styles
    HighUsageStyle = lipgloss.NewStyle().Foreground(DangerColor).Bold(true)
    MediumUsageStyle = lipgloss.NewStyle().Foreground(WarningColor)
    LowUsageStyle = lipgloss.NewStyle().Foreground(SuccessColor)

    // Help styles
    HelpStyle = lipgloss.NewStyle().
        Foreground(MutedColor).
        Italic(true)
)

func GetUsageStyle(value float64) lipgloss.Style {
    switch {
    case value > 80:
        return HighUsageStyle
    case value > 60:
        return MediumUsageStyle
    default:
        return LowUsageStyle
    }
}
