package status

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	docker "github.com/docker/docker/client"
)

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Width(7)

type Model struct {
	fields map[string]interface{}
}

func NewModel(dockerClient *docker.Client) Model {
	serverVersion, _ := dockerClient.ServerVersion(context.Background())
	clientApiVersion := dockerClient.ClientVersion()
	fields := map[string]interface{}{
		"Host":   os.Getenv("DOCKER_HOST"),
		"Engine": fmt.Sprintf("v%s", serverVersion.Version),
		"Client": fmt.Sprintf("v%s", clientApiVersion),
	}
	return Model{
		fields: fields,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var s string

    s += fmt.Sprintf("%s %s\n", keyStyle.Render("Host:"), m.fields["Host"] )
    s += fmt.Sprintf("%s %s\n", keyStyle.Render("Engine:"), m.fields["Engine"])
    s += fmt.Sprintf("%s %s\n", keyStyle.Render("Client:"), m.fields["Client"] )
	return s
}
