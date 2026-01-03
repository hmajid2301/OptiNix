package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var boxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Padding(1, 2).
	MarginTop(1)

func (m Model) View() string {
	if len(m.list.Items()) == 0 {
		loadingMsg := "Loading options..."
		spinnerView := fmt.Sprintf("\n\n   %s %s\n\n", m.spinner.View(), loadingMsg)
		return spinnerView
	}

	if m.view.showDetail {
		if !m.view.detailReady {
			return "Loading detail view..."
		}

		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			MarginTop(1)

		view := m.view.detailViewport.View()

		if m.view.statusMessage != "" {
			view += "\n" + statusStyle.Render(m.view.statusMessage)
		}

		scrollInfo := ""
		if m.view.detailViewport.TotalLineCount() > m.view.detailViewport.Height {
			scrollInfo = fmt.Sprintf("\n↑/↓ to scroll (%d%%)", int(m.view.detailViewport.ScrollPercent()*100))
		}

		return boxStyle.Render(view + scrollInfo)
	}

	listView := m.list.View()
	filterQuery := ""
	if m.list.FilterState() == list.Filtering || m.list.FilterState() == list.FilterApplied {
		filterQuery = m.list.FilterValue()
	}
	currentIndex := m.list.Index()
	statusBar := renderStatusBar(m.totalOptions, len(m.list.VisibleItems()), currentIndex, m.display.width, m.sourceFilter, filterQuery)

	return lipgloss.JoinVertical(lipgloss.Left, listView, statusBar)
}
