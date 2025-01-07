package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

type Model struct {
	ctx *appstate.State
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "content"
}
