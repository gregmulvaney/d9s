package images

import (
	"context"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/image"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx    *appcontext.Context
	table  table.Model
	images []image.Summary
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	images, err := m.ctx.Docker.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	m.images = images

	rows := renderRows(images)

	cols := []table.Column{
		{Title: "ID", Hidden: true},
		{Title: "#", Width: 4},
		{Title: "Name", Flex: true},
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocus(true),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func renderRows(images []image.Summary) []table.Row {
	var rows []table.Row

	for i, image := range images {
		id := image.ID
		index := strconv.Itoa(i)
		var name string
		if len(image.RepoTags) > 0 {
			name = image.RepoTags[0]
		} else {
			name = "no name"
		}

		rows = append(rows, table.Row{id, index, name})
	}
	return rows
}
