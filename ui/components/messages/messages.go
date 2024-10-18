package messages

import tea "github.com/charmbracelet/bubbletea"

type CommandMsg string

func SendCommand(command string) tea.Cmd {
    return func() tea.Msg {
        return CommandMsg(command)
    }
}
