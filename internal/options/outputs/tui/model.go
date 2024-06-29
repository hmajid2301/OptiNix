package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type DoneMsg struct {
	List []list.Item
}

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
	keys       *keyMap
	list       list.Model
	help       help.Model
	docStyle   lipgloss.Style
	glammy     glamour.TermRenderer
	getOptions tea.Cmd
	showGlammy bool
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
	// l.SetShowHelp(false)
	key := newKeyMap()
	help := styledHelp(help.New())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{help: help,
		docStyle:   docStyle,
		list:       l,
		glammy:     *glammy,
		keys:       key,
		getOptions: getOptions,
		spinner:    s,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.getOptions)
}
