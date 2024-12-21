package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/components/logo"
	"github.com/gregmulvaney/d9s/pkg/components/status"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

type Model struct {
	ctx    *appcontext.Context
	height int
	width  int
	status status.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	m.height = 0
	m.width = 0
	m.status = status.New(ctx)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	return m, nil
}

func (m Model) View() string {
	status := lipgloss.NewStyle().Width(m.width / 3).MaxWidth(m.width / 3).Render(m.status.View())

	keymap := lipgloss.NewStyle().Width(m.width / 3).MaxWidth(m.width / 3).Render("keymap")

	logo := lipgloss.NewStyle().
		Width(m.width / 3).
		MaxWidth(m.width / 3).
		Foreground(lipgloss.Color(constants.LOGOCOLOR)).
		Align(lipgloss.Right).
		Render(logo.LOGO_SMALL)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
