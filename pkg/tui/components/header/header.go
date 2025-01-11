package header

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/keylist"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/commands"
	"github.com/gregmulvaney/d9s/pkg/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/views/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/views/images"
	"github.com/gregmulvaney/d9s/pkg/tui/views/networks"
	"github.com/gregmulvaney/d9s/pkg/tui/views/volumes"
)

var keymaps = map[commands.PromptMsg][][]string{
	"containers": containers.Keymap,
	"images":     images.Keymap,
	"networks":   networks.Keymap,
	"volumes":    volumes.Keymap,
}

type Model struct {
	// State
	ctx   *appstate.State
	width int

	// Elements
	status keylist.Model
	keymap keylist.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	clientVersion := m.ctx.DockerClient.ClientVersion()
	serverVersion, err := m.ctx.DockerClient.ServerVersion(context.Background())
	if err != nil {
		panic(err)
	}

	statusItems := [][]string{
		{"Host", os.Getenv("DOCKER_HOST")},
		{"Server", fmt.Sprintf("v%s", serverVersion.Version)},
		{"Client", fmt.Sprintf("v%s", clientVersion)},
	}

	m.status = keylist.New(
		keylist.WithItems(statusItems),
		keylist.WithGrid(true),
		keylist.WithSeparator(":"),
	)

	keymapStyles := keylist.Styles{
		Key:   lipgloss.NewStyle().Foreground(constants.KEYMAP_KEY_COLOR).Bold(true),
		Value: lipgloss.NewStyle().Foreground(constants.KEYMAP_VALUE_COLOR),
	}

	m.keymap = keylist.New(
		keylist.WithItems(keymaps["containers"]),
		keylist.WithGrid(true),
		keylist.WithMaxRows(6),
		keylist.WithStyles(keymapStyles),
	)

	m.keymap.SetItems(keymaps["containers"])

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Set keymap based on the PromptMsg
	case commands.PromptMsg:
		m.keymap.SetItems(keymaps[msg])

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	status := lipgloss.NewStyle().
		Width(m.width / 3).
		MaxWidth(m.width / 3).
		Render(m.status.View())

	keymap := lipgloss.NewStyle().
		Width(m.width / 3).
		MaxWidth(m.width / 3).
		Render(m.keymap.View())

	logo := lipgloss.NewStyle().
		Width(m.width / 3).
		MaxWidth(m.width / 3).
		Align(lipgloss.Right).
		Foreground(constants.LOGO_COLOR).
		Render(constants.LOGO_SMALL)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
