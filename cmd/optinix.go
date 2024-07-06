package cmd

import (
	"context"
	"database/sql"
	"fmt"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options/config"
	"gitlab.com/hmajid2301/optinix/internal/options/store"

	"github.com/spf13/cobra"
)

func NewRootCmd(ctx context.Context, ddl string) (*cobra.Command, error) {
	db, err := getDB(ctx, ddl)
	defer func() {
		err = db.Close()
	}()

	if err != nil {
		return nil, err
	}

	rootCmd := &cobra.Command{
		Version: "v0.1.0",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
	}

	updateCmd := getUpdateCmd(ctx, db)
	rootCmd.AddCommand(updateCmd)
	getCmd := getGetCmd(ctx, db)
	rootCmd.AddCommand(getCmd)

	return rootCmd, nil
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
