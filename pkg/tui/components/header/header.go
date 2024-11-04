package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

var logoRaw = `________  ________       
 \______ \/   __   \______
  |    |  \____    /  ___/
  |    '   \ /    /\___ \ 
 /_______  //____//____  >
         \/            \//
    `

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

func (m *Model) UpdateProgramContenxt(ctx *context.ProgramContext) {
	m.ctx = ctx
}

func (m Model) View() string {
	status := lipgloss.NewStyle().Align(lipgloss.Left).Width(m.ctx.ScreenWidth / 3).Render("status")
	keymap := lipgloss.NewStyle().Width(m.ctx.ScreenWidth / 3).Render("keymap")
	logo := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.ctx.ScreenWidth / 3).Foreground(lipgloss.Color("3")).Render(logoRaw)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
