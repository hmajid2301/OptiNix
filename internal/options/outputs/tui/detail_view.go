package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginTop(1).
			MarginBottom(1)

	codeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("235")).
			Padding(1, 2)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Underline(true)

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

func renderDetailedView(item Item) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(item.OptionName))
	b.WriteString("\n\n")

	metaInfo := fmt.Sprintf("Type: %s  •  From: %s",
		item.OptionType,
		item.OptionFrom,
	)
	b.WriteString(metaStyle.Render(metaInfo))
	b.WriteString("\n")

	if item.Desc != "" {
		b.WriteString(sectionStyle.Render("Description"))
		b.WriteString("\n")
		b.WriteString(item.Desc)
		b.WriteString("\n")
	}

	if item.DefaultValue != "" && item.DefaultValue != "null" {
		b.WriteString(sectionStyle.Render("Default Value"))
		b.WriteString("\n")
		b.WriteString(codeStyle.Render(item.DefaultValue))
		b.WriteString("\n")
	}

	if item.Example != "" && item.Example != "null" {
		b.WriteString(sectionStyle.Render("Example"))
		b.WriteString("\n")
		b.WriteString(codeStyle.Render(item.Example))
		b.WriteString("\n")
	}

	if len(item.Sources) > 0 {
		b.WriteString(sectionStyle.Render("Sources"))
		b.WriteString("\n")
		for _, source := range item.Sources {
			url := convertToGitHubURL(source)
			b.WriteString("  • ")
			b.WriteString(urlStyle.Render(url))
			b.WriteString("\n")
		}
	}

	// Add official documentation links
	b.WriteString(sectionStyle.Render("Documentation"))
	b.WriteString("\n")

	switch item.OptionFrom {
	case "NixOS":
		searchURL := fmt.Sprintf("https://search.nixos.org/options?channel=unstable&query=%s", item.OptionName)
		b.WriteString("  • ")
		b.WriteString(urlStyle.Render(searchURL))
		b.WriteString("\n")
	case "Home Manager":
		hmURL := fmt.Sprintf("https://nix-community.github.io/home-manager/options.xhtml#opt-%s",
			strings.ReplaceAll(item.OptionName, ".", "-"))
		b.WriteString("  • ")
		b.WriteString(urlStyle.Render(hmURL))
		b.WriteString("\n")
	case "Darwin":
		// Darwin options are typically documented in the nix-darwin manual
		darwinURL := "https://daiderd.com/nix-darwin/manual/index.html"
		b.WriteString("  • ")
		b.WriteString(urlStyle.Render(darwinURL))
		b.WriteString("\n")
	}

	helpText := metaStyle.Render("\nPress 't' to return • 'y' copy • 'r' related • 'o' docs • 'O' source • 'q' quit")
	b.WriteString("\n")
	b.WriteString(helpText)

	return b.String()
}


func convertToGitHubURL(source string) string {
	if strings.Contains(source, "github.com") {
		return source
	}

	index := strings.Index(source, "nixos/modules")
	if index != -1 {
		part := source[index:]
		return "https://github.com/nixos/nixpkgs/blob/master/" + part
	}

	index = strings.Index(source, "modules/")
	if index != -1 {
		part := source[index:]
		if strings.Contains(source, "home-manager") {
			return "https://github.com/nix-community/home-manager/blob/master/" + part
		}
		if strings.Contains(source, "darwin") || strings.Contains(source, "nix-darwin") {
			return "https://github.com/LnL7/nix-darwin/blob/master/" + part
		}
	}

	return source
}
