package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/glamour"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

type item struct {
	title        string
	description  string
	defaultValue string
	optionType   string
	sources      []string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type Model struct {
	keys       *listKeyMap
	list       list.Model
	docStyle   lipgloss.Style
	glammy     glamour.TermRenderer
	showGlammy bool
}

func NewTUI(options []entities.Option) Model {
	//nolint: mnd
	docStyle := lipgloss.NewStyle().Margin(1, 2)

	optsList := []list.Item{}
	for _, opt := range options {
		newDescription := strings.ReplaceAll(opt.Description, ".", ".\n")
		listItem := item{
			title:        opt.Name,
			description:  newDescription,
			defaultValue: opt.Default,
			optionType:   opt.Type,
			sources:      opt.Sources,
		}
		optsList = append(optsList, listItem)
	}
	l := list.New(optsList, list.NewDefaultDelegate(), 0, 0)

	// TODO: Handle errors
	// TODO: Terminal Width set to 120
	glammy, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		//nolint: mnd
		glamour.WithWordWrap(120))

	listKeys := newListKeyMap()
	return Model{docStyle: docStyle, list: l, glammy: *glammy, keys: listKeys}
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
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit

		case key.Matches(msg, m.keys.openModal):
			cmd := m.list.NewStatusMessage(fmt.Sprint(m.list.Index()))
			m.showGlammy = !m.showGlammy
			return m, tea.Batch(cmd)
		}
	case tea.WindowSizeMsg:
		h, v := m.docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.showGlammy {
		selectedItem := m.list.SelectedItem().(item)
		markdown := renderMarkdown(selectedItem)
		markdownString, _ := m.glammy.Render(markdown)
		return markdownString
	}
	return m.docStyle.Render(m.list.View())
}

func renderMarkdown(item item) string {
	template := `
# %s

## Description

%s

## Type

%s

## Default

%s

## Sources

Declared in:
`
	markdown := fmt.Sprintf(template, item.title, item.description, item.optionType, item.defaultValue)

	for _, source := range item.sources {
		tea.Printf(source)
		sourceMarkdown := fmt.Sprintf(" - %s\n", source)
		markdown += sourceMarkdown
	}
	return markdown
}
