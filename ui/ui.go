package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/components/command"
	"github.com/gregmulvaney/d9s/ui/components/content"
	"github.com/gregmulvaney/d9s/ui/components/header"
)

type sessionState int

const (
	contentView sessionState = iota
	commandView
)

type Model struct {
	state           sessionState
	width, height   int
	showCommandView bool
	header          header.Model
	command         command.Model
	content         content.Model
}

func New() Model {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	return Model{
		state:           contentView,
		showCommandView: false,
		header:          header.New(dockerClient),
		command:         command.New(),
		content:         content.New(dockerClient),
	}
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
		case tea.KeyCtrlC.String():
			return m, tea.Quit
        case tea.KeyEsc.String():
           if m.showCommandView && m.state == commandView {
                // Pass ESC to command model to clear input
                m.command, _ = m.command.Update(msg)
                // Hide command view 
                m.showCommandView = false
                // Swap focus back to content view
                m.state = contentView
            } 
		case ":":
			if !m.showCommandView {
                // Show command view
				m.showCommandView = true
                // Focus command model
				m.state = commandView
                // Focus input model inside command model
				return m, command.Focused()
			}
		}

    // Route command msg to content model
	case content.CommandMsg:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	switch m.state {
	case commandView:
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
	header := lipgloss.NewStyle().Width(m.width).Height(7).Render(m.header.View())
	var command string
	if m.showCommandView {
		command = lipgloss.NewStyle().
			Width(m.width - 2).
			PaddingLeft(2).
			PaddingRight(2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#a6e3a1")).
			Render(m.command.View())
	}
	content := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - lipgloss.Height(command) - lipgloss.Height(header) - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#74c7ec")).
		Render(m.content.View())
	return lipgloss.JoinVertical(lipgloss.Top, header, command, content)
}
