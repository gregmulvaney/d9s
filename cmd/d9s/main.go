package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/ui"
)

func main() {
	program := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
