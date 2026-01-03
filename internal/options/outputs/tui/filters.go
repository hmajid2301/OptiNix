package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) applySourceFilter(items []list.Item) []list.Item {
	if m.sourceFilter.NixOS && m.sourceFilter.HomeManager && m.sourceFilter.Darwin {
		return items
	}

	filtered := []list.Item{}
	for _, item := range items {
		i := item.(Item)
		if (i.OptionFrom == "NixOS" && m.sourceFilter.NixOS) ||
			(i.OptionFrom == "Home Manager" && m.sourceFilter.HomeManager) ||
			(i.OptionFrom == "Darwin" && m.sourceFilter.Darwin) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (m Model) filterByPrefix(optionName string) tea.Cmd {
	prefix := extractPrefix(optionName)
	if prefix == "" {
		return nil
	}

	return func() tea.Msg {
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(prefix),
		}
	}
}

func extractPrefix(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) >= 2 {
		return strings.Join(parts[:2], ".")
	}
	return ""
}
