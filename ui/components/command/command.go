package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	height, width int
	input         textinput.Model
}

func NewModel() Model {
	input := textinput.New()
	input.PromptStyle.Border(lipgloss.NormalBorder(), true).Foreground(lipgloss.Color("#74c7ec"))
	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
        case "enter":
            switch m.input.Value() {
            case "q":
                return m, tea.Quit
            }
	    case "esc":
            m.input.SetValue("")
		case ":":
			m.input.Focus()
            if m.input.Focused() {
                m.input.SetValue("")
            }
			m.input.Width = 100
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().
        Width(m.width).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("#74c7ec")).
        PaddingLeft(2).
        PaddingRight(2).
        Render(m.input.View())
}
