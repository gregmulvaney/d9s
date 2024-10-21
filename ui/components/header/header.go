package header

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/d9s/ui/components/keymap"
	"os"
)

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Width(7)

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
	statusField   map[string]interface{}
}

func New(dockerClient *docker.Client) Model {
	serverVersion, _ := dockerClient.ServerVersion(context.Background())
	clientApiVersion := dockerClient.ClientVersion()

	statusField := make(map[string]interface{})
	statusField["host"] = os.Getenv("DOCKER_HOST")
	statusField["engine"] = serverVersion.Version
	statusField["client"] = clientApiVersion

	return Model{
		keymap:      keymap.New(),
		statusField: statusField,
	}
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
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Host:"), m.statusField["host"])
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Engine:"), m.statusField["engine"])
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Client:"), m.statusField["client"])
	status := lipgloss.NewStyle().Align(lipgloss.Left).Width(m.width / 3).Render(statusFields)
	keymap := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width / 3).Render(m.keymap.View())
	logo := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.width / 3).Foreground(lipgloss.Color("#f9e2af")).Render(logoRaw)
	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
