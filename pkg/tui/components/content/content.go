package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/commands"
	"github.com/gregmulvaney/d9s/pkg/tui/views/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/views/images"
	"github.com/gregmulvaney/d9s/pkg/tui/views/logs"
	"github.com/gregmulvaney/d9s/pkg/tui/views/networks"
	"github.com/gregmulvaney/d9s/pkg/tui/views/volumes"
)

type sessionState int

const (
	containersMode sessionState = iota
	imagesMode
	logsMode
	networksMode
	volumesMode
)

var modes = map[commands.PromptMsg]sessionState{
	"containers": containersMode,
	"images":     imagesMode,
	"networks":   networksMode,
	"volumes":    volumesMode,
}

type Model struct {
	// State
	ctx   *appstate.State
	state sessionState

	// Views
	containers containers.Model
	images     images.Model
	logs       logs.Model
	networks   networks.Model
	volumes    volumes.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx
	// Set initial state
	m.state = containersMode

	// Initialize all views
	m.containers = containers.New(m.ctx)
	m.images = images.New(m.ctx)
	m.networks = networks.New(m.ctx)
	m.volumes = volumes.New(m.ctx)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.containers.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.PromptMsg:
		m.state = modes[msg]
	}

	// Route messages to the active mode
	switch m.state {
	case volumesMode:
		m.volumes, cmd = m.volumes.Update(msg)
		cmds = append(cmds, cmd)
	case networksMode:
		m.networks, cmd = m.networks.Update(msg)
		cmds = append(cmds, cmd)
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
	// Render the view for the active mode
	switch m.state {
	case volumesMode:
		return m.volumes.View()
	case networksMode:
		return m.networks.View()
	case imagesMode:
		return m.images.View()
	default:
		return m.containers.View()
	}
}
