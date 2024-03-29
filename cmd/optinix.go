package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"github.com/spf13/cobra"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
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
	db, err := GetDB()
	if err != nil {
		return err
	}

	s, err := store.New(db)
	if err != nil {
		return err
	}

	opt := options.New(s)

	// TODO: make config
	sources := map[options.Source]string{
		options.NixOSSource:       "http://docker:8080/manual/nixos/unstable/options",
		options.HomeManagerSource: "http://docker:8080/home-manager/options.xhtml",
	}
	err = opt.SaveOptions(ctx, sources)
	if err != nil {
		return err
	}

	optionName := args[0]
	matchingOpts, err := opt.GetOptions(ctx, optionName)
	if err != nil {
		return err
	}

	// TODO: format this nicely
	for _, o := range matchingOpts {
		cmd.Println(o.Name)
		cmd.Println(o.Type)
		cmd.Println(o.Description)
		cmd.Println(o.DefaultValue)
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

func GetDB() (*sql.DB, error) {
	// TODO: what if it not set
	state := os.Getenv("XDG_STATE_HOME")
	configPath := filepath.Join(state, "optinix")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		permissions := 0755
		err = os.Mkdir(configPath, fs.FileMode(permissions))
		if err != nil {
			return nil, err
		}
	}

	dbPath := filepath.Join(configPath, "options.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	return db, err
}
