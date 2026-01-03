package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	typeBadgeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			MarginRight(1)

	nixosBadgeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("33")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			MarginRight(1)

	hmBadgeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("99")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			MarginRight(1)

	darwinBadgeStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("208")).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1).
				MarginRight(1)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 3 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	var sourceBadge string
	switch i.OptionFrom {
	case "NixOS":
		sourceBadge = nixosBadgeStyle.Render("NixOS")
	case "Home Manager":
		sourceBadge = hmBadgeStyle.Render("HM")
	case "Darwin":
		sourceBadge = darwinBadgeStyle.Render("Darwin")
	}

	typeBadge := typeBadgeStyle.Render(i.OptionType)

	title := i.OptionName
	desc := i.Desc
	if len(desc) > 80 {
		desc = desc[:77] + "..."
	}

	if index == m.Index() {
		title = selectedItemStyle.Render("▸ " + title)
		badges := lipgloss.JoinHorizontal(lipgloss.Left, sourceBadge, typeBadge)
		title = lipgloss.JoinHorizontal(lipgloss.Left, title, " ", badges)
		desc = descStyle.Render(desc)
	} else {
		title = itemStyle.Render(title)
		desc = descStyle.Render(desc)
	}

	fmt.Fprint(w, title+"\n"+desc)
}
