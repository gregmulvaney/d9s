package command

import (
	"slices"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var commands = []string{
	"containers",
	"images",
	"networks",
	"volumes",
	"q",
}

type Model struct {
	input textinput.Model
}

type CommadMsg string

func Execute(command string) tea.Cmd {
	return func() tea.Msg {
		return CommadMsg(command)
	}
}

func New() (m Model) {
	m.input = textinput.New()
	m.input.SetSuggestions(commands)
	m.input.ShowSuggestions = true

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			value := m.input.Value()
			if !slices.Contains(commands, value) {
				return m, tea.WindowSize()
			}
			if value == "q" {
				return m, tea.Quit
			}
			m.input.Reset()
			return m, Execute(value)
		}
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Focus() {
	m.input.Focus()
	m.input.Cursor.SetMode(cursor.CursorBlink)
}

func (m *Model) Reset() {
	m.input.Reset()
}

func (m Model) View() string {
	return m.input.View()
}
