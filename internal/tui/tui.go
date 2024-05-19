package tui

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type Model struct {
	userInput textinput.Model
	// selected  int
	done  bool
	ctx   context.Context
	db    *sql.DB
	flags Flags
}

type Flags struct {
	ForceRefresh bool
	Limit        int64
}

func New(ctx context.Context, db *sql.DB, flags Flags) Model {
	ti := textinput.New()
	ti.CharLimit = 100
	ti.Placeholder = "Type the option to search for"
	ti.Width = 30
	ti.Focus()

	return Model{userInput: ti, ctx: ctx, db: db, flags: flags}
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
		opts, err := FindOptions(m.ctx, m.db, m.userInput.Value(), m.flags)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(opts)
	}

	return fmt.Sprintf(
		"Search for options:\n\n%s\n\n%s",
		m.userInput.View(),
		"(esc to quit)",
	) + "\n"
}

func FindOptions(ctx context.Context,
	db *sql.DB,
	optionName string,
	flags Flags,
) (opts []store.OptionWithSources, err error) {
	s, err := store.NewStore(db)
	if err != nil {
		return nil, err
	}

	// TODO: should this be setup with constructors
	cmdExecutor := NixCmdExecutor{}
	nixReader := NixReader{}
	fetcher := options.NewFetcher(cmdExecutor, nixReader)

	opt := options.NewOptions(s, fetcher)

	// TODO: should I read file and evalute expression?
	nixosPath, err := nixReader.Read("nix/nixos-options.nix")
	if err != nil {
		return nil, err
	}

	homeManagerPath, err := nixReader.Read("nix/hm-options.nix")
	if err != nil {
		return nil, err
	}

	darwinPath, err := nixReader.Read("nix/darwin-options.nix")
	if err != nil {
		return nil, err
	}

	sources := options.Sources{
		NixOS:       string(nixosPath),
		HomeManager: string(homeManagerPath),
		Darwin:      string(darwinPath),
	}
	err = opt.SaveOptions(ctx, sources, flags.ForceRefresh)
	if err != nil {
		return nil, err
	}

	opts, err = opt.GetOptions(ctx, optionName, flags.Limit)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

type NixReader struct{}

func (f NixReader) Read(pathToExpression string) ([]byte, error) {
	nixExpression, err := os.ReadFile(pathToExpression)
	return nixExpression, err
}

type NixCmdExecutor struct{}

func (n NixCmdExecutor) Executor(expression string) (string, error) {
	cmd := exec.Command("nix-build", "-E", expression)
	cmd.Env = append(cmd.Env,
		"NIXPKGS_ALLOW_UNFREE=1",
		"NIXPKGS_ALLOW_BROKEN=1",
		"NIXPKGS_ALLOW_INSECURE=1",
		"NIXPKGS_ALLOW_UNSUPPORTED_SYSTEM=1",
		"NIX_PATH=/etc/nix/inputs",
		"--no-out-link",
	)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	trimmedOuput := strings.TrimSpace(string(output))
	return trimmedOuput, nil
}
