package store_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func TestIntegrationAddOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should add options to DB successfully", func(t *testing.T) {
		ctx := context.Background()
		db := createDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		optionsToAdd := []store.OptionWithSources{
			{
				Name:         "option",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example.com"},
			},
			{
				Name:         "option1",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example.com"},
			},
		}
		err = s.AddOptions(ctx, optionsToAdd)
		assert.NoError(t, err)

		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM options").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count, "Two entries should have been added to table")
	})

	t.Run("Should get options from DB successfully", func(t *testing.T) {
		ctx := context.Background()
		db := createDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		optionsToAdd := []store.OptionWithSources{
			{
				Name:         "option",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example.com"},
			},
			{
				Name:         "option1",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example.com"},
			},
			{
				Name:         "other_name",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example.com"},
			},
		}
		err = s.AddOptions(ctx, optionsToAdd)
		assert.NoError(t, err)

		options, err := s.FindOptions(ctx, "option")
		assert.NoError(t, err)
		assert.Len(t, options, 2)
	})
}

func createDB(ctx context.Context, t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	dir, err := os.Getwd()
	assert.NoError(t, err)
	schemaFile := filepath.Join(dir, "../", "../", "../", "db/schema.sql")
	content, err := os.ReadFile(schemaFile)
	assert.NoError(t, err)
	ddl := string(content)
	_, err = db.ExecContext(ctx, ddl)
	assert.NoError(t, err)
	return db
}
