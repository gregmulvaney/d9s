package images

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

var Keymap = [][]string{
	{"<C+p>", "Prune"},
}

type (
	fetchImageMsg bool
	apiError      error
)

type Model struct {
	ctx *appstate.State

	table table.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	cols := []table.Column{
		{Title: "#", Width: 4},
		{Title: "ID", Flex: true},
		{Title: "Tags", Flex: true},
		{Title: "Size", Width: 20},
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
			//
		}

	case fetchImageMsg:
		_, err := m.ctx.Api.ImagesFetch()
		if err != nil {
			panic(err)
		}

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
