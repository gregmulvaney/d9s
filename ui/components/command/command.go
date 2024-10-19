package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/ui/components/content"
)

type FocusMsg bool

func Focus() tea.Cmd {
	return func() tea.Msg {
		return FocusMsg(true)
	}
}

type ClearMsg bool

func Clear() tea.Cmd {
    return func() tea.Msg {
        return ClearMsg(true)
    }
}

type Model struct {
	width, height int
	input         textinput.Model
}

func New() Model {
	input := textinput.New()
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
			default:
                command := m.input.Value()
                m.input.SetValue("")
				return m, content.Command(command)
			}
		default:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case FocusMsg:
		m.input.Focus()
    case ClearMsg:
        m.input.SetValue("`")
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}

func (m Model) Focus() {
	m.input.Focus()
}
