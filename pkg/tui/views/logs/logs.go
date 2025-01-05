package logs

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

type Model struct {
	ctx *appstate.State

	viewport viewport.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx
	m.viewport = viewport.New(10, 20)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 12
		m.viewport.Width = msg.Width - 2
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
