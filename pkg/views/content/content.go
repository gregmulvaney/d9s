package content

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/views/containers"
)

type sessionState int

const (
	containersMode sessionState = iota
)

type Model struct {
	ctx        *appcontext.Context
	state      sessionState
	containers containers.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.state = containersMode

	m.containers = containers.New(ctx)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.state {
	default:
		m.containers, cmd = m.containers.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	default:
		return m.containers.View()
	}
}
