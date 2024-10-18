package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input textinput.Model
}

func New() Model {
	input := textinput.New()
	input.Width = 20
	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.input.Value() {
			case "q":
				return m, tea.Quit
            case "Networks":
			}
		case "esc":
			m.input.SetValue("")
		case ":":
			m.input.Focus()
		default:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, nil
}

func (m Model) View() string {
	return m.input.View()
}
