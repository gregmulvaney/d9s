package constants

import "github.com/charmbracelet/lipgloss"

const (
	PROMPT_BORDER_COLOR  = lipgloss.Color("2")
	CONTENT_BORDER_COLOR = lipgloss.Color("4")
	LOGO_COLOR           = lipgloss.Color("214")
	KEYMAP_KEY_COLOR     = lipgloss.Color("13")
	KEYMAP_VALUE_COLOR   = lipgloss.Color("7")
)

var HeaderStyle = lipgloss.NewStyle()

var PromptStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(PROMPT_BORDER_COLOR)

var ContentStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(CONTENT_BORDER_COLOR)

var CrumbsStyle = lipgloss.NewStyle()
