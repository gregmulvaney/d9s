package keymap

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMapItem []string

type sessionState int

const (
	containersView sessionState = iota
	networksView
	volumesView
)

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

type Model struct {
	state         sessionState
	containerKeys []keyMapItem
}

func New() (m Model) {
	containerKeys := []keyMapItem{
		{"l", "Logs"},
		{"r", "Restart"},
		{"d", "Delete"},
		{"k", "Stops"},
	}
	m.state = containersView
	m.containerKeys = containerKeys
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		break
	}
	return m, nil
}

func (m Model) View() string {
    // var keyset1 string
    // var keyset2 string
    // var keyset3 string

	switch m.state {
	default:
        rows := make([]string, 0, len(m.containerKeys))
	    for i := 0; i < len(m.containerKeys); i++ {
            key := fmt.Sprintf("<%s>", m.containerKeys[i][0])
            rows = append(rows, fmt.Sprintf("%s %s", keyStyle.Render(key), m.containerKeys[i][1]))
        }
        return lipgloss.JoinVertical(lipgloss.Top, rows...)
	}
}

