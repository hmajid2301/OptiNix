package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"gitlab.com/hmajid2301/optinix/cmd"
)

//go:embed db/schema.sql
var ddl string

func main() {
	ctx := context.TODO()
	db, err := cmd.GetDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	cmd.Execute()
}
