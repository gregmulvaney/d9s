package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/ui/components/content"
)

var Commands = []string{
	"networks",
	"containers",
	"volumes",
}

type FocusedMsg bool

func Focused() tea.Cmd {
	return func() tea.Msg {
		return FocusedMsg(true)
	}
}

type Model struct {
	width, height int
	input         textinput.Model
}

func New() Model {
	input := textinput.New()
	input.ShowSuggestions = true
	input.SetSuggestions(Commands)
	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
        case tea.KeyEsc.String():
            m.input.SetValue("")
		case tea.KeyEnter.String():
            // TODO: Validate command
            command := m.input.Value()
            if command == "q" {
                return m, tea.Quit
            }
            m.input.SetValue("")
            return m, content.Command(command)
		default:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}

	case FocusedMsg:
		m.input.Focus()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}
