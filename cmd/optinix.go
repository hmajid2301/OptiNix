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
	"gitlab.com/hmajid2301/optinix/internal/options/store"
	"gitlab.com/hmajid2301/optinix/internal/options/tui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func Execute(ctx context.Context, ddl string) error {
	var forceRefresh bool
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
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, _ []string) {
			forceRefresh, _ = cmd.Flags().GetBool("force-refresh")
			limit, _ = cmd.Flags().GetInt64("limit")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			flags := ArgsAndFlags{
				OptionName:   args[0],
				ForceRefresh: forceRefresh,
				Limit:        limit,
			}

			getOptionsFunc := GetOptions(ctx, db, flags)
			myTui, err := tui.NewTUI(getOptionsFunc)
			if err != nil {
				return err
			}

			p := tea.NewProgram(myTui)
			if _, err := p.Run(); err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}

	rootCmd.Flags().BoolVar(&forceRefresh, "force-refresh", false, "If set will force a refresh of the options")

	//nolint: mnd
	rootCmd.Flags().Int64Var(&limit, "limit", 10, "Limit the number of results returned")
	return rootCmd.ExecuteContext(ctx)
}

func GetOptions(ctx context.Context, db *sql.DB, flag ArgsAndFlags) tea.Cmd {
	return func() tea.Msg {
		options, err := FindOptions(ctx, db, flag)
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

type ArgsAndFlags struct {
	OptionName   string
	Limit        int64
	ForceRefresh bool
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

	opts, err = opt.GetOptions(ctx, flags.OptionName, flags.Limit)
	if err != nil {
		return nil, err
	}

	return opts, nil
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
