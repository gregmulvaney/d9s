package volumes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

var Keymap = [][]string{}

type Model struct {
	ctx *appstate.State

	table table.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

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
	return "volumes"
}
