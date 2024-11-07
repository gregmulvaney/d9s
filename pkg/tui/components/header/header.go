package header

import (
	"fmt"
	"os"

	dctx "context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/keylist"
	"github.com/gregmulvaney/d9s/pkg/tui/components/containers"
	"github.com/gregmulvaney/d9s/pkg/tui/components/content"
	"github.com/gregmulvaney/d9s/pkg/tui/context"
)

var logoRaw = `________  ________       
 \______ \/   __   \______
  |    |  \____    /  ___/
  |    '   \ /    /\___ \ 
 /_______  //____//____  >
         \/            \//
    `

type Model struct {
	ctx          *context.ProgramContext
	statusFields []statusField
	keymap       keylist.Model
}

type statusField struct {
	Key   string
	Value string
}

func New(ctx context.ProgramContext) (m Model) {
	m.ctx = &ctx

	engineVer, _ := m.ctx.DockerClient.ServerVersion(dctx.Background())
	clientVer := m.ctx.DockerClient.ClientVersion()

	host := statusField{Key: "Host:", Value: os.Getenv("DOCKER_HOST")}
	engine := statusField{Key: "Engine:", Value: fmt.Sprintf("v%s", engineVer.Version)}
	client := statusField{Key: "Client:", Value: fmt.Sprintf("v%s", clientVer)}
	m.statusFields = []statusField{
		host,
		engine,
		client,
	}

	m.keymap = keylist.New(
		keylist.WithItems(containers.Keymap),
	)

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case content.KeymapMsg:
		m.keymap.SetValues(msg)
	}
	return m, nil
}

func (m *Model) UpdateProgramContenxt(ctx *context.ProgramContext) {
	m.ctx = ctx
}

func (m Model) View() string {
	var statusFields []string
	for _, status := range m.statusFields {
		key := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).PaddingRight(1).Render(status.Key)
		statusFields = append(statusFields, lipgloss.JoinHorizontal(lipgloss.Left, key, status.Value))
	}
	status := lipgloss.NewStyle().Align(lipgloss.Left).Width(m.ctx.ScreenWidth / 3).Render(lipgloss.JoinVertical(lipgloss.Top, statusFields...))
	keymap := lipgloss.NewStyle().Width(m.ctx.ScreenWidth / 3).Height(6).MaxHeight(6).Render(m.keymap.View())
	logo := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.ctx.ScreenWidth / 3).Foreground(lipgloss.Color("214")).Render(logoRaw)

	return lipgloss.JoinHorizontal(lipgloss.Left, status, keymap, logo)
}
