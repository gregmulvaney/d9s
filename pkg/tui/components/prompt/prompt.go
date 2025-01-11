package prompt

import (
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/commands"
)

var modes = []string{
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

	// Init the textinput bubble and set suggstions
	m.input = textinput.New()
	m.input.ShowSuggestions = true
	m.input.SetSuggestions(modes)

	return m
}

func (m *Model) Focus() {
	m.input.Focus()
}

func (m *Model) Reset() {
	m.input.Reset()
}

func (m Model) Init() tea.Cmd {
	// TODO: Does this actually do anything since the prompt isn't rendered by default?
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
			// Quit the app
			if value == "q" {
				return m, tea.Quit
			}
			// Check if the value of the prompt exist in the valid modes array
			// TODO: add some kind of visual effect to communicate invalid commands
			if !slices.Contains(modes, value) {
				return m, tea.WindowSize()
			}
			m.input.Reset()
			return m, commands.PromptExecute(value)
		}
	}

	// Pass all key messages to the input bubble
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.input.View()
}
