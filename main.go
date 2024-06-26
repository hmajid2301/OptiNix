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
	if err := cmd.Execute(ctx, ddl); err != nil {
		fmt.Println("Error executing command falfd")
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
