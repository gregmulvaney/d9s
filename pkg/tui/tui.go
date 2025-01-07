package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	docker "github.com/docker/docker/client"

	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/breadcrumbs"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/commands"
	"github.com/gregmulvaney/d9s/pkg/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/components/content"
	"github.com/gregmulvaney/d9s/pkg/tui/components/header"
	"github.com/gregmulvaney/d9s/pkg/tui/components/prompt"
)

type sessionState int

const (
	contentMode sessionState = iota
	promptMode
)

type Model struct {
	// State values
	ctx        appstate.State
	state      sessionState
	height     int
	width      int
	showPrompt bool

	// Components
	header  header.Model
	prompt  prompt.Model
	content content.Model
	crumbs  breadcrumbs.Model
}

func New() (m Model) {
	// Initialize app state
	m.ctx = appstate.State{}
	// Set initial focus state
	m.state = contentMode

	// Connect to the remote daemon and initiliaze the Client
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}
	m.ctx.DockerClient = dockerClient

	// Initialize all elements
	m.header = header.New(&m.ctx)
	m.prompt = prompt.New(&m.ctx)
	m.content = content.New(&m.ctx)

	// Set initial breadcrumbs state
	crumbs := []string{"Containers"}
	m.crumbs = breadcrumbs.New(breadcrumbs.WithCrumbs(crumbs))

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.prompt.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Handle key messages
	case tea.KeyMsg:
		// Handle keys with a known type
		switch msg.Type {
		// Quit the app
		case tea.KeyCtrlC:
			return m, tea.Quit

		// Handle the escape key
		case tea.KeyEsc:
			// Exit the prompt if it is shown
			if m.showPrompt {
				m.showPrompt = false
				m.state = contentMode
				m.prompt.Reset()
				return m, tea.WindowSize()
			}
			// TODO: Handle navigating backwards from subviews
		}

		// Handle keys by string
		switch msg.String() {
		// Show prompt and focus the input
		case ":":
			if !m.showPrompt {
				m.showPrompt = true
				m.state = promptMode
				m.prompt.Focus()
				return m, tea.WindowSize()
			}
		}

	case commands.PromptMsg:
		// Pass the PromptMsg to the content element
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
		m.showPrompt = false
		m.state = contentMode

	// Track window size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Route messages based on active state
	switch m.state {
	case promptMode:
		m.prompt, cmd = m.prompt.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.content, cmd = m.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Pass all messages to header component
	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := constants.HeaderStyle.
		Width(m.width).
		MaxWidth(m.width).
		Render("header")

	// Measure the height of the header element
	_, headerHeight := lipgloss.Size(header)

	var prompt string
	promptHeight := 0

	// If prompt is shown render the prompt element and measure its height
	if m.showPrompt {
		prompt = constants.PromptStyle.
			Width(m.width - 3).
			PaddingLeft(1).
			Render(m.prompt.View())
		_, promptHeight = lipgloss.Size(prompt)
	}

	crumbs := constants.CrumbsStyle.
		Width(m.width).
		MaxWidth(m.width).
		Render(m.crumbs.View())

	// Measure the height of the crumbs element
	_, crumbsHeight := lipgloss.Size(crumbs)

	// Subtract the rendered elements height from the window height and an additional 2 units for the border
	contentHeight := m.height - headerHeight - promptHeight - crumbsHeight - 2
	m.ctx.ContentHeight = contentHeight

	content := constants.ContentStyle.
		Width(m.width - 2).
		Height(contentHeight).
		Render(m.content.View())

	// Hack to circumvent blank elements still taking 1 unit of height
	if m.showPrompt {
		return lipgloss.JoinVertical(lipgloss.Top, header, prompt, content, crumbs)
	}
	return lipgloss.JoinVertical(lipgloss.Top, header, content, crumbs)
}
