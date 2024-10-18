package logo

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387"))

type Model struct {}

func NewModel() Model {
    return Model{}
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil 
}

func (m Model) View() string {
   var s string
    s +=`________    ________       
\______ \  /  _____/ ______
 |    |  \/   __  \ /  ___/
 |    '   \  |__\  \\___ \ 
/_______  /\_____  /____  >
        \/       \/     \/ `

    return logoStyle.Render(s) 
}
