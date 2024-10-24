package content

import (
	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/tables/containers"
)

type sessionState int

const (
	containersView sessionState = iota
	networksView
	volumesView
)

type Model struct {
	width, height   int
	state           sessionState
	ContainerHeight int
	containers      containers.Model
}

func New(dockerClient *docker.Client) (m Model) {
	m.state = containersView
	m.containers = containers.New(dockerClient)
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

	switch m.state {
	default:
		m.containers, cmd = m.containers.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.containers.SetHeight(m.ContainerHeight)
	return m.containers.View()
}
