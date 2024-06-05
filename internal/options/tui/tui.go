package tui

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/glamour"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/nix"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type item struct {
	title        string
	description  string
	defaultValue string
	optionType   string
	optionFrom   string
	sources      []string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type Model struct {
	spinner    spinner.Model
	keys       *listKeyMap
	list       list.Model
	docStyle   lipgloss.Style
	glammy     glamour.TermRenderer
	showGlammy bool
	ctx        context.Context
	db         *sql.DB
	flag       ArgsAndFlags
}

type doneMsg struct {
	list []list.Item
}

func NewTUI(ctx context.Context, db *sql.DB, flag ArgsAndFlags) *Model {
	//nolint: mnd
	docStyle := lipgloss.NewStyle().Margin(1, 2)

	// TODO: Handle errors
	// TODO: Terminal Width set to 120
	glammy, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		//nolint: mnd
		glamour.WithWordWrap(120))

	optsList := []list.Item{}
	l := list.New(optsList, list.NewDefaultDelegate(), 0, 0)
	listKeys := newListKeyMap()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{docStyle: docStyle, list: l, glammy: *glammy, keys: listKeys, ctx: ctx, db: db, flag: flag, spinner: s}
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

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, startLongRunningProcess(m))
}

func startLongRunningProcess(m *Model) tea.Cmd {
	return func() tea.Msg {
		options, err := FindOptions(m.ctx, m.db, m.flag)
		if err != nil {
			tea.Printf("Failed to get options %s\n", err)
		}

		optsList := []list.Item{}
		for _, opt := range options {
			newDescription := strings.ReplaceAll(opt.Description, ".", ".\n")
			listItem := item{
				title:        opt.Name,
				description:  newDescription,
				defaultValue: opt.Default,
				optionType:   opt.Type,
				optionFrom:   opt.OptionFrom,
				sources:      opt.Sources,
			}
			optsList = append(optsList, listItem)
		}

		return doneMsg{
			list: optsList,
		}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case doneMsg:
		cmds := []tea.Cmd{}
		for _, newItem := range msg.list {
			insCmd := m.list.InsertItem(0, newItem)
			cmds = append(cmds, insCmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if len(m.list.Items()) == 0 {
		return m.spinner.View()
	}

	if m.showGlammy {
		selectedItem := m.list.SelectedItem().(item)
		markdown := renderMarkdown(selectedItem)
		markdownString, _ := m.glammy.Render(markdown)
		return markdownString
	}

	return m.list.View()
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

From: %s

Declared in
`
	markdown := fmt.Sprintf(template, item.title, item.description, item.optionType, item.defaultValue, item.optionFrom)

	// INFO: Convert a source from this path:
	// /nix/store/sdfiiqwrf78i47gzld1favdx9m5ms1cj5pb1dx0brbrbigy8ij-source/nixos/modules/programs/wayland/hyprland.nix
	// to this URL:
	// https://github.com/nixos/nixpkgs/blob/master/nixos/modules/programs/wayland/hyprland.nix
	for _, source := range item.sources {
		index := strings.Index(source, "nixos/modules")
		if index == -1 {
			continue
		}

		part := source[index:]
		url := "https://github.com/nixos/nixpkgs/blob/master/" + part
		sourceMarkdown := fmt.Sprintf(" - %s\n", url)
		markdown += sourceMarkdown
	}
	return markdown
}

type ArgsAndFlags struct {
	OptionName   string
	Limit        int64
	ForceRefresh bool
}

// TODO: better name
type Updater struct{}

func (u Updater) SendMessage(msg string) {
	tea.Println(msg)
}

func FindOptions(ctx context.Context,
	db *sql.DB,
	flags ArgsAndFlags,
) (opts []entities.Option, err error) {
	myStore, err := store.NewStore(db)
	if err != nil {
		return nil, err
	}

	nixExecutor := nix.NewCmdExecutor()
	nixReader := nix.NewReader()
	updater := Updater{}
	fetcher := fetch.NewFetcher(nixExecutor, nixReader, updater)

	opt := options.NewSearcher(myStore, fetcher)

	sources := entities.Sources{
		NixOS:       "nix/nixos-options.nix",
		HomeManager: "nix/hm-options.nix",
		Darwin:      "nix/darwin-options.nix",
	}
	err = opt.SaveOptions(ctx, sources, flags.ForceRefresh)
	if err != nil {
		return nil, err
	}

	opts, err = opt.GetOptions(ctx, flags.OptionName, flags.Limit)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
