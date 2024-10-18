package header

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
)

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Width(7)

var logoRaw = `________    ________       
\______ \  /  _____/ ______
 |    |  \/   __  \ /  ___/
 |    '   \  |__\  \\___ \ 
/_______  /\_____  /____  >
        \/       \/     \/ `

type Model struct {
	width, height int
	dockerClient  *docker.Client
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
		dockerClient: dockerClient,
		statusField:  statusField,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	var statusFields string
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Host:"), m.statusField["host"])
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Engine:"), m.statusField["engine"])
	statusFields += fmt.Sprintf("%s %s\n", keyStyle.Render("Client:"), m.statusField["client"])
	status := lipgloss.NewStyle().Align(lipgloss.Left).Width(m.width / 3).Render(statusFields)
	keymap := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width / 3).Render("Keymap")
	logo := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.width / 3).Foreground(lipgloss.Color("#89b4fa")).Render(logoRaw)
	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
