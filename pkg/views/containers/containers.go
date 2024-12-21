package containers

import (
	"context"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx        *appcontext.Context
	containers []types.Container
	table      table.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx

	containers, err := m.ctx.Docker.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	m.containers = containers

	cols := []table.Column{
		{Title: "id", Hidden: true},
		{Title: "#", Width: 3},
		{Title: "Name", Flex: true},
	}

	var rows []table.Row

	for i, container := range containers {
		name := container.Names[0]
		id := container.ID
		index := strconv.Itoa(i)
		rows = append(rows, table.Row{id, index, name})
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(10)
		m.table.SetWidth(msg.Width - 2)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
	// return "containers"
}
