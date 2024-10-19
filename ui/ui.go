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
	dockerClient    *docker.Client
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
		dockerClient:    dockerClient,
		showCommandView: false,
		header:          header.New(dockerClient),
		command:         command.New(),
		content:         content.New(dockerClient),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Handle key messages
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case ":":
			if !m.showCommandView {
				m.state = commandView
				m.showCommandView = true
				return m, command.Focus()
			}
        case "esc":
            if m.showCommandView {
                m.state = contentView
                m.showCommandView = false
            }
		}
	// Pass command to content
	case content.CommandMsg:
		m.state = contentView
		m.showCommandView = false
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
		return m, command.Clear()
	// Handle window size messages
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
