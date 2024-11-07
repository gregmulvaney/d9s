package networks

import (
	dctx "context"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/network"
	"github.com/gregmulvaney/bubbles/keylist"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/tui/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

var Keymap = constants.Keymap{
	keylist.Item{Key: "<l>", Value: "logs"},
}

type Model struct {
	ctx   *context.ProgramContext
	table table.Model
}

func New(ctx *context.ProgramContext) (m Model) {
	m.ctx = ctx
	networks, _ := m.ctx.DockerClient.NetworkList(dctx.Background(), network.ListOptions{})

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 4},
		{Title: "Name", Flex: true},
	}

	rows := make([]table.Row, 0, len(networks))
	for i, network := range networks {
		index := strconv.Itoa(i)
		id := network.ID
		name := network.Name
		row := table.Row{id, index, name}
		rows = append(rows, row)
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.table.UpdateDimensions(m.ctx.ScreenWidth-2, m.ctx.MainContentHeight)
}

func (m Model) View() string {
	return m.table.View()
}
