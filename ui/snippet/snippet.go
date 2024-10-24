package snippet


import tea "github.com/charmbracelet/bubbletea"

type Model struct {
    width, height int
}

func New() (m Model) {
    return m
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmds []tea.Cmd

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    return "model"
}
