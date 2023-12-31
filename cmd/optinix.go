package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"gitlab.com/majiy00/go/clis/optinix/internal/options"
	"gitlab.com/majiy00/go/clis/optinix/internal/options/store"
)

var rootCmd = &cobra.Command{
	Short: "optinix - a CLI tool to search nix options",
	Long: `OptiNix is tool you can use on the command line to search options for both NixOS and home-manager
	rather than needing to go to a website i.e. nixos.org or mynixos.com.`,
	Args: cobra.ExactArgs(1),
	RunE: FindOptions,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func FindOptions(cmd *cobra.Command, args []string) error {
	ctx := gracefulShutdown()
	db, err := getDB()
	if err != nil {
		return err
	}

	s, err := store.New(db)
	if err != nil {
		return err
	}

	opt := options.New(s)

	var count int64
	db.Count(&count)

	if count == 0 {
		err = opt.SaveOptions(ctx)
		if err != nil {
			return err
		}
	}

	optionName := args[0]
	matchingOpts, err := opt.GetOptions(ctx, optionName)
	if err != nil {
		return err
	}

	for _, o := range matchingOpts {
		cmd.Println(o.Name)
		cmd.Println(o.Type)
		cmd.Println(o.Description)
		cmd.Println(o.Default)
		cmd.Println(o.Example)

		for _, s := range o.Sources {
			cmd.Println(s)
		}
	}

	return nil
}

func gracefulShutdown() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	return ctx
}

func getDB() (*gorm.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "optinix")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		permissions := 0755
		err = os.Mkdir(configPath, fs.FileMode(permissions))
		if err != nil {
			return nil, err
		}
	}

	dbPath := filepath.Join(configPath, "options.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	return db, err
}
