package header

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/keylist"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/constants"
	"github.com/gregmulvaney/d9s/pkg/tui/components/command"
	"github.com/gregmulvaney/d9s/pkg/tui/views/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/views/images"
	"github.com/gregmulvaney/d9s/pkg/tui/views/networks"
)

var keymaps = map[command.CommadMsg][][]string{
	"containers": containers.Keymap,
	"images":     images.Keymap,
	"networks":   networks.Keymap,
}

type Model struct {
	ctx   *appstate.State
	width int

	status keylist.Model
	keymap keylist.Model
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	clientVersion := m.ctx.Api.Client.ClientVersion()
	serverVersion, err := m.ctx.Api.Client.ServerVersion(context.Background())
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

	m.keymap = keylist.New(
		keylist.WithItems(keymaps["containers"]),
		keylist.WithGrid(true),
		keylist.WithMaxRows(6),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case command.CommadMsg:
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
		Foreground(lipgloss.Color(lipgloss.Color("214"))).
		Align(lipgloss.Right).
		Render(constants.LOGO_SMALL)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
