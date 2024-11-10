package command

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

type Model struct {
	ctx   *appcontext.Context
	input textinput.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	m.input = textinput.New()
	m.input.ShowSuggestions = true
	m.input.SetSuggestions(constants.Commands)
	return m
}

func (m Model) Init() tea.Cmd {
	// TODO: Probably redundant
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			command := m.input.Value()
			if command == "q" {
				return m, tea.Quit
			}
			if validateCommand(command) {
				return m, Execute(command)
			}
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

func validateCommand(command string) bool {
	for _, value := range constants.Commands {
		if value == command {
			return true
		}
	}
	return false
}

func Execute(command string) tea.Cmd {
	return func() tea.Msg {
		return constants.CommandMsg(command)
	}
}
