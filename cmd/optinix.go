package cmd

import (
	"context"
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"

	"github.com/spf13/cobra"
)

func NewRootCmd(ctx context.Context, db *sql.DB, logger *slog.Logger) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Version: "v0.2.0",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
	}

	baseSourcesTemplate := entities.Sources{}

	updateCmd := getUpdateCmd(ctx, db, baseSourcesTemplate, logger)
	rootCmd.AddCommand(updateCmd)
	getCmd := getGetCmd(ctx, db, baseSourcesTemplate, logger)
	rootCmd.AddCommand(getCmd)

	return rootCmd, nil
}
