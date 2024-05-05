package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"gitlab.com/hmajid2301/optinix/cmd"
	"gitlab.com/hmajid2301/optinix/internal/options/config"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

//go:embed db/schema.sql
var ddl string

func main() {
	var exitCode int
	ctx := gracefulShutdown()

	// TODO: proper error messages return back to CLI
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	// TODO: move to a better lib
	db, err := store.GetDB(conf.DBFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops.  '%s'", err)
		return
	}

	defer func() {
		err = db.Close()
	}()

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	defer func() { os.Exit(exitCode) }()
	if err := cmd.Execute(ctx, db); err != nil {
		exitCode = 1
		return
	}
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
