package constants

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	InitMsg     bool
	LogsMsg     *io.Reader
	ApiErrorMsg error
)

func ApiError(err error) tea.Cmd {
	return func() tea.Msg {
		return ApiErrorMsg(err)
	}
}

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
