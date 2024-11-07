package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/tui/components/command"
	"github.com/gregmulvaney/d9s/pkg/tui/components/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/components/networks"
	"github.com/gregmulvaney/d9s/pkg/tui/components/volumes"
	"github.com/gregmulvaney/d9s/pkg/tui/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type sessionState int

type KeymapMsg constants.Keymap

func UpdateKeymap(keymap constants.Keymap) tea.Cmd {
	return func() tea.Msg {
		return KeymapMsg(keymap)
	}
}

const (
	containersMode sessionState = iota
	networksMode
	volumesMode
)

type Model struct {
	ctx        *context.ProgramContext
	state      sessionState
	containers containers.Model
	networks   networks.Model
	volumes    volumes.Model
}

func New(ctx context.ProgramContext) (m Model) {
	m.ctx = &ctx
	m.state = containersMode
	m.containers = containers.New(m.ctx)
	m.networks = networks.New(m.ctx)
	m.volumes = volumes.New(m.ctx)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case command.CommandMsg:
		switch msg {
		case "containers":
			m.state = containersMode
			return m, UpdateKeymap(containers.Keymap)
		case "networks":
			m.state = networksMode
			return m, UpdateKeymap(networks.Keymap)
		case "volumes":
			m.state = volumesMode
		}
	}

	switch m.state {
	case volumesMode:
		m.volumes, cmd = m.volumes.Update(msg)
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

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.containers.UpdateProgramContext(ctx)
	m.networks.UpdateProgramContext(ctx)
	m.volumes.UpdateProgramContext(ctx)
}

func (m Model) View() string {
	var content string

	switch m.state {
	case volumesMode:
		content = m.volumes.View()
	case networksMode:
		content = m.networks.View()
	default:
		content = m.containers.View()
	}

	return lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth - 2).
		MaxWidth(m.ctx.ScreenWidth).
		Height(m.ctx.MainContentHeight).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("4")).
		Render(content)
}
