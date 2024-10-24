package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/views/command"
	"github.com/gregmulvaney/d9s/ui/views/content"
	"github.com/gregmulvaney/d9s/ui/views/header"
)

type sessionState int

const (
	contentView sessionState = iota
	commandView
)

type Model struct {
	width, height   int
	state           sessionState
	showCommandView bool
	header          header.Model
	command         command.Model
	content         content.Model
}

func New() (m Model) {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	m.state = contentView
	m.showCommandView = false
	m.header = header.New(dockerClient)
	m.command = command.New()
	m.content = content.New(dockerClient)

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
		// Quit key combo
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		case ":":
			if !m.showCommandView {
				m.showCommandView = true
				m.command.Focus()
				m.state = commandView
			}
		case tea.KeyEsc.String():
			if m.showCommandView {
				m.showCommandView = false
				m.command.Clear()
				m.state = contentView
			}
		}

	case command.CommandMsg:
		m.showCommandView = false
		m.state = contentView

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
	header := lipgloss.NewStyle().MaxHeight(7).Height(7).Render(m.header.View())
	var command string
	if m.showCommandView {
		command = lipgloss.NewStyle().
			Width(m.width - 2).
			PaddingLeft(2).
			PaddingRight(2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("2")).
			Render(m.command.View())
	}

	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(command) - 2

	m.content.WrapperHeight = contentHeight
	m.content.WrapperWidth = m.width - 2

	content := lipgloss.NewStyle().Height(contentHeight).Width(m.width).Render(m.content.View())

	return lipgloss.JoinVertical(lipgloss.Top, header, command, content)
}
