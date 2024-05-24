package options

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	flags     Flags
	userInput textinput.Model
	list      list.Model
	db        *sql.DB
	ctx       context.Context
	docStyle  lipgloss.Style
}

// TODO: rename this
type Flags struct {
	OptionName   string
	Limit        int64
	ForceRefresh bool
}

func NewTUI(ctx context.Context, db *sql.DB, flags Flags) Model {
	ti := textinput.New()
	ti.SetValue(flags.OptionName)

	//nolint: mnd
	docStyle := lipgloss.NewStyle().Margin(1, 2)
	opts, err := FindOptions(ctx, db, flags)
	if err != nil {
		log.Fatal(err)
	}

	optsList := []list.Item{}
	for _, opt := range opts {
		newDescription := strings.ReplaceAll(opt.Description, ".", ".\n")
		listItem := item{title: opt.Name, desc: newDescription}
		optsList = append(optsList, listItem)
	}
	l := list.New(optsList, list.NewDefaultDelegate(), 0, 0)

	return Model{docStyle: docStyle, userInput: ti, ctx: ctx, db: db, flags: flags, list: l}
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

// TODO: move to a proper place
func FindOptions(ctx context.Context,
	db *sql.DB,
	flags Flags,
) (opts []OptionWithSources, err error) {
	s, err := NewStore(db)
	if err != nil {
		return nil, err
	}

	// TODO: should this be setup with constructors
	cmdExecutor := NixCmdExecutor{}
	nixReader := NixReader{}
	fetcher := NewFetcher(cmdExecutor, nixReader)

	opt := NewOptions(s, fetcher)

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

	sources := Sources{
		NixOS:       string(nixosPath),
		HomeManager: string(homeManagerPath),
		Darwin:      string(darwinPath),
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
		// TODO: why does this needed on NixOS but not on Ubuntu
		"NIX_PATH=/etc/nix/inputs",
		"--no-out-link",
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	trimmedOuput := strings.TrimSpace(string(output))
	return trimmedOuput, nil
}
