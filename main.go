package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"gitlab.com/hmajid2301/optinix/cmd"
	"gitlab.com/hmajid2301/optinix/internal/options/config"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

//go:embed db/schema.sql
var ddl string

func main() {
	var exitCode int
	ctx := gracefulShutdown()
	defer func() { os.Exit(exitCode) }()

	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config")
		fmt.Print(err)
		exitCode = 1
		return
	}

	logger, err := initLogger(conf.DBFolder)
	if err != nil {
		fmt.Println("Error initializing logger")
		fmt.Print(err)
		exitCode = 1
		return
	}

	db, err := getDB(ctx, ddl)
	if err != nil {
		fmt.Println("Error creating db command")
		fmt.Print(err)
		exitCode = 1
		return
	}

	rootCmd, err := cmd.NewRootCmd(ctx, db, logger)
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

func initLogger(logDir string) (*slog.Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, "optinix.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return slog.New(slog.NewTextHandler(file, nil)), nil
}

func getDB(ctx context.Context, ddl string) (*sql.DB, error) {
	args := os.Args
	if len(args) > 0 && args[1] == "completion" {
		return nil, nil
	}

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
