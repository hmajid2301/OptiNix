package options_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/majiy00/go/clis/optinix/internal/options"
	"gitlab.com/majiy00/go/clis/optinix/internal/options/store"
)

func TestIntegrationSaveOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should save options", func(t *testing.T) {
		db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Should not save options because latest in db not a week old", func(t *testing.T) {
		db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)
	})
}

func TestIntegrationGetOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should get option with `name` in option name", func(t *testing.T) {
		db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)

		nixOpts, err := opt.GetOptions(context.Background(), "name")
		assert.NoError(t, err)

		assert.Equal(t, len(nixOpts), 10)
	})
}
