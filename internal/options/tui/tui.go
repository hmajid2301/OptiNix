package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	list     list.Model
	docStyle lipgloss.Style
}

func NewTUI(options []entities.Option) Model {
	//nolint: mnd
	docStyle := lipgloss.NewStyle().Margin(1, 2)

	optsList := []list.Item{}
	for _, opt := range options {
		newDescription := strings.ReplaceAll(opt.Description, ".", ".\n")
		listItem := item{title: opt.Name, desc: newDescription}
		optsList = append(optsList, listItem)
	}
	l := list.New(optsList, list.NewDefaultDelegate(), 0, 0)

	return Model{docStyle: docStyle, list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
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
	return m.docStyle.Render(m.list.View())
}
