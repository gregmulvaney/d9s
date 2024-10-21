package content

import (
	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/tables/containers"
	"github.com/gregmulvaney/d9s/ui/tables/networks"
)

type sessionState int

// Messages

type CommandMsg string

func Command(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
	}
}

const (
	containerView sessionState = iota
	networkView
	volumeView
)

type Model struct {
	state      sessionState
	containers containers.Model
	networks   networks.Model
}

func New(dockerClient *docker.Client) Model {
	return Model{
		state:      containerView,
		containers: containers.New(dockerClient),
		networks:   networks.New(dockerClient),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CommandMsg:
		switch msg {
		case "Containers", "containers":
			m.state = containerView
		case "Networks", "networks":
			m.state = networkView
		}
	}

	switch m.state {
	case containerView:
		m.containers, cmd = m.containers.Update(msg)
		cmds = append(cmds, cmd)
	case networkView:
		m.networks, cmd = m.networks.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case containerView:
		return m.containers.View()
	case networkView:
		return m.networks.View()
	}
	return "content"
}
