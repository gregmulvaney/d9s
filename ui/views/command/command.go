package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var Commands = []string{
	"networks",
	"containers",
	"volumes",
	"q",
}

var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

type CommandMsg string

func Command(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
	}
}

type Model struct {
	width, height int
	input         textinput.Model
}

func New() (m Model) {
	m.input = textinput.New()
	m.input.SetSuggestions(Commands)
	m.input.ShowSuggestions = true
	return m
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
		// Handle command submission
		case tea.KeyEnter.String():
			command := m.input.Value()
			if err := validateCommand(command); err != nil {
				m.input.TextStyle = errStyle
				break
			}
			switch command {
			// Quit command
			case "q":
				return m, tea.Quit

			default:
				m.Clear()
				return m, Command(command)
			}

		// Ignore colons from command toggle
		case ":":
			break

		// Pass all other keys to the textinput
		default:
			m.input.TextStyle = lipgloss.NewStyle()
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}

// Focus the input model
func (m *Model) Focus() {
	m.input.Focus()
}

// Clear the input
func (m *Model) Clear() {
	m.input.Reset()
	m.input.Blur()
}

type validateError struct{}

func (e *validateError) Error() string {
	return "failed"
}

func validateCommand(command string) error {
	for _, valid := range Commands {
		if command == valid {
			return nil
		}
	}
	return &validateError{}
}
