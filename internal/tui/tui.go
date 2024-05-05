package tui

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/config"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type Model struct {
	userInput textinput.Model
	// selected  int
	done bool
	ctx  context.Context
	db   *sql.DB
}

func New(ctx context.Context, db *sql.DB) Model {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Placeholder = "Type the option to search for"
	ti.Width = 30
	ti.Focus()

	return Model{userInput: ti, ctx: ctx, db: db}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		}

	default:
		return m, nil
	}
	m.userInput, cmd = m.userInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.done {
		err := FindOptions(m.ctx, m.db, m.userInput.Value())
		if err != nil {
			log.Fatal(err)
		}
	}

	return fmt.Sprintf(
		"Search for options:\n\n%s\n\n%s",
		m.userInput.View(),
		"(esc to quit)",
	) + "\n"
}

func FindOptions(ctx context.Context, db *sql.DB, optionName string) (err error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	s, err := store.New(db)
	if err != nil {
		return err
	}

	opt := options.New(s)

	sources := map[options.Source]string{
		options.NixOSSource:       conf.Sources.NixOSURL,
		options.HomeManagerSource: conf.Sources.HomeManagerURL,
	}
	err = opt.SaveOptions(ctx, sources)
	if err != nil {
		return err
	}

	matchingOpts, err := opt.GetOptions(ctx, optionName)
	if err != nil {
		return err
	}

	// TODO: format this nicely
	for _, o := range matchingOpts {
		fmt.Println(o.Name)
		fmt.Println(o.Type)
		fmt.Println(o.Description)
		fmt.Println(o.DefaultValue)
		fmt.Println(o.Example)

		for _, s := range o.Sources {
			fmt.Println(s)
		}
	}

	return nil
}
