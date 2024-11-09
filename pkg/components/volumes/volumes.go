package volumes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx   *appcontext.Context
	table table.Model
}

func New(ctx *appcontext.Context) (m Model) {
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
	return "secrets"
}

func (m *Model) SyncProgramContext(ctx *appcontext.Context) {
	m.ctx = ctx
}
