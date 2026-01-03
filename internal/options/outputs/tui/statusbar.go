package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	statusItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62"))
)

func renderStatusBar(total, filtered, currentIndex int, width int, sourceFilter SourceFilter, filterQuery string) string {
	left := fmt.Sprintf("Total: %d", total)

	var middle string
	if filtered != total {
		middle = fmt.Sprintf("Showing: %d (%d/%d)", filtered, currentIndex+1, filtered)
	} else if filtered > 0 {
		middle = fmt.Sprintf("%d/%d", currentIndex+1, filtered)
	}

	right := "? for help"

	leftStr := statusItemStyle.Render(left)
	rightStr := statusItemStyle.Render(right)
	middleStr := statusItemStyle.Render(middle)

	leftWidth := lipgloss.Width(leftStr)
	rightWidth := lipgloss.Width(rightStr)
	middleWidth := lipgloss.Width(middleStr)

	padding := width - leftWidth - rightWidth - middleWidth
	if padding < 0 {
		padding = 0
	}

	var statusBar string
	if middle != "" {
		leftPad := padding / 2
		rightPad := padding - leftPad
		statusBar = statusBarStyle.Render(
			leftStr +
				strings.Repeat(" ", leftPad) +
				middleStr +
				strings.Repeat(" ", rightPad) +
				rightStr,
		)
	} else {
		statusBar = statusBarStyle.Render(
			leftStr +
				strings.Repeat(" ", padding) +
				rightStr,
		)
	}

	filterBar := renderFilterBar(sourceFilter, filterQuery, width)
	if filterBar != "" {
		return filterBar + "\n" + statusBar
	}

	return statusBar
}

func renderFilterBar(sourceFilter SourceFilter, filterQuery string, width int) string {
	if sourceFilter.NixOS && sourceFilter.HomeManager && sourceFilter.Darwin && filterQuery == "" {
		return ""
	}

	activeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Background(lipgloss.Color("235")).
		Padding(0, 1).
		Bold(true)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	var parts []string

	nixosBtn := "NixOS"
	if sourceFilter.NixOS {
		parts = append(parts, activeStyle.Render(nixosBtn))
	} else {
		parts = append(parts, inactiveStyle.Render(nixosBtn))
	}

	hmBtn := "Home Manager"
	if sourceFilter.HomeManager {
		parts = append(parts, activeStyle.Render(hmBtn))
	} else {
		parts = append(parts, inactiveStyle.Render(hmBtn))
	}

	darwinBtn := "Darwin"
	if sourceFilter.Darwin {
		parts = append(parts, activeStyle.Render(darwinBtn))
	} else {
		parts = append(parts, inactiveStyle.Render(darwinBtn))
	}

	filterBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	content := strings.Join(parts, " ")

	if filterQuery != "" {
		searchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Italic(true)
		content += "  " + searchStyle.Render(fmt.Sprintf("Filter: %s", filterQuery))
	}

	return filterBarStyle.Width(width).Render(content)
}
