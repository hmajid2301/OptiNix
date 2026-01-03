package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/nix"
	"gitlab.com/hmajid2301/optinix/internal/options/outputs/plaintext"
	"gitlab.com/hmajid2301/optinix/internal/options/outputs/tui"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type GetArgsAndFlags struct {
	OptionName string
	NoTUI      bool
}

func getGetCmd(ctx context.Context, db *sql.DB, baseSourcesTemplate entities.Sources, logger *slog.Logger) *cobra.Command {
	var noTUI bool

	getCmd := &cobra.Command{
		Use:     "get",
		Short:   "Finds options",
		Long:    `This command will find options based on the name. If no name is provided it will return all options.`,
		Example: "optinix get niri",
		PreRun: func(cmd *cobra.Command, _ []string) {
			noTUI, _ = cmd.Flags().GetBool("no-tui")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sources := baseSourcesTemplate
			flags := GetArgsAndFlags{
				NoTUI: noTUI,
			}

			if len(args) > 0 {
				flags.OptionName = args[0]
			}

			if flags.OptionName == "" && !flags.NoTUI {
				return errors.New("option name is required when using the TUI, pass --no-tui to disable the TUI")
			}

			outputOptions(ctx, flags, db, sources, logger)
			return nil
		},
	}

	getCmd.Flags().BoolVar(&noTUI, "no-tui", false, "If set will not show the TUI and just print the options to stdout")
	return getCmd
}

func outputOptions(ctx context.Context, flags GetArgsAndFlags, db *sql.DB, sources entities.Sources, logger *slog.Logger) {
	if flags.NoTUI {
		options, err := findOptions(ctx, db, flags, sources, logger)
		if err != nil {
			fmt.Printf("Failed to get options %s\n", err)
		}
		plaintext.Output(options)
	} else {
		getOptionsFunc := getOptions(ctx, db, flags, sources, logger)
		myTui, err := tui.NewTUI(getOptionsFunc)
		if err != nil {
			fmt.Println(err)
		}

		p := tea.NewProgram(myTui)
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
		}
	}
}

func getOptions(ctx context.Context, db *sql.DB, flags GetArgsAndFlags, sources entities.Sources, logger *slog.Logger) tea.Cmd {
	return func() tea.Msg {
		options, err := findOptions(ctx, db, flags, sources, logger)
		if err != nil {
			tea.Printf("Failed to get options %s\n", err)
		}

		opts := []list.Item{}
		for _, opt := range options {
			newDescription := strings.ReplaceAll(opt.Description, ".", ".\n")
			listItem := tui.Item{
				OptionName:   opt.Name,
				Desc:         newDescription,
				DefaultValue: opt.Default,
				Example:      opt.Example,
				OptionType:   opt.Type,
				OptionFrom:   opt.OptionFrom,
				Sources:      opt.Sources,
			}
			opts = append(opts, listItem)
		}

		return tui.DoneMsg{
			List: opts,
		}
	}
}

func findOptions(ctx context.Context,
	db *sql.DB,
	flags GetArgsAndFlags,
	sources entities.Sources,
	logger *slog.Logger,
) (opts []entities.Option, err error) {
	myStore, err := store.NewStore(db)
	if err != nil {
		return nil, err
	}

	nixExecutor := nix.NewCmdExecutor()
	nixReader := nix.NewReader()
	messenger := nix.NewMessenger()
	fetcher := fetch.NewFetcher(nixExecutor, nixReader, messenger, logger)

	option := options.NewSearcher(myStore, fetcher, messenger, logger)

	optionCount, err := myStore.CountOptions(ctx)
	if err != nil {
		return nil, err
	}

	if optionCount == 0 {
		err = option.Save(ctx, sources)
		if err != nil {
			return nil, err
		}
	}

	if flags.OptionName == "" {
		opts, err = option.GetAll(ctx)
	} else {
		opts, err = option.Find(ctx, flags.OptionName)
	}

	if err != nil {
		return nil, err
	}
	return opts, nil
}
