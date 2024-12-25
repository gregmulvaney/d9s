package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/components/command"
	"github.com/gregmulvaney/d9s/pkg/views/containers"
	"github.com/gregmulvaney/d9s/pkg/views/images"
	"github.com/gregmulvaney/d9s/pkg/views/networks"
)

type sessionState int

const (
	containersMode sessionState = iota
	networksMode
	imagesMode
)

var commands = map[command.CommandMsg]sessionState{
	"containers": containersMode,
	"networks":   networksMode,
	"images":     imagesMode,
}

type Model struct {
	ctx        *appcontext.Context
	state      sessionState
	containers containers.Model
	networks   networks.Model
	images     images.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.state = containersMode

	m.containers = containers.New(ctx)
	m.networks = networks.New(ctx)
	m.images = images.New(ctx)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case command.CommandMsg:
		m.state = commands[msg]
		return m, tea.WindowSize()
	}

	switch m.state {
	case imagesMode:
		m.images, cmd = m.images.Update(msg)
		cmds = append(cmds, cmd)
	case networksMode:
		m.networks, cmd = m.networks.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.containers, cmd = m.containers.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case imagesMode:
		return m.images.View()
	case networksMode:
		return m.networks.View()
	default:
		return m.containers.View()
	}
}
