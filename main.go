package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"gitlab.com/hmajid2301/optinix/cmd"
)

//go:embed db/schema.sql
var ddl string

//go:embed nix
var _ embed.FS

func main() {
	var exitCode int
	ctx := gracefulShutdown()
	defer func() { os.Exit(exitCode) }()

	rootCmd, err := cmd.NewRootCmd(ctx, ddl)
	if err != nil {
		fmt.Println("Error creating root command")
		fmt.Print(err)
		exitCode = 1
		return
	}

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Println("Error executing command failed")
		fmt.Print(err)
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
