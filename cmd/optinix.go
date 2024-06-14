package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/config"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/nix"
	"gitlab.com/hmajid2301/optinix/internal/options/outputs/plaintext"
	"gitlab.com/hmajid2301/optinix/internal/options/outputs/tui"
	"gitlab.com/hmajid2301/optinix/internal/options/store"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type ArgsAndFlags struct {
	OptionName   string
	Limit        int64
	ForceRefresh bool
	NoTUI        bool
}

func Execute(ctx context.Context, ddl string) error {
	var forceRefresh bool
	var noTUI bool
	var limit int64

	db, err := getDB(ctx, ddl)
	defer func() {
		err = db.Close()
	}()

	if err != nil {
		return err
	}

	rootCmd := &cobra.Command{
		Version: "v0.1.0",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
		Example: "optinix hyprland",
		PreRun: func(cmd *cobra.Command, _ []string) {
			forceRefresh, _ = cmd.Flags().GetBool("force-refresh")
			noTUI, _ = cmd.Flags().GetBool("no-tui")
			limit, _ = cmd.Flags().GetInt64("limit")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			flags := ArgsAndFlags{
				ForceRefresh: forceRefresh,
				Limit:        limit,
				NoTUI:        noTUI,
			}

			if len(args) > 0 {
				flags.OptionName = args[0]
			}

			outputOptions(ctx, flags, db)
			return nil
		},
	}

	rootCmd.Flags().BoolVar(&forceRefresh, "force-refresh", false, "If set will force a refresh of the options")
	rootCmd.Flags().BoolVar(&noTUI, "no-tui", false, "If set will not show the TUI and just print the options to stdout")
	//nolint: mnd
	rootCmd.Flags().Int64Var(&limit, "limit", 10, "Limit the number of results returned")
	return rootCmd.ExecuteContext(ctx)
}

func getDB(ctx context.Context, ddl string) (*sql.DB, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := store.GetDB(conf.DBFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, fmt.Errorf("failed to create database schema: %w", err)
	}
	return db, nil
}

func GetOptions(ctx context.Context, db *sql.DB, flags ArgsAndFlags) tea.Cmd {
	return func() tea.Msg {
		options, err := FindOptions(ctx, db, flags)
		if err != nil {
			tea.Printf("Failed to get options %s\n", err)
		}

		optsList := []list.Item{}
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
			optsList = append(optsList, listItem)
		}

		return tui.DoneMsg{
			List: optsList,
		}
	}
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
	messenger := nix.NewMessenger()
	fetcher := fetch.NewFetcher(nixExecutor, nixReader, messenger)

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

	// INFO: If not option name passed then get all options.
	if flags.OptionName == "" {
		opts, err = opt.GetAllOptions(ctx)
	} else {
		opts, err = opt.FindOptions(ctx, flags.OptionName, flags.Limit)
	}

	if err != nil {
		return nil, err
	}
	return opts, nil
}

func outputOptions(ctx context.Context, flags ArgsAndFlags, db *sql.DB) {
	if flags.NoTUI {
		options, err := FindOptions(ctx, db, flags)
		if err != nil {
			fmt.Printf("Failed to get options %s\n", err)
		}
		plaintext.Output(options)
	} else {
		getOptionsFunc := GetOptions(ctx, db, flags)
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
