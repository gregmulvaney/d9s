package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var Commands = []string{
	"networks",
	"containers",
	"volumes",
}

type Model struct {
	width, height int
	input         textinput.Model
}

func New() (m Model) {
	m.input = textinput.New()
	m.input.ShowSuggestions = true
	m.input.SetSuggestions(Commands)
	return m
}

type CommandMsg string

func Command(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
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

		case tea.KeyEnter.String():
			command := m.input.Value()
			// Quit command
			if command == "q" {
				return m, tea.Quit
			}
			// TODO: Validate command
			m.Clear()
			return m, Command(command)

		case ":":
			// Ignore colons in commands
			break

		default:
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

func (m *Model) Clear() {
	m.input.SetValue("")
}

func (m *Model) Focus() {
	m.input.Focus()
}
