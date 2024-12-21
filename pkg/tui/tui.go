package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/components/command"
	"github.com/gregmulvaney/d9s/pkg/views/content"
	"github.com/gregmulvaney/d9s/pkg/views/header"
)

type sessionState int

const (
	contentMode sessionState = iota
	splashMode
	commandMode
)

type Model struct {
	state       sessionState
	height      int
	width       int
	showCommand bool
	ctx         appcontext.Context

	header  header.Model
	command command.Model
	content content.Model
}

func New() (m Model) {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	// Initialize app context struct
	m.ctx = appcontext.Context{}
	m.ctx.Docker = client

	// Set initial interaction state
	m.state = contentMode

	m.header = header.New(&m.ctx)
	m.command = command.New()
	m.content = content.New(&m.ctx)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Show command prompt
		case ":":
			m.showCommand = true
			m.command.Focus()
			m.state = commandMode
			return m, tea.WindowSize()
		}
		switch msg.Type {
		// Quit app
		case tea.KeyCtrlC:
			return m, tea.Quit

		// Hide command prompt and reset input
		case tea.KeyEsc:
			if m.showCommand {
				m.showCommand = false
				m.command.Reset()
				m.state = contentMode
				return m, tea.WindowSize()
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	// Switch which view is receiving messages
	switch m.state {
	case commandMode:
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().
		Width(m.width).
		MaxWidth(m.width).
		Render(m.header.View())

	var command string
	if m.showCommand {
		command = lipgloss.
			NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("2")).
			Width(m.width - 2).PaddingLeft(1).
			Render(m.command.View())
	}

	_, headerHeight := lipgloss.Size(header)
	_, commandHeight := lipgloss.Size(command)

	content := lipgloss.NewStyle().
		Height(m.height - headerHeight - commandHeight - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("4")).
		Width(m.width - 2).
		Render(m.content.View())

	return lipgloss.JoinVertical(lipgloss.Left, header, command, content)
}
