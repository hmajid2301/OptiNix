package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"gitlab.com/hmajid2301/optinix/cmd"
	"gitlab.com/hmajid2301/optinix/internal/options/config"
)

//go:embed db/schema.sql
var ddl string

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	db, err := cmd.GetDB(conf.DBFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	ctx := context.TODO()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	cmd.Execute()
}
