package containers

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

type sessionState int

const (
	containersView sessionState = iota
	containerLogsView
)

type Model struct {
	width, height   int
	state           sessionState
	ContainerHeight int
	containers      []types.Container
	table           table.Model
	dockerClient    *docker.Client
	logData         string
}

func New(dockerClient *docker.Client) (m Model) {
	containers, _ := dockerClient.ContainerList(context.Background(), container.ListOptions{})

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "NAME", Width: 15},
		{Title: "STATE", Width: 20},
		{Title: "STATUS", Width: 20},
		{Title: "IMAGE", Width: 50},
	}

	var rows []table.Row
	for index, ctr := range containers {
		i := strconv.Itoa(index)
		rows = append(rows, table.Row{i, strings.Trim(ctr.Names[0], "/"), ctr.State, ctr.Status, ctr.Image})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#74c7ec"))
	t.SetStyles(s)

	m.table = t
	m.dockerClient = dockerClient
	m.containers = containers
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			index := m.table.SelectedRow()[0]
            m.state = containerLogsView
			m.logData = m.getContainerLogs(index)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.state {
	case containerLogsView:

	default:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.table.SetHeight(m.ContainerHeight)
	switch m.state {
	case containerLogsView:
		return lipgloss.NewStyle().Width(m.width-4).MaxWidth(m.width - 4).Render(m.logData)
	default:
		return m.table.View()
	}
}

func (m *Model) SetHeight(height int) {
	m.ContainerHeight = height
}

func (m *Model) getContainerLogs(index string) string {
	i, _ := strconv.Atoi(index)
	ctr := m.containers[i]
	buf := new(strings.Builder)
	reader, _ := m.dockerClient.ContainerLogs(context.Background(), ctr.ID, container.LogsOptions{ShowStdout: true})
	defer reader.Close()
	_, _ = io.Copy(buf, reader)
	return buf.String()
}
