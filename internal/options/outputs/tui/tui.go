package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/glamour"
)

type Item struct {
	OptionName   string
	OptionType   string
	OptionFrom   string
	Desc         string
	Example      string
	DefaultValue string
	Sources      []string
}

func (i Item) Title() string       { return i.OptionName }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.OptionName }

type Model struct {
	spinner    spinner.Model
	keys       *listKeyMap
	list       list.Model
	docStyle   lipgloss.Style
	glammy     glamour.TermRenderer
	getOptions tea.Cmd
	showGlammy bool
}

type DoneMsg struct {
	List []list.Item
}

func NewTUI(getOptions tea.Cmd) (Model, error) {
	docStyle := lipgloss.NewStyle()

	wordWrap := 80
	glammy, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wordWrap))

	if err != nil {
		return Model{}, err
	}

	optsList := []list.Item{}
	l := list.New(optsList, list.NewDefaultDelegate(), 0, 0)
	listKeys := newListKeyMap()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{docStyle: docStyle, list: l, glammy: *glammy, keys: listKeys, getOptions: getOptions, spinner: s}, nil
}

type listKeyMap struct {
	openModal key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		openModal: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.getOptions)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit

		case key.Matches(msg, m.keys.openModal):
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
