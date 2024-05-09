package tui

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	// if m.done {
	// 	// err := FindOptions(m.ctx, m.db, m.userInput.Value())
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// }
	//
	return fmt.Sprintf(
		"Search for options:\n\n%s\n\n%s",
		m.userInput.View(),
		"(esc to quit)",
	) + "\n"
}
