package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/tui/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

type CommandMsg string

type Model struct {
	ctx   *context.ProgramContext
	input textinput.Model
}

func New(ctx context.ProgramContext) (m Model) {
	m.ctx = &ctx
	m.input = textinput.New()
	m.input.ShowSuggestions = true
	m.input.SetSuggestions(constants.Commands)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			command := m.input.Value()
			if command == "q" {
				return m, tea.Quit
			}
			Execute(command)
			m.Reset()
		default:
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
}

func (m *Model) Focus() {
	m.input.Focus()
}

func (m *Model) Reset() {
	m.input.Reset()
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("2")).
		Width(m.ctx.ScreenWidth - 2).
		PaddingLeft(1).
		Render(m.input.View())
}

func Execute(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandMsg(command)
	}
}
