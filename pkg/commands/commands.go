package commands

import tea "github.com/charmbracelet/bubbletea"

type (
	ApiErrorMsg error
	PromptMsg   string
)

func ApiError(err error) tea.Cmd {
	return func() tea.Msg {
		return ApiErrorMsg(err)
	}
}

func PromptExecute(command string) tea.Cmd {
	return func() tea.Msg {
		return PromptMsg(command)
	}
}
