package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	if len(m.list.Items()) == 0 {
		return m.spinner.View()
	}

	if m.showGlammy {
		selectedItem := m.list.SelectedItem().(Item)
		markdown := renderMarkdown(selectedItem)
		markdownString, _ := m.glammy.Render(markdown)
		return markdownString
	}

	// TODO: work out how to show custom help menu
	return m.list.View()
}

func renderMarkdown(item Item) string {
	template := `
# %s

## Description

%s

## Type

%s

## Default

%s

## Example

%s

## Sources

From: %s

Declared in
`
	markdown := fmt.Sprintf(
		template,
		item.OptionName,
		item.Desc,
		item.Example,
		item.OptionType,
		item.DefaultValue,
		item.OptionFrom,
	)

	// INFO: Convert a source from this path:
	// /nix/store/sdfiiqwrf78i47gzld1favdx9m5ms1cj5pb1dx0brbrbigy8ij-source/nixos/modules/programs/wayland/hyprland.nix
	// to this URL:
	// https://github.com/nixos/nixpkgs/blob/master/nixos/modules/programs/wayland/hyprland.nix
	for _, source := range item.Sources {
		url := source
		index := strings.Index(source, "nixos/modules")
		if index != -1 {
			part := source[index:]
			url = "https://github.com/nixos/nixpkgs/blob/master/" + part
		}

		sourceMarkdown := fmt.Sprintf(" - %s\n", url)
		markdown += sourceMarkdown
	}
	return markdown
}
