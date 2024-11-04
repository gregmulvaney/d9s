package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type Model struct {
	ctx *context.ProgramContext
}

func New(ctx context.ProgramContext) (m Model) {
	m.ctx = &ctx
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth - 2).
		MaxWidth(m.ctx.ScreenWidth).
		Height(m.ctx.MainContentHeight).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("4")).
		Render("content")
}
