package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit

		case key.Matches(msg, m.keys.toggle):
			m.showGlammy = !m.showGlammy
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := m.docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case DoneMsg:
		cmds := []tea.Cmd{}
		for _, newItem := range msg.List {
			insCmd := m.list.InsertItem(0, newItem)
			cmds = append(cmds, insCmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
