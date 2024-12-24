package networks

import (
	"context"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/network"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx   *appcontext.Context
	table table.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx

	networks, err := m.ctx.Docker.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		panic(err)
	}

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "Name", Width: 20},
		{Title: "Driver", Width: 16},
		{Title: "Attachable", Width: 15},
		{Title: "IPAM Driver", Width: 20},
		{Title: "IPv4 Subnet", Width: 20},
		{Title: "IPv4 Gateway", Width: 20},
	}

	var rows []table.Row

	for _, network := range networks {
		id := network.ID
		name := network.Name
		driver := network.Driver
		attach := strconv.FormatBool(network.Attachable)
		ipamDriver := network.IPAM.Driver

		var ipv4Subnet string
		var ipv4Gateway string
		if driver != "host" && driver != "null" {
			ipv4Subnet = network.IPAM.Config[0].Subnet
			ipv4Gateway = network.IPAM.Config[0].Gateway
		} else {
			ipv4Gateway = "-"
			ipv4Subnet = "-"
		}
		rows = append(rows, table.Row{id, name, driver, attach, ipamDriver, ipv4Subnet, ipv4Gateway})
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocus(true),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 11)
		m.table.SetWidth(msg.Width - 2)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}
