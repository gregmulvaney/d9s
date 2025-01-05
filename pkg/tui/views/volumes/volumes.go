package volumes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/volume"
	"github.com/gregmulvaney/bubbles/table"
	"github.com/gregmulvaney/d9s/pkg/appstate"
	"github.com/gregmulvaney/d9s/pkg/constants"
)

type (
	fetchVolumesMsg bool
	apiErrorMsg     error
)

type Model struct {
	ctx   *appstate.State
	table table.Model
}

func fetchVolumes() tea.Cmd {
	return func() tea.Msg {
		return fetchVolumesMsg(true)
	}
}

func apiError(err error) tea.Cmd {
	return func() tea.Msg {
		return apiErrorMsg(err)
	}
}

func New(ctx *appstate.State) (m Model) {
	m.ctx = ctx

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlP:
			//
		}
	case fetchVolumesMsg:
		volumes, err := m.ctx.Api.VolumesFetch()
		if err != nil {
			return m, apiErrorMsg(err)
		}
		m.table.SetRows(renderRows(volumes))
		return m, tea.WindowSize()
	case constants.InitMsg:
		return m, tea.Batch(tea.WindowSize(), fetchVolumes())

	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 12)
		m.table.SetWidth(msg.Width - 2)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func renderRows(volumes []*volume.Volume) []table.Row {
	var rows []table.Row

	return rows
}
