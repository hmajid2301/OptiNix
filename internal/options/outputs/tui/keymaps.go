package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// default keybind definitions
type keyMap struct {
	filter     key.Binding
	quit       key.Binding
	more       key.Binding
	choose     key.Binding
	toggle     key.Binding
	selectDown key.Binding
	selectUp   key.Binding
	up         key.Binding
	down       key.Binding
	// nextPage   key.Binding
	// prevPage   key.Binding
	home key.Binding
	end  key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q/esc", "quit"),
		),
		more: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "copy"),
		),
		selectDown: key.NewBinding(
			key.WithKeys("ctrl+down", "ctrl+j"),
			key.WithHelp("ctrl+↓/↑", "select"),
		),
		selectUp: key.NewBinding(
			key.WithKeys("ctrl+up", "ctrl+k"),
			key.WithHelp("ctrl+↓/↑", "select"),
		),
		up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		home: key.NewBinding(
			key.WithKeys("home", "g"),
		),
		end: key.NewBinding(
			key.WithKeys("end", "G"),
		),
		toggle: key.NewBinding(
			key.WithKeys("t"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.choose, k.filter, k.toggle, k.up, k.down, k.more,
	}
}

// not currently in use as intentionally being overridden by the default
// full help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.home, k.end},
		{k.choose, k.toggle},
		{k.selectDown, k.selectUp},
		{k.filter, k.quit},
	}
}
