package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/d9s/pkg/ui"
)

func main() {
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
