package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type BrowserOpenMsg struct {
	success bool
	url     string
}

func openInBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "darwin":
			cmd = exec.Command("open", url)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		default:
			return BrowserOpenMsg{success: false, url: url}
		}

		err := cmd.Start()
		return BrowserOpenMsg{success: err == nil, url: url}
	}
}

func getDocURL(item Item) string {
	switch item.OptionFrom {
	case "NixOS":
		return fmt.Sprintf("https://search.nixos.org/options?channel=unstable&query=%s", item.OptionName)
	case "Home Manager":
		return fmt.Sprintf("https://nix-community.github.io/home-manager/options.xhtml#opt-%s",
			strings.ReplaceAll(item.OptionName, ".", "-"))
	case "Darwin":
		return "https://daiderd.com/nix-darwin/manual/index.html"
	}
	return ""
}

func getSourceURL(item Item) string {
	if len(item.Sources) > 0 {
		source := item.Sources[0]

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
	}
	return ""
}
