package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/components/containers"
	"github.com/gregmulvaney/d9s/pkg/components/networks"
	"github.com/gregmulvaney/d9s/pkg/components/secrets"
	"github.com/gregmulvaney/d9s/pkg/components/volumes"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

var contentStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("4"))

type sessionState int

const (
	containersMode sessionState = iota
	networksMode
	volumesMode
	secretsMode
)

type Model struct {
	ctx        *appcontext.Context
	state      sessionState
	containers containers.Model
	networks   networks.Model
	volumes    volumes.Model
	secrets    secrets.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx

	m.state = containersMode

	m.containers = containers.New(m.ctx)
	m.networks = networks.New(m.ctx)
	m.volumes = volumes.New(m.ctx)
	m.secrets = secrets.New(m.ctx)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case constants.CommandMsg:
		switch msg {
		case "containers":
			m.state = containersMode
		case "networks":
			m.state = networksMode
		case "volumes":
			m.state = volumesMode
		case "secrets":
			m.state = secretsMode
		}
	}

	switch m.state {
	default:
		m.containers, cmd = m.containers.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var content string
	switch m.state {
	case networksMode:
		content = m.networks.View()
	case volumesMode:
		content = m.volumes.View()
	case secretsMode:
		content = m.secrets.View()
	default:
		content = m.containers.View()
	}

	return contentStyle.Width(m.ctx.ScreenWidth - 3).
		MaxWidth(m.ctx.ScreenWidth).
		Height(m.ctx.MainContentHeight - 3).
		MaxHeight(m.ctx.MainContentHeight).
		Render(content)
}

func (m *Model) SyncAppContext(ctx *appcontext.Context) {
	m.ctx = ctx
	m.containers.SyncAppContext(m.ctx)
	m.networks.SyncProgramContext(m.ctx)
	m.volumes.SyncProgramContext(m.ctx)
	m.secrets.SyncProgramContext(m.ctx)
}
