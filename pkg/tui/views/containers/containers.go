package containers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/commands"
)

var Keymap = [][]string{
	{"<C+s>", "Start"},
	{"<C+k>", "Stop"},
	{"<C+f>", "Kill"},
	{"<C+r>", "Restart"},
	{"<C+j>", "Delete"},
	{"<C+p>", "Pause"},
	{"<C+u>", "Unpause"},
	{"<l>", "Logs"},
}

type fetchContainersMsg []types.Container

type Model struct {
	ctx *appstate.State

	table table.Model
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
	return fetchContainers(m.ctx)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		selected := m.table.SelectedRow()
		id := selected[0]
		// Handle keys with defined types
		switch msg.Type {
		case tea.KeyCtrlS:
			return m, startContainer(m.ctx, id)
		case tea.KeyCtrlR:
			return m, restartContainer(m.ctx, id)
		case tea.KeyCtrlJ:
			return m, deleteContainer(m.ctx, id)
		case tea.KeyCtrlP:
			return m, pauseContainers(m.ctx, id)
		case tea.KeyCtrlU:
			return m, resumeContainers(m.ctx, id)
		}

	case fetchContainersMsg:
		m.table.SetRows(renderRows(msg))
		return m, tea.WindowSize()

	// Set the table dimensions
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width - 2)
		m.table.SetHeight(msg.Height - 13)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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

func fetchContainers(ctx *appstate.State) tea.Cmd {
	return func() tea.Msg {
		containers, err := ctx.DockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainersMsg(containers)
	}
}

func startContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerStart(context.Background(), id, container.StartOptions{})
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}

func stopContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerStop(context.Background(), id, container.StopOptions{})
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}

func killContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerKill(context.Background(), id, "SIGKILL")
		if err != nil {
			return commands.ApiError(err)
		}
		return fetchContainers(ctx)
	}
}

func restartContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerRestart(context.Background(), id, container.StopOptions{})
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}

func deleteContainer(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerRemove(context.Background(), id, container.RemoveOptions{})
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}

func pauseContainers(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerPause(context.Background(), id)
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}

func resumeContainers(ctx *appstate.State, id string) tea.Cmd {
	return func() tea.Msg {
		err := ctx.DockerClient.ContainerUnpause(context.Background(), id)
		if err != nil {
			return commands.ApiErrorMsg(err)
		}
		return fetchContainers(ctx)
	}
}
