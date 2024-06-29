package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle()

func styledHelp(help help.Model) help.Model {
	help.Styles.ShortKey = style.Foreground(lipgloss.Color("#ffffff"))
	help.Styles.ShortDesc = style.Foreground(lipgloss.Color("#ffffff"))
	help.Styles.FullKey = style.Foreground(lipgloss.Color("#ffffff"))
	help.Styles.FullDesc = style.Foreground(lipgloss.Color("#ffffff"))
	help.FullSeparator = style.Foreground(lipgloss.Color("#ffffff")).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	help.ShortSeparator = style.
		Foreground(lipgloss.Color("#ffffff")).
		PaddingLeft(1).
		PaddingRight(1).
		Render("•")
	return help
}
