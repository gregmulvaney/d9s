package networks

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

type Model struct {
	width, height int
	table         table.Model
	networks      []network.Inspect
}

func New(dockerClient *docker.Client) Model {
	networks, _ := dockerClient.NetworkList(context.Background(), network.ListOptions{})

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
	return Model{
		networks: networks,
		table:    t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}
