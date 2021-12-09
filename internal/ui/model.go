package ui

import (
	"fmt"
	"strings"

	"docker-manager/internal/docker"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	dockerClient *docker.DockerClient
	table        table.Model
	containers   []docker.ContainerInfo
	err          error
	loading      bool
}

func NewModel(dc *docker.DockerClient) Model {
	cols := []table.Column{
		{Title: "ID", Width: 13},
		{Title: "Name", Width: 20},
		{Title: "Image", Width: 25},
		{Title: "Status", Width: 15},
	}
	t := table.New(table.WithColumns(cols), table.WithFocused(true), table.WithHeight(15))
	return Model{dockerClient: dc, table: t, loading: true}
}

type containersMsg []docker.ContainerInfo
type errMsg struct{ error }

func (m Model) Init() tea.Cmd {
	return fetchContainers(m.dockerClient)
}

func fetchContainers(dc *docker.DockerClient) tea.Cmd {
	return func() tea.Msg {
		containers, err := dc.ListContainers(nil)
		if err != nil {
			return errMsg{err}
		}
		return containersMsg(containers)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case containersMsg:
		m.loading = false
		m.containers = msg
		rows := make([]table.Row, len(msg))
		for i, c := range msg {
			rows[i] = table.Row{c.ID, c.Name, c.Image, c.Status}
		}
		m.table.SetRows(rows)
	case errMsg:
		m.loading = false
		m.err = msg.error
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return "Loading containers..."
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}
	title := lipgloss.NewStyle().Bold(true).Render("Docker Manager")
	help := "\n  ↑/↓ navigate  q quit"
	return title + "\n\n" + m.table.View() + help
}
