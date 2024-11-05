package containers

import (
	dctx "context"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/container"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type Model struct {
	ctx   *context.ProgramContext
	table table.Model
}

func New(ctx *context.ProgramContext) (m Model) {
	m.ctx = ctx
	containers, _ := m.ctx.DockerClient.ContainerList(dctx.Background(), container.ListOptions{})

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 5},
		{Title: "NAME", Width: 15},
		{Title: "State", Width: 20},
		{Title: "Status", Width: 30},
		{Title: "Image", Flex: true},
	}

	rows := make([]table.Row, 0, len(containers))
	for i, container := range containers {
		id := container.ID
		index := strconv.Itoa(i)
		name := strings.Trim(container.Names[0], "/")
		state := container.State
		status := container.Status
		image := container.Image

		row := table.Row{id, index, name, state, status, image}
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

func (m Model) View() string {
	return m.table.View()
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.table.UpdateDimensions(m.ctx.ScreenWidth-2, m.ctx.MainContentHeight)
}
