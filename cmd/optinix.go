package cmd

import (
	"context"
	"database/sql"
	"embed"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"

	"github.com/spf13/cobra"
)

func NewRootCmd(ctx context.Context, db *sql.DB, nixExpressions embed.FS) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Version: "v0.1.2",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
	}

	no, err := nixExpressions.ReadFile("nix/nixos-options.nix")
	if err != nil {
		return nil, err
	}

	ho, err := nixExpressions.ReadFile("nix/hm-options.nix")
	if err != nil {
		return nil, err
	}

	do, err := nixExpressions.ReadFile("nix/darwin-options.nix")
	if err != nil {
		return nil, err
	}

	sources := entities.Sources{
		NixOS:       string(no),
		HomeManager: string(ho),
		Darwin:      string(do),
	}

	updateCmd := getUpdateCmd(ctx, db, sources)
	rootCmd.AddCommand(updateCmd)
	getCmd := getGetCmd(ctx, db, sources)
	rootCmd.AddCommand(getCmd)

	return rootCmd, nil
}
