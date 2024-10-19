package content

import (
	"context"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

type CommandMsg string

func Command(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
	}
}

type Model struct {
	width, height int
	table         table.Model
	containers    []types.Container
	dockerClient  *docker.Client
}

func New(dockerClient *docker.Client) Model {

	t, containers := getContainerTable(dockerClient)

	return Model{
		dockerClient: dockerClient,
		containers:   containers,
		table:        t,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case CommandMsg:
		switch msg {
		case "networks","Networks":
			m.table = m.getNetworks()
        case "containers","Containers":
            m.table, _ = getContainerTable(m.dockerClient)
		}
	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) getNetworks() table.Model {
	networks, _ := m.dockerClient.NetworkList(context.Background(), network.ListOptions{})

	columns := []table.Column{
		{Title: "NAME", Width: 14},
		{Title: "DRIVER", Width: 14},
		{Title: "IPv4 Subnet", Width: 20},
		{Title: "IPv4 Gateway", Width: 21},
	}

	var rows []table.Row
	for _, ntr := range networks {
		var subnet string
		var gateway string
		if len(ntr.IPAM.Config) > 0 {
			subnet = ntr.IPAM.Config[0].Subnet
			gateway = ntr.IPAM.Config[0].Gateway
		} else {
            subnet = "-"
            gateway = "-"
        }
		rows = append(rows, table.Row{ntr.Name, ntr.Driver, subnet, gateway})
	}
	s := table.DefaultStyles()
	s.Selected = s.Selected.Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#74c7ec"))
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	t.SetStyles(s)
	return t
}

func getContainerTable(dockerClient *docker.Client) (table.Model, []types.Container) {
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
	return t, containers
}
