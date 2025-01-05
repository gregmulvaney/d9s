package constants

import (
	tea "github.com/charmbracelet/bubbletea"
)

type InitMsg bool

func InitView() tea.Cmd {
	return func() tea.Msg {
		return InitMsg(true)
	}
}

const LOGO_SMALL = `________  ________       
 \______ \/   __   \______
  |    |  \____    /  ___/
  |    '   \ /    /\___ \ 
 /_______  //____//____  >
         \/            \//
    `
