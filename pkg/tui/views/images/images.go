package images

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/image"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

var Keymap = [][]string{
	{"<C+p>", "Prune"},
}

type (
	fetchImageMsg bool
)

type Model struct {
	ctx *appstate.State

	table table.Model
}

func FetchImages() tea.Cmd {
	return func() tea.Msg {
		return fetchImageMsg(true)
	}
}

func pruneImages(ctx *appstate.State) tea.Cmd {
	return func() tea.Msg {
		_, err := ctx.Api.ImagesPrune()
		if err != nil {
			return constants.ApiErrorMsg(err)
		}
		return fetchImageMsg(true)
	}
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	cols := []table.Column{
		{Title: "#", Width: 5},
		{Title: "ID", Flex: true},
		{Title: "Tags", Flex: true},
		{Title: "Size", Width: 20},
		{Title: "Unused", Width: 8},
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithFocus(true),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlP:
			return m, pruneImages(m.ctx)
		}

	case constants.InitMsg:
		return m, FetchImages()

	case fetchImageMsg:
		images, err := m.ctx.Api.ImagesFetch()
		if err != nil {
			panic(err)
		}
		m.table.SetRows(renderRows(images))
		return m, tea.WindowSize()

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 12)
		m.table.SetWidth(msg.Width - 2)
	}

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
		index := strconv.Itoa(i)
		id := image.ID
		tags := strings.Join(image.RepoTags, ",")
		size := strconv.Itoa(int(image.Size))
		// TODO: Why is this -1 and what the fuck are the docs talking about
		unused := false
		if image.Containers == 0 {
			unused = true
		}
		rows = append(rows, table.Row{index, id, tags, size, strconv.FormatBool(unused)})
	}
	return rows
}
