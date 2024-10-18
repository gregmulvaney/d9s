package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/components/command"
	"github.com/gregmulvaney/d9s/ui/components/data"
	"github.com/gregmulvaney/d9s/ui/components/logo"
	"github.com/gregmulvaney/d9s/ui/components/status"
)

type sessionState int

const (
	dataView sessionState = iota
	commandView
)

type Model struct {
	width, height   int
	state           sessionState
	status          status.Model
	logo            logo.Model
	command         command.Model
	showCommandView bool
	data            data.Model
	dockerClient    *docker.Client
}

func InitModel() Model {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()
	m := Model{
		state:  dataView,
		logo:   logo.NewModel(),
		status: status.NewModel(dockerClient),
	}
	m.data = data.NewModel(dockerClient)
	m.command = command.NewModel()
	m.showCommandView = false
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
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.showCommandView = false
			m.state = dataView
		case ":":
			if !m.showCommandView {
				m.showCommandView = true
				m.state = commandView
			}
		case "tab":
			if m.state == commandView {
				m.state = dataView
			} else {
				m.state = commandView
			}
		}
		switch m.state {
		case commandView:
			m.command, cmd = m.command.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.data, cmd = m.data.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	status := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(m.width / 3).
		Render(m.status.View())
	keymap := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width / 3).
		Render("Keymap")
	logo := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(m.width / 3).
        Height(7).
		Render(m.logo.View())
	header := lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
	var command string
	if m.showCommandView {
		command = lipgloss.NewStyle().Width(m.width).Render(m.command.View())
	}
	data := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - lipgloss.Height(logo)-4).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#0000FF")).
		Render(m.data.View())

	return lipgloss.JoinVertical(lipgloss.Top, header, command, data)
}
