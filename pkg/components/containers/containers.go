package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx   *appcontext.Context
	table table.Model
}

var Keymap = [][]string{
	{"<C+r>", "Restart container"},
	{"<C+k>", "Stop Container"},
	{"<C+s>", "Start Container"},
	{"<l>", "Logs"},
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 5},
		{Title: "NAME", Width: 20},
		{Title: "STATE", Width: 12},
		{Title: "CPU", Width: 10},
		{Title: "MEM/LIMIT", Width: 20},
		{Title: "MEM", Width: 10},
		{Title: "STATUS", Width: 35},
		{Title: "IMAGE", Flex: true},
	}

	containers := m.queryContainersList()

	var rows = make([]table.Row, 0, len(containers))
	for i, container := range containers {
		index := strconv.Itoa(i)
		statsReader, _ := m.ctx.DockerClient.ContainerStatsOneShot(context.Background(), container.ID)
		defer statsReader.Body.Close()
		var stats types.Stats
		json.NewDecoder(statsReader.Body).Decode(&stats)
		// cpu := strconv.Itoa(int(stats.CPUStats.CPUUsage.TotalUsage))
		memlimit := fmt.Sprintf("%d/%d", (stats.MemoryStats.Usage/1024)/1024, (stats.MemoryStats.Limit/1024)/1024)
		memFloat := (float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)) * 100
		mem := fmt.Sprintf("%s%s", strconv.FormatFloat(memFloat, 'f', 2, 64), "%")
		name := strings.Trim(container.Names[0], "/")
		row := table.Row{container.ID, index, name, container.State, "0%", memlimit, mem, container.Status, container.Image}
		rows = append(rows, row)
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...,
	)
}

func (m Model) View() string {
	return m.table.View()
}

func (m *Model) SyncAppContext(ctx *appcontext.Context) {
	m.ctx = ctx
	m.table.UpdateDimensions(m.ctx.ScreenWidth-3, m.ctx.MainContentHeight)
}

func (m *Model) queryContainersList() []types.Container {
	containers, err := m.ctx.DockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}
	return containers
}
