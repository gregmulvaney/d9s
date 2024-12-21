package status

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/keylist"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
)

type Model struct {
	ctx *appcontext.Context

	items keylist.Model
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	clientVersion := m.ctx.Docker.ClientVersion()
	host := os.Getenv("DOCKER_HOST")
	server, _ := m.ctx.Docker.ServerVersion(context.Background())
	itemsList := [][]string{
		{"Host", host},
		{"Server", fmt.Sprintf("v%s", server.Version)},
		{"Client", fmt.Sprintf("v%s", clientVersion)},
	}

	styles := keylist.Styles{
		Key: lipgloss.NewStyle().PaddingRight(0).Bold(true).Foreground(lipgloss.Color("214")),
	}

	m.items = keylist.New(
		keylist.WithItems(itemsList),
		keylist.WithGrid(true),
		keylist.WithSeparator(":"),
		keylist.WithStyles(styles),
	)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.items.View()
}
