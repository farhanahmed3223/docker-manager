package ui

import (
    "context"
    "fmt"
    "io"
    "strings"
    "time"

    "docker-manager/internal/docker"

    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    dockerClient *docker.DockerClient
    table        table.Model
    viewport     viewport.Model
    textinput    textinput.Model
    containers   []docker.ContainerInfo
    selectedID   string
    currentView  ViewType
    err          error
    loading      bool
    filter       string
    compactMode  bool
    width        int
    height       int
}

type ViewType int

const (
    ContainersView ViewType = iota
    LogsView
    FilterView
)

type tickMsg time.Time
type containersMsg []docker.ContainerInfo
type errorMsg struct{ error }

func NewModel(dockerClient *docker.DockerClient, compact bool) Model {
    // Initialize table
    columns := []table.Column{
        {Title: "ID", Width: 12},
        {Title: "Name", Width: 20},
        {Title: "Image", Width: 25},
        {Title: "Status", Width: 15},
        {Title: "Ports", Width: 20},
        {Title: "CPU%", Width: 8},
        {Title: "Memory%", Width: 10},
        {Title: "Network", Width: 15},
        {Title: "Uptime", Width: 15},
    }

    if compact {
        columns = []table.Column{
            {Title: "ID", Width: 12},
            {Title: "Name", Width: 20},
            {Title: "Status", Width: 15},
            {Title: "CPU%", Width: 8},
            {Title: "Memory%", Width: 10},
        }
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithFocused(true),
        table.WithHeight(10),
    )

    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        BorderBottom(true).
        Bold(true)
    s.Selected = SelectedStyle
    t.SetStyles(s)

    // Initialize viewport for logs
    vp := viewport.New(80, 20)
    vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(PrimaryColor)

    // Initialize text input for filtering
    ti := textinput.New()
    ti.Placeholder = "Filter by name, status..."
    ti.CharLimit = 50
    ti.Width = 50

    return Model{
        dockerClient: dockerClient,
        table:        t,
        viewport:     vp,
        textinput:    ti,
        currentView:  ContainersView,
        compactMode:  compact,
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.refreshContainers(),
        tickCmd(),
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.currentView {
        case ContainersView:
            return m.updateContainersView(msg)
        case LogsView:
            return m.updateLogsView(msg)
        case FilterView:
            return m.updateFilterView(msg)
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.table.SetHeight(msg.Height - 10)
        m.viewport.Height = msg.Height - 10
        m.viewport.Width = msg.Width - 4

    case containersMsg:
        m.loading = false
        m.containers = msg
        m.updateTableRows()
        cmds = append(cmds, tickCmd())

    case errorMsg:
        m.err = msg.error
        m.loading = false

    case tickMsg:
        cmds = append(cmds, m.refreshContainers(), tickCmd())
    }

    return m, tea.Batch(cmds...)
}

func (m *Model) updateContainersView(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch {
    case key.Matches(msg, Keys.Quit):
        return m, tea.Quit

    case key.Matches(msg, Keys.Logs):
        if m.table.SelectedRow() != nil {
            m.currentView = LogsView
            return m, m.loadLogs()
        }

    case key.Matches(msg, Keys.Filter):
        m.currentView = FilterView
        m.textinput.Focus()
        return m, nil

    case key.Matches(msg, Keys.Start):
        return m, m.startContainer()

    case key.Matches(msg, Keys.Stop):
        return m, m.stopContainer()

    case key.Matches(msg, Keys.Restart):
        return m, m.restartContainer()

    case key.Matches(msg, Keys.Remove):
        return m, m.removeContainer()

    case key.Matches(msg, Keys.Refresh):
        return m, m.refreshContainers()

    case key.Matches(msg, Keys.Help):
        // Toggle help
        return m, nil
    }

    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m *Model) updateLogsView(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch {
    case key.Matches(msg, Keys.Back):
        m.currentView = ContainersView
        m.viewport.SetContent("")
    case key.Matches(msg, Keys.Quit):
        return m, tea.Quit
    }

    var cmd tea.Cmd
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

func (m *Model) updateFilterView(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch {
    case key.Matches(msg, Keys.Enter):
        m.filter = m.textinput.Value()
        m.currentView = ContainersView
        m.textinput.Blur()
        return m, m.refreshContainers()

    case key.Matches(msg, Keys.Back):
        m.currentView = ContainersView
        m.textinput.Blur()
        return m, nil
    }

    var cmd tea.Cmd
    m.textinput, cmd = m.textinput.Update(msg)
    return m, cmd
}

func (m Model) View() string {
    if m.err != nil {
        return fmt.Sprintf("Error: %v\nPress q to quit", m.err)
    }

    var view string
    switch m.currentView {
    case ContainersView:
        view = m.containersView()
    case LogsView:
        view = m.logsView()
    case FilterView:
        view = m.filterView()
    }

    return view
}

func (m Model) containersView() string {
    var b strings.Builder

    // Title
    b.WriteString(TitleStyle.Render("🐳 Docker Container Manager"))
    b.WriteString("\n\n")

    // Table
    b.WriteString(m.table.View())
    b.WriteString("\n\n")

    // Status bar
    status := fmt.Sprintf("Containers: %d", len(m.containers))
    if m.filter != "" {
        status += fmt.Sprintf(" | Filter: %s", m.filter)
    }
    if m.loading {
        status += " | Refreshing..."
    }
    b.WriteString(StatusBarStyle.Render(status))
    b.WriteString("\n\n")

    // Help
    b.WriteString(m.helpView())

    return b.String()
}

func (m Model) logsView() string {
    var b strings.Builder

    b.WriteString(TitleStyle.Render("📋 Logs - " + m.selectedID))
    b.WriteString("\n\n")

    b.WriteString(m.viewport.View())
    b.WriteString("\n\n")

    b.WriteString(HelpStyle.Render("↑/↓: Scroll • q/esc: Back to containers"))

    return b.String()
}

func (m Model) filterView() string {
    var b strings.Builder

    b.WriteString(TitleStyle.Render("🔍 Filter Containers"))
    b.WriteString("\n\n")

    b.WriteString("Enter filter (name, status):\n")
    b.WriteString(m.textinput.View())
    b.WriteString("\n\n")

    b.WriteString(HelpStyle.Render("Enter: Apply • esc: Cancel"))

    return b.String()
}

func (m Model) helpView() string {
    if m.currentView == ContainersView {
        return HelpStyle.Render(
            "←/→/↑/↓: Navigate • s: Start • t: Stop • r: Restart • d: Remove • l: Logs • f: Filter • F5: Refresh • q: Quit",
        )
    }
    return ""
}

// Command functions
func (m *Model) refreshContainers() tea.Cmd {
    return func() tea.Msg {
        m.loading = true
        containers, err := m.dockerClient.ListContainers(true)
        if err != nil {
            return errorMsg{err}
        }

        // Apply filter
        if m.filter != "" {
            filtered := []docker.ContainerInfo{}
            for _, c := range containers {
                if strings.Contains(strings.ToLower(c.Name), strings.ToLower(m.filter)) ||
                    strings.Contains(strings.ToLower(c.Status), strings.ToLower(m.filter)) ||
                    strings.Contains(strings.ToLower(c.Image), strings.ToLower(m.filter)) {
                    filtered = append(filtered, c)
                }
            }
            containers = filtered
        }

        return containersMsg(containers)
    }
}

func (m *Model) loadLogs() tea.Cmd {
    return func() tea.Msg {
        if m.table.SelectedRow() == nil {
            return errorMsg{fmt.Errorf("no container selected")}
        }

        containerID := m.table.SelectedRow()[0]
        m.selectedID = containerID

        logs, err := m.dockerClient.GetContainerLogs(containerID)
        if err != nil {
            return errorMsg{err}
        }

        m.viewport.SetContent(logs)
        m.viewport.GotoBottom()

        return nil
    }
}

func (m *Model) startContainer() tea.Cmd {
    return func() tea.Msg {
        if m.table.SelectedRow() == nil {
            return errorMsg{fmt.Errorf("no container selected")}
        }

        containerID := m.table.SelectedRow()[0]
        if err := m.dockerClient.StartContainer(containerID); err != nil {
            return errorMsg{err}
        }

        return m.refreshContainers()()
    }
}

func (m *Model) stopContainer() tea.Cmd {
    return func() tea.Msg {
        if m.table.SelectedRow() == nil {
            return errorMsg{fmt.Errorf("no container selected")}
        }

        containerID := m.table.SelectedRow()[0]
        if err := m.dockerClient.StopContainer(containerID); err != nil {
            return errorMsg{err}
        }

        return m.refreshContainers()()
    }
}

func (m *Model) restartContainer() tea.Cmd {
    return func() tea.Msg {
        if m.table.SelectedRow() == nil {
            return errorMsg{fmt.Errorf("no container selected")}
        }

        containerID := m.table.SelectedRow()[0]
        if err := m.dockerClient.RestartContainer(containerID); err != nil {
            return errorMsg{err}
        }

        return m.refreshContainers()()
    }
}

func (m *Model) removeContainer() tea.Cmd {
    return func() tea.Msg {
        if m.table.SelectedRow() == nil {
            return errorMsg{fmt.Errorf("no container selected")}
        }

        containerID := m.table.SelectedRow()[0]
        // In a real implementation, we'd show a confirmation dialog
        if err := m.dockerClient.RemoveContainer(containerID); err != nil {
            return errorMsg{err}
        }

        return m.refreshContainers()()
    }
}

func (m *Model) updateTableRows() {
    var rows []table.Row
    for _, c := range m.containers {
        var statusStyle lipgloss.Style
        switch {
        case strings.Contains(c.Status, "Up"):
            statusStyle = ContainerRunningStyle
        case strings.Contains(c.Status, "Exited"):
            statusStyle = ContainerStoppedStyle
        default:
            statusStyle = ContainerPausedStyle
        }

        uptime := time.Since(c.Created).Truncate(time.Second).String()

        if m.compactMode {
            rows = append(rows, table.Row{
                c.ID,
                c.Name,
                statusStyle.Render(c.Status),
                GetUsageStyle(c.CPU).Render(fmt.Sprintf("%.1f", c.CPU)),
                GetUsageStyle(c.Memory).Render(fmt.Sprintf("%.1f", c.Memory)),
            })
        } else {
            rows = append(rows, table.Row{
                c.ID,
                c.Name,
                c.Image,
                statusStyle.Render(c.Status),
                c.Ports,
                GetUsageStyle(c.CPU).Render(fmt.Sprintf("%.1f", c.CPU)),
                GetUsageStyle(c.Memory).Render(fmt.Sprintf("%.1f", c.Memory)),
                c.Network,
                uptime,
            })
        }
    }
    m.table.SetRows(rows)
}

func tickCmd() tea.Cmd {
    return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
