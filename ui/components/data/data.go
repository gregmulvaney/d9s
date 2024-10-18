package data

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

type Model struct {
	viewport      viewport.Model
	table         table.Model
	width, height int
	dockerClient  *docker.Client
	containers    []types.Container
}

func NewModel(dockerClient *docker.Client) Model {
	vp := viewport.New(100, 20)

	containers, _ := dockerClient.ContainerList(context.Background(), container.ListOptions{})

	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 14},
		{Title: "Image", Width: 50},
	}

	var rows []table.Row
	for _, ctr := range containers {
		rows = append(rows, table.Row{ctr.Names[0], ctr.State, ctr.Image})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return Model{
		viewport:     vp,
		table:        t,
		dockerClient: dockerClient,
		containers:   containers,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

