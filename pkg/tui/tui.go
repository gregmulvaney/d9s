package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/breadcrumbs"
	"github.com/gregmulvaney/d9s/pkg/api"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/components/command"
	"github.com/gregmulvaney/d9s/pkg/tui/components/content"
	"github.com/gregmulvaney/d9s/pkg/tui/components/header"
)

type sessionState int

const (
	contentMode sessionState = iota
	commandMode
)

type Model struct {
	ctx         appstate.State
	state       sessionState
	height      int
	width       int
	showCommand bool

	header  header.Model
	command command.Model
	content content.Model
	crumbs  breadcrumbs.Model
}

func New() (m Model) {
	m.ctx = appstate.State{}
	m.ctx.Api = api.Init()

	m.state = contentMode

	m.header = header.New(&m.ctx)
	m.command = command.New()
	m.content = content.New(&m.ctx)

	crumbs := []string{"Containers"}
	m.crumbs = breadcrumbs.New(breadcrumbs.WithCrumbs(crumbs))
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.content.Init(), m.command.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.showCommand {
				m.showCommand = false
				m.state = contentMode
				m.command.Reset()
				return m, tea.WindowSize()
			}
		}

		switch msg.String() {
		case ":":
			if !m.showCommand {
				m.showCommand = true
				m.state = commandMode
				m.command.Focus()
				return m, tea.WindowSize()
			}
		}

	case command.CommadMsg:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
		m.showCommand = false
		m.state = contentMode

	case constants.ApiErrorMsg:
		// TODO

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	switch m.state {
	case commandMode:
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().
		Width(m.width).
		MaxWidth(m.width).
		Render(m.header.View())

	var command string
	commandHeight := 0

	if m.showCommand {
		command = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("2")).
			Width(m.width - 2).
			Render(m.command.View())
		// Measure size of the command component
		_, commandHeight = lipgloss.Size(command)
	}

	// Measure the size of the header component
	_, headerHeight := lipgloss.Size(header)

	// Set size of the content component - 3 for borders and breadcrumb height
	contentHeight := m.height - commandHeight - headerHeight - 3

	content := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("4")).
		Width(m.width - 2).
		Height(contentHeight).
		Render(m.content.View())

	breadcrumbs := lipgloss.NewStyle().Render(m.crumbs.View())

	if m.showCommand {
		return lipgloss.JoinVertical(lipgloss.Top, header, command, content, breadcrumbs)
	}
	return lipgloss.JoinVertical(lipgloss.Top, header, content, breadcrumbs)
}
