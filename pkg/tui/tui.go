package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
	docker "github.com/docker/docker/client"
	"github.com/gregmulvaney/bubbles/breadcrumbs"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/components/command"
	"github.com/gregmulvaney/d9s/pkg/components/content"
	"github.com/gregmulvaney/d9s/pkg/components/header"
	"github.com/gregmulvaney/d9s/pkg/components/splash"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

type sessionState int

const (
	contentMode sessionState = iota
	splashMode
	commandMode
)

type Model struct {
	ctx   appcontext.Context
	state sessionState
	ready bool

	splash      splash.Model
	header      header.Model
	command     command.Model
	content     content.Model
	breadcrumbs breadcrumbs.Model
}

func New() (m Model) {
	m.ctx = appcontext.Context{}
	m.state = splashMode
	m.ready = false

	client, err := client.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}
	m.ctx.DockerClient = client

	m.splash = splash.New(&m.ctx)
	m.header = header.New(&m.ctx)
	m.command = command.New(&m.ctx)
	m.content = content.New(&m.ctx)
	crumbs := []string{"<containers>"}
	m.breadcrumbs = breadcrumbs.New(
		breadcrumbs.WithCrumbs(crumbs),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return constants.ReadyMsg(true)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ":":
			if !m.ctx.ShowCommandView {
				m.ctx.ShowCommandView = true
				m.state = commandMode
				m.command.Focus()
				cmds = append(cmds, tea.WindowSize(), textinput.Blink)
				return m, tea.Batch(cmds...)
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.ctx.ShowCommandView {
				m.state = contentMode
				m.ctx.ShowCommandView = false
				m.command.Reset()
				return m, tea.WindowSize()
			}
		}
	case constants.CommandMsg:
		m.ctx.ShowCommandView = false
		m.command.Reset()
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd, tea.WindowSize())
		return m, tea.Batch(cmds...)

	case constants.ReadyMsg:
		m.ready = true
		m.state = contentMode

	case tea.WindowSizeMsg:
		m.syncWindowSize(msg)
	}

	switch m.state {
	case splashMode:
		m.splash, cmd = m.splash.Update(msg)
		return m, cmd
	case commandMode:
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.syncAppContext()

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return m.splash.View()
	}

	header := lipgloss.NewStyle().Height(7).MaxHeight(constants.HEADERHEIGHT).Render(m.header.View())
	var command string
	if m.ctx.ShowCommandView {
		command = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(m.ctx.ScreenWidth - 3).
			MaxWidth(m.ctx.ScreenWidth).
			PaddingLeft(1).
			BorderForeground(lipgloss.Color("2")).
			Render(m.command.View())
	}

	content := m.content.View()

	return lipgloss.JoinVertical(lipgloss.Left, header, command, content, m.breadcrumbs.View())
}

func (m *Model) syncWindowSize(msg tea.WindowSizeMsg) {
	m.ctx.ScreenHeight = msg.Height
	m.ctx.ScreenWidth = msg.Width
	if m.ctx.ShowCommandView {
		m.ctx.MainContentHeight = m.ctx.ScreenHeight - constants.HEADERHEIGHT - constants.COMMANDHEIGHT - constants.CRUMBSHEIGHT
	} else {
		m.ctx.MainContentHeight = m.ctx.ScreenHeight - constants.HEADERHEIGHT - constants.CRUMBSHEIGHT
	}

	m.syncAppContext()
}

func (m *Model) syncAppContext() {
	m.splash.SyncAppContext(&m.ctx)
	m.header.SyncAppContext(&m.ctx)
	m.content.SyncAppContext(&m.ctx)
}
