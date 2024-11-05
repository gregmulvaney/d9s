package volumes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type Model struct {
	ctx   *context.ProgramContext
	table table.Model
}

func New(ctx *context.ProgramContext) (m Model) {
	m.ctx = ctx

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
}

func (m Model) View() string {
	return "volumes"
}
