package header

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
}

func New() (m Model) {
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "header"
}
