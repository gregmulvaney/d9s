package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/tui/components/command"
	"github.com/gregmulvaney/d9s/pkg/tui/components/content"
	"github.com/gregmulvaney/d9s/pkg/tui/components/header"
	"github.com/gregmulvaney/d9s/pkg/tui/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type sessionState int

const (
	contentMode sessionState = iota
	commandMode
)

type Model struct {
	ctx     context.ProgramContext
	state   sessionState
	header  header.Model
	command command.Model
	content content.Model
}

func New() (m Model) {
	m.ctx = context.ProgramContext{}
	m.state = contentMode
	m.header = header.New(m.ctx)
	m.command = command.New(m.ctx)
	m.content = content.New(m.ctx)
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case ":":
			if !m.ctx.ShowCommandView {
				m.ctx.ShowCommandView = true
				m.state = commandMode
				m.command.Focus()
				return m, tea.WindowSize()
			}
		case tea.KeyEsc.String():
			if m.ctx.ShowCommandView {
				m.ctx.ShowCommandView = false
				m.state = contentMode
				m.command.Reset()
				return m, tea.WindowSize()
			}
		}
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	}

	switch m.state {
	case commandMode:
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.syncProgramContext()

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().Height(constants.HeaderHeight).MaxHeight(constants.HeaderHeight).Render(m.header.View())

	var command string
	if m.ctx.ShowCommandView {
		command = m.command.View()
	}

	content := m.content.View()

	return lipgloss.JoinVertical(lipgloss.Top, header, command, content)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.ctx.ScreenHeight = msg.Height
	m.ctx.ScreenWidth = msg.Width
	if m.ctx.ShowCommandView {
		m.ctx.MainContentHeight = msg.Height - constants.HeaderHeight - constants.CommandHeight - 2
	} else {
		m.ctx.MainContentHeight = msg.Height - constants.HeaderHeight - 3
	}
}

func (m *Model) syncProgramContext() {
	m.header.UpdateProgramContenxt(&m.ctx)
	m.command.UpdateProgramContext(&m.ctx)
	m.content.UpdateProgramContext(&m.ctx)
}
