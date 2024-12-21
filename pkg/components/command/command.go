package command

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

type CommandMsg string

type Model struct {
	height int
	width  int
	input  textinput.Model
}

func New() (m Model) {
	m.height = 0
	m.width = 0

	m.input = textinput.New()
	m.input.SetSuggestions(constants.Commands)
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
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) Focus() {
	m.input.Focus()
	m.input.Cursor.SetMode(cursor.CursorBlink)
}

func (m *Model) Blur() {
	m.input.Blur()
}

func (m *Model) Reset() {
	m.input.Reset()
}

func (m Model) View() string {
	return m.input.View()
}

func Execute(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
	}
}
