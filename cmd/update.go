package cmd

import (
	"context"
	"database/sql"

	"github.com/spf13/cobra"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/nix"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func getUpdateCmd(ctx context.Context, db *sql.DB) *cobra.Command {
	updateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Update and fetch the latest options.",
		Long:    `This command will fetch the latest options from all the sources.`,
		Example: "optinix update",
		RunE: func(_ *cobra.Command, _ []string) error {
			myStore, err := store.NewStore(db)
			if err != nil {
				return err
			}

			nixExecutor := nix.NewCmdExecutor()
			nixReader := nix.NewReader()
			messenger := nix.NewMessenger()
			fetcher := fetch.NewFetcher(nixExecutor, nixReader, messenger)

			option := options.NewSearcher(myStore, fetcher, messenger)

			sources := entities.Sources{
				NixOS:       "nix/nixos-options.nix",
				HomeManager: "nix/hm-options.nix",
				Darwin:      "nix/darwin-options.nix",
			}

			err = option.Save(ctx, sources)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return updateCmd
}
