package networks

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/network"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

var Keymap = [][]string{
	{"<C+p>", "Prune"},
	{"<C+k>", "Delete"},
}

type (
	fetchNetworksMsg bool
	apiErrorMsg      error
)

type Model struct {
	ctx *appstate.State

	table table.Model
}

func apiError(err error) tea.Cmd {
	return func() tea.Msg {
		return apiError(err)
	}
}

func fetchNetworks() tea.Cmd {
	return func() tea.Msg {
		return fetchNetworksMsg(true)
	}
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 4},
		{Title: "Name", Flex: true},
		{Title: "Driver", Width: 16},
		{Title: "Attachable", Width: 15},
		{Title: "IPAM Driver", Width: 20},
		{Title: "IPv4 Subnet", Width: 20},
		{Title: "IPv4 Gateway", Width: 20},
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithFocus(true),
	)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlP:
			//
		case tea.KeyCtrlK:
			//
		}

	case apiErrorMsg:
		//

	case constants.InitMsg:
		return m, fetchNetworks()

	case fetchNetworksMsg:
		networks, err := m.ctx.Api.NetworksFetch()
		if err != nil {
			return m, apiError(err)
		}
		m.table.SetRows(renderRows(networks))
		return m, tea.WindowSize()

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 12)
		m.table.SetWidth(msg.Width - 2)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func renderRows(networks []network.Summary) []table.Row {
	var rows []table.Row
	for i, network := range networks {
		id := network.ID
		index := strconv.Itoa(i)
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
		rows = append(rows, table.Row{id, index, name, driver, attach, ipamDriver, ipv4Subnet, ipv4Gateway})
	}

	return rows
}
