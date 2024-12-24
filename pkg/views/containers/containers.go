package containers

import (
	"context"
	"strconv"
	"strings"

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

	containers, err := m.ctx.Docker.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	m.containers = containers

	cols := []table.Column{
		{Title: "id", Hidden: true},
		{Title: "#", Width: 4},
		{Title: "Name", Width: 25},
		{Title: "Status", Width: 12},
		{Title: "State", Width: 35},
		{Title: "Image", Flex: true},
	}

	rows := renderRows(containers)

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
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.fetchContainers()
			return m, tea.WindowSize()
		}
		switch msg.Type {
		case tea.KeyCtrlK:
			selected := m.table.SelectedRow()
			m.ctx.Docker.ContainerStop(context.Background(), selected[0], container.StopOptions{})
			m.fetchContainers()
			return m, tea.WindowSize()
		case tea.KeyCtrlS:
			selected := m.table.SelectedRow()
			m.ctx.Docker.ContainerStart(context.Background(), selected[0], container.StartOptions{})
			m.fetchContainers()
			return m, tea.WindowSize()
		}
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

func (m *Model) fetchContainers() {
	containers, err := m.ctx.Docker.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	rows := renderRows(containers)
	m.table.SetRows(rows)

	m.containers = containers
}

func renderRows(containers []types.Container) []table.Row {
	var rows []table.Row

	for i, container := range containers {
		name := strings.Trim(container.Names[0], "/")
		id := container.ID
		index := strconv.Itoa(i)
		rows = append(rows, table.Row{id, index, name, container.State, container.Status, container.Image})
	}

	return rows
}
