package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/constants"
	"github.com/gregmulvaney/d9s/pkg/data"
	"github.com/gregmulvaney/d9s/pkg/ui/command"
	"github.com/gregmulvaney/d9s/pkg/ui/content"
	"github.com/gregmulvaney/d9s/pkg/ui/header"
	"github.com/gregmulvaney/d9s/pkg/ui/splash"
)

type sessionState int

const (
	mainContext sessionState = iota
	commandContext
)

type Model struct {
	ctx appcontext.Context

	splash  splash.Model
	header  header.Model
	command command.Model
	content content.Model

	ready bool
}

func New() (m Model) {
	m.ready = false
	m.ctx = appcontext.Context{}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.SyncWindowSize(msg)
	}

	if !m.ready {
		m.ctx.Docker = data.NewDockerClient()
		m.ready = true
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Splash"
	}
	header := constants.HeaderStyle.Render(m.header.View())
	return lipgloss.JoinVertical(lipgloss.Top, header)
}

func (m *Model) SyncWindowSize(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
}
