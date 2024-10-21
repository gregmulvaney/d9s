package containers

import (
	"context"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
)

type sessionState int

const (
	containerOverview sessionState = iota
	containerLogs
)

type Model struct {
	width, height  int
	state          sessionState
	containersMaps []map[string]interface{}
	table          table.Model
	dockerClient   *docker.Client
	logs           string
}

func New(dockerClient *docker.Client) Model {
	containers, _ := dockerClient.ContainerList(context.Background(), container.ListOptions{})

	containersMaps := make([]map[string]interface{}, 1, len(containers))

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "NAME", Width: 15},
		{Title: "STATE", Width: 20},
		{Title: "STATUS", Width: 20},
		{Title: "IMAGE", Width: 50},
	}

	var rows []table.Row
	for index, ctr := range containers {
		ctrMap := map[string]interface{}{
			"Name": strings.Trim(ctr.Names[0], "/"),
			"ID":   ctr.ID,
		}
		containersMaps = append(containersMaps, ctrMap)
		i := strconv.Itoa(index)
		rows = append(rows, table.Row{i, strings.Trim(ctr.Names[0], "/"), ctr.State, ctr.Status, ctr.Image})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#74c7ec"))
	t.SetStyles(s)

	return Model{
		state:          containerOverview,
		containersMaps: containersMaps,
		table:          t,
		dockerClient:   dockerClient,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
        case "esc":
            m.logs = ""
            m.state = containerOverview
		case "L":
            m.state = containerLogs
			name := m.table.SelectedRow()[1]
			i := slices.IndexFunc(m.containersMaps, func(ctr map[string]interface{}) bool {
				return ctr["Name"] == name
			})
			ctr := m.containersMaps[i]
			m.getContainerLogs(ctr["ID"].(string))

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case containerLogs:
		return m.logs
	default:
		return m.table.View()
	}
}

func (m *Model) getContainerLogs(containerID string) {

	reader, err := m.dockerClient.ContainerLogs(context.Background(), containerID, container.LogsOptions{
        ShowStdout: true,
    })
	if err != nil {
		log.Fatal(err)
	}
    buf := new(strings.Builder)
    _, err = io.Copy(buf, reader)
    m.logs = buf.String()
}
