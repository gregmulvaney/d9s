package splash

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

func Ready() tea.Cmd {
	return func() tea.Msg {
		return constants.ReadyMsg(true)
	}
}

type Model struct {
	ctx *appcontext.Context
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	logo := lipgloss.NewStyle().Foreground(constants.LOGOCOLOR).Bold(true).Render(constants.RAWLOGO)
	splash := lipgloss.JoinVertical(lipgloss.Top, logo)

	return lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth).
		Height(m.ctx.ScreenHeight).
		AlignHorizontal(lipgloss.Center).
		PaddingTop(10).
		Render(splash)
}

func (m *Model) SyncAppContext(ctx *appcontext.Context) {
	m.ctx = ctx
}
