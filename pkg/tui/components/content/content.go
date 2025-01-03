package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/tui/components/command"
	"github.com/gregmulvaney/d9s/pkg/tui/views/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/views/images"
)

type sessionState int

const (
	containersMode sessionState = iota
	imagesMode
	volumesMode
	networksMode
	logsMode
)

var commands = map[command.CommadMsg]sessionState{
	"containers": containersMode,
	"images":     imagesMode,
	"networks":   networksMode,
	"volumes":    volumesMode,
}

type Model struct {
	ctx   *appstate.State
	state sessionState

	containers containers.Model
	images     images.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx
	m.state = containersMode

	m.containers = containers.New(m.ctx)
	m.images = images.New(m.ctx)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.containers.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case command.CommadMsg:
		m.state = commands[msg]
		return m, tea.WindowSize()
	}

	switch m.state {
	case imagesMode:
		m.images, cmd = m.images.Update(msg)
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
	default:
		return m.containers.View()
	}
}
