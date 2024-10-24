package header

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/components/keymap"
)

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Width(7)

var logoRaw = `________  ________       
 \______ \/   __   \______
  |    |  \____    /  ___/
  |    '   \ /    /\___ \ 
 /_______  //____//____  >
         \/            \//
    `

type Model struct {
	width, height int
	keymap        keymap.Model
	statusFields  [][]string
}

func New(dockerClient *docker.Client) (m Model) {
	serverVersion, _ := dockerClient.ServerVersion(context.Background())
	clientApiVersion := dockerClient.ClientVersion()

	m.statusFields = [][]string{
		{"Host:", os.Getenv("DOCKER_HOST")},
		{"Server:", fmt.Sprintf("v%s", serverVersion.Version)},
		{"Client:", fmt.Sprintf("v%s", clientApiVersion)},
	}
	m.keymap = keymap.New()

	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.keymap, cmd = m.keymap.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var statusFields string
	for _, field := range m.statusFields {
		statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render(field[0]), field[1])
	}

	status := lipgloss.NewStyle().Width(m.width / 3).Align(lipgloss.Left).Render(statusFields)
	keymap := lipgloss.NewStyle().Width(m.width / 3).Align(lipgloss.Center).Render(m.keymap.View())
	logo := lipgloss.NewStyle().Width(m.width / 3).Align(lipgloss.Right).Foreground(lipgloss.Color("11")).Render(logoRaw)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
