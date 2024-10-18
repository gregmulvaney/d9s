package content

import (
	"context"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/components/messages"
)

type Model struct {
	height, width int
	dockerClient  *docker.Client
	containers    []dockerTypes.Container
	table         table.Model
}

func New(dockerClient *docker.Client) Model {
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

	return Model{
		dockerClient: dockerClient,
		containers:   containers,
		table:        t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
    case messages.CommandMsg:
        m.getNetworks()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) getNetworks() {
	networks, _ := m.dockerClient.NetworkList(context.Background(), network.ListOptions{})

	columns := []table.Column{
		{Title: "NAME", Width: 14},
	}

	var rows []table.Row
	for _, ntr := range networks {
		rows = append(rows, table.Row{ntr.Name})
	}
	s := table.DefaultStyles()
	s.Selected = s.Selected.Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#74c7ec"))

	m.table.SetStyles(s)
	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
}
