package command

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx   *appcontext.Context
	input textinput.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	m.input = textinput.New()
	return m
}

// TODO: Probably redundant
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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
