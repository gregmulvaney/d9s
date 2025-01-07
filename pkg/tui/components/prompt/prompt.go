package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
)

var commands = []string{
	"containers",
	"images",
	"networks",
	"volumes",
}

type Model struct {
	// State
	ctx *appstate.State

	// Elements
	input textinput.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	m.input = textinput.New()
	m.input.ShowSuggestions = true
	m.input.SetSuggestions(commands)

	return m
}

func (m *Model) Focus() {
	m.input.Focus()
}

func (m *Model) Reset() {
	m.input.Reset()
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
			// TODO
		}
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}
