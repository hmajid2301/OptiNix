package cmd

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/nix"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func getUpdateCmd(ctx context.Context, db *sql.DB, baseSourcesTemplate entities.Sources, logger *slog.Logger) *cobra.Command {
	var noTUI bool

	updateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Update and fetch the latest options.",
		Long:    `This command will fetch the latest options from all the sources.`,
		Example: "optinix update",
		PreRun: func(cmd *cobra.Command, _ []string) {
			noTUI, _ = cmd.Flags().GetBool("no-tui")
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			sources := baseSourcesTemplate
			myStore, err := store.NewStore(db)
			if err != nil {
				return err
			}

			nixExecutor := nix.NewCmdExecutor()
			nixReader := nix.NewReader()

			var messenger entities.Messager
			var progressMsg entities.ProgressMessager

			if noTUI {
				messenger = nix.NewMessenger()
			} else {
				progressMsg = nix.NewProgressMessenger(true)
				defer progressMsg.Stop()
				messenger = progressMsg
			}

			messenger.Send("Updating flake inputs...")
			updateCmd := nixExecutor.ExecuteCommand(ctx, "nix", "flake", "update", "--flake", "./nix")
			var stderr bytes.Buffer
			updateCmd.Stderr = &stderr
			if err := updateCmd.Run(); err != nil {
				if progressMsg != nil {
					progressMsg.Stop()
				}
				return fmt.Errorf("failed to update flake: %w\n%s", err, stderr.String())
			}

			fetcher := fetch.NewFetcher(nixExecutor, nixReader, messenger, logger)
			option := options.NewSearcher(myStore, fetcher, messenger, logger)

			err = option.Save(ctx, sources)
			if err != nil {
				if progressMsg != nil {
					progressMsg.Stop()
				}
				return err
			}

			if progressMsg != nil {
				progressMsg.Finish("Update complete! Options successfully loaded to database")
			}

			return nil
		},
	}

	updateCmd.Flags().BoolVar(&noTUI, "no-tui", false, "Disable interactive progress display")
	return updateCmd
}
