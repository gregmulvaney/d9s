package containers

import (
	"context"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

type Model struct {
	width, height int
	WrapperHeight int
	WrapperWidth  int
	containers    []types.Container
	table         table.Model
}

// TODO: Move this to a constants pkg for reuse
var contentStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("4"))

func New(dockerClient *docker.Client) (m Model) {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "#", Width: 3},
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

	m.containers = containers
	m.table = t

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return contentStyle.Width(m.width - 2).Height(m.WrapperHeight).Render(m.table.View())
}
