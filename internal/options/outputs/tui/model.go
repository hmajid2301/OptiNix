package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
func (i Item) FilterValue() string {
	return i.OptionName + " " + i.Desc + " " + i.OptionType + " " + i.OptionFrom
}

type SourceFilter struct {
	NixOS       bool
	HomeManager bool
	Darwin      bool
}

type ViewState struct {
	showDetail     bool
	detailViewport viewport.Model
	detailReady    bool
	statusMessage  string
}

type DisplayConfig struct {
	width  int
	height int
}

type Model struct {
	spinner      spinner.Model
	keys         *keyMap
	list         list.Model
	help         help.Model
	docStyle     lipgloss.Style
	glammy       glamour.TermRenderer
	getOptions   tea.Cmd
	totalOptions int
	sourceFilter SourceFilter
	allItems     []list.Item
	view         ViewState
	display      DisplayConfig
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
	delegate := itemDelegate{}
	l := list.New(optsList, delegate, 0, 0)
	l.Title = "NixOS Options"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	key := newKeyMap()
	help := styledHelp(help.New())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		help:       help,
		docStyle:   docStyle,
		list:       l,
		glammy:     *glammy,
		keys:       key,
		getOptions: getOptions,
		spinner:    s,
		sourceFilter: SourceFilter{
			NixOS:       true,
			HomeManager: true,
			Darwin:      true,
		},
		allItems: []list.Item{},
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.getOptions)
}
