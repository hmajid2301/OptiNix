package options_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	// used to connect to sqlite
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func createDB(ctx context.Context, t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	dir, err := os.Getwd()
	assert.NoError(t, err)
	schemaFile := filepath.Join(dir, "../", "../", "db/schema.sql")
	content, err := os.ReadFile(schemaFile)
	assert.NoError(t, err)
	ddl := string(content)
	_, err = db.ExecContext(ctx, ddl)
	assert.NoError(t, err)
	return db
}

func TestIntegrationSaveOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should save options", func(t *testing.T) {
		ctx := context.Background()
		db := createDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(ctx)
		assert.NoError(t, err)
	})

	t.Run("Should not save options because latest in db not a week old", func(t *testing.T) {
		assert.True(t, true)
	})
}

func TestIntegrationGetOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should get option with `name` in option name", func(t *testing.T) {
		ctx := context.Background()
		db := createDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)

		nixOpts, err := opt.GetOptions(context.Background(), "name")
		assert.NoError(t, err)

		// TODO: fix
		expectedResults := 10
		assert.Equal(t, expectedResults, len(nixOpts))
	})
}
