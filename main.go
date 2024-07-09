package main

import (
	"context"
	"database/sql"
	"embed"
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

//go:embed nix
var embeddedFiles embed.FS

func main() {
	var exitCode int
	ctx := gracefulShutdown()
	defer func() { os.Exit(exitCode) }()

	db, err := getDB(ctx, ddl)
	if err != nil {
		fmt.Println("Error creating db command")
		fmt.Print(err)
		exitCode = 1
		return
	}

	rootCmd, err := cmd.NewRootCmd(ctx, db, embeddedFiles)
	if err != nil {
		fmt.Println("Error creating root command")
		fmt.Print(err)
		exitCode = 1
		return
	}

	err = rootCmd.ExecuteContext(ctx)

	defer func() {
		err = db.Close()
	}()

	if err != nil {
		fmt.Println("Error executing command failed")
		fmt.Print(err)
		exitCode = 1
		return
	}
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
