package containers

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

var Keymap = [][]string{
	{"<C+s>", "Start"},
	{"<C+k>", "Stop"},
	{"<C+r>", "Restart"},
	{"<C+j>", "Delete"},
	{"<C+p>", "Pause"},
	{"<C+u>", "Unpause"},
	{"<l>", "Logs"},
}

type (
	fetchContainersMsg bool
	apiErrorMsg        error
)

type Model struct {
	ctx   *appstate.State
	table table.Model
}

func apiError(err error) tea.Cmd {
	return func() tea.Msg {
		return apiErrorMsg(err)
	}
}

func fetchContainers(ctx *appstate.State) ([]types.Container, error) {
	containers, err := ctx.Api.FetchContainers()
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func initContainers() tea.Cmd {
	return func() tea.Msg {
		return fetchContainersMsg(true)
	}
}

func startContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.Api.StartContainer(id)
		if err != nil {
			return apiErrorMsg(err)
		}

		return fetchContainersMsg(true)
	}
}

func stopContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.Api.StopContainer(id)
		if err != nil {
			return apiErrorMsg(err)
		}
		return fetchContainersMsg(true)
	}
}

func restartContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.Api.RestartContainer(id)
		if err != nil {
			return apiErrorMsg(err)
		}
		return fetchContainersMsg(true)
	}
}

func removeContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.Api.RemoveContainer(id)
		if err != nil {
			return apiErrorMsg(err)
		}
		return fetchContainersMsg(true)
	}
}

func pauseContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		return fetchContainersMsg(true)
	}
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "State", Width: 12},
		{Title: "Image", Flex: true},
		{Title: "IP Address", Flex: true},
		{Title: "Ports", Flex: true},
	}

	rows := []table.Row{}
	m.table = table.New(
		table.WithColumns(cols),
		table.WithFocus(true),
		table.WithRows(rows),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return initContainers()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		selected := m.table.SelectedRow()
		switch msg.Type {
		case tea.KeyCtrlK:
			return m, stopContainer(m.ctx, selected[0])
		case tea.KeyCtrlS:
			return m, startContainer(m.ctx, selected[0])
		case tea.KeyCtrlR:
			return m, restartContainer(m.ctx, selected[0])
		case tea.KeyCtrlJ:
			return m, removeContainer(m.ctx, selected[0])
		case tea.KeyCtrlP:
			return m, pauseContainer(m.ctx, selected[0])
		}

	case apiErrorMsg:
		panic(msg)

	case fetchContainersMsg:
		containers, err := fetchContainers(m.ctx)
		if err != nil {
			return m, apiError(err)
		}
		rows := renderRows(containers)
		m.table.SetRows(rows)
		return m, tea.WindowSize()

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 12)
		m.table.SetWidth(msg.Width - 2)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, nil
}

func (m Model) View() string {
	return m.table.View()
}

func renderRows(containers []types.Container) []table.Row {
	var rows []table.Row
	for i, container := range containers {
		index := strconv.Itoa(i)
		name := strings.Trim(container.Names[0], "/")

		// TODO: Look into this more
		networks := container.NetworkSettings.Networks
		var ip string
		for _, network := range networks {
			ip = network.IPAddress
		}

		var ports []string
		for _, port := range container.Ports {
			if port.PublicPort == 0 {
				continue
			}
			ports = append(ports, fmt.Sprintf("%v:%v", port.PrivatePort, port.PublicPort))
		}

		// TODO: Remove duplicated ports
		var portsList string
		if len(ports) == 0 {
			portsList = "-"
		} else {
			portsList = strings.Join(ports, ",")
		}

		row := table.Row{
			container.ID,
			index,
			name,
			container.State,
			container.Image,
			ip,
			portsList,
		}

		rows = append(rows, row)
	}
	return rows
}
