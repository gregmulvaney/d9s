package content

import tea "github.com/charmbracelet/bubbletea"

type sessionContenxt int

const (
	containersContext sessionContenxt = iota
	volumesContext
	secretsContext
	networksContext
)

type Model struct {
}

func New() (m Model) {
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "content"
}
