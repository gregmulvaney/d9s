package header

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/d9s/pkg/appcontext"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

// Using a struct because maps render weird
type statusField struct {
	Key   string
	Value string
}

type Model struct {
	ctx          *appcontext.Context
	statusFields []statusField
}

func New(ctx *appcontext.Context) (m Model) {
	m.ctx = ctx
	host := strings.Split(os.Getenv("DOCKER_HOST"), "//")
	engineVer, err := m.ctx.DockerClient.ServerVersion(context.Background())
	if err != nil {
		panic(err)
	}
	clientVer := m.ctx.DockerClient.ClientVersion()

	m.statusFields = []statusField{
		{Key: "Host:", Value: host[1]},
		{Key: "Engine:", Value: fmt.Sprintf("v%s", engineVer.Version)},
		{Key: "Client:", Value: fmt.Sprintf("v%s", clientVer)},
		{Key: "CPU:", Value: "N/A"},
		{Key: "Mem:", Value: "N/A"},
	}

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var statusFields []string
	for _, field := range m.statusFields {
		key := lipgloss.NewStyle().Bold(true).Foreground(constants.STATUSCOLOR).PaddingRight(1).Render(field.Key)
		statusFields = append(statusFields, lipgloss.JoinHorizontal(lipgloss.Left, key, field.Value))
	}

	status := lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth / 3).
		MaxWidth(m.ctx.ScreenWidth / 3).
		Align(lipgloss.Left).
		Render(lipgloss.JoinVertical(lipgloss.Top, statusFields...))

	keymap := lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth / 3).
		MaxWidth(m.ctx.ScreenWidth / 3).
		Align(lipgloss.Center).
		Render("keymap")

	logo := lipgloss.NewStyle().
		Width(m.ctx.ScreenWidth / 3).
		MaxWidth(m.ctx.ScreenWidth / 3).
		Align(lipgloss.Right).
		Bold(true).
		Foreground(constants.LOGOCOLOR).
		Render(constants.RAWLOGO)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}

func (m *Model) SyncAppContext(ctx *appcontext.Context) {
	m.ctx = ctx
}
