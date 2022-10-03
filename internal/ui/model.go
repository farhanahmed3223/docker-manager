package ui

import (
	"context"
	"fmt"
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

type ViewType int

const (
	ContainersView ViewType = iota
	LogsView
	FilterView
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

type tickMsg time.Time
type containersMsg []docker.ContainerInfo
type errorMsg struct{ error }

func NewModel(dc *docker.DockerClient, compact bool) Model {
	cols := []table.Column{
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
	t := table.New(table.WithColumns(cols), table.WithFocused(true), table.WithHeight(15))
	ti := textinput.New()
	ti.Placeholder = "Filter containers..."
	ti.CharLimit = 64
	return Model{
		dockerClient: dc,
		table:        t,
		textinput:    ti,
		loading:      true,
		compactMode:  compact,
	}
}

func tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func fetchContainers(dc *docker.DockerClient) tea.Cmd {
	return func() tea.Msg {
		cs, err := dc.ListContainers(context.Background())
		if err != nil {
			return errorMsg{err}
		}
		return containersMsg(cs)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchContainers(m.dockerClient), tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tea.Batch(fetchContainers(m.dockerClient), tick())
	case containersMsg:
		m.loading = false
		m.containers = msg
		m.refreshTable()
	case errorMsg:
		m.loading = false
		m.err = msg.error
	case tea.WindowSizeMsg:
		// fix: update table height when terminal is resized
		m.width, m.height = msg.Width, msg.Height
		tableHeight := m.height - 6
		if tableHeight < 4 {
			tableHeight = 4
		}
		m.table.SetHeight(tableHeight)
		m.viewport = viewport.New(m.width, tableHeight)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Refresh):
			m.loading = true
			return m, fetchContainers(m.dockerClient)
		case key.Matches(msg, keys.Filter):
			m.currentView = FilterView
			m.textinput.Focus()
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) refreshTable() {
	rows := make([]table.Row, 0, len(m.containers))
	for _, c := range m.containers {
		if m.filter != "" && !strings.Contains(strings.ToLower(c.Name), strings.ToLower(m.filter)) {
			continue
		}
		rows = append(rows, table.Row{c.ID, c.Name, c.Image, c.Status, "", "—", "—", "—", "—"})
	}
	m.table.SetRows(rows)
}

func (m Model) View() string {
	if m.loading {
		return "\n  ⟳ Fetching containers...\n"
	}
	if m.err != nil {
		return fmt.Sprintf("\n  ❌ Error: %v\n\n  Press q to quit.\n", m.err)
	}
	header := titleStyle.Render("🐋 Docker Manager")
	if m.filter != "" {
		header += "  " + lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Render("filter: "+m.filter)
	}
	help := helpStyle.Render("↑/↓ navigate  s start  x stop  l logs  / filter  r refresh  q quit")
	return "\n" + header + "\n\n" + m.table.View() + "\n" + help
}
