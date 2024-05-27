package store_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func setupSubtest(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()
	db := optionstest.CreateDB(ctx, t)

	return db, func() {
		db.Close()
	}
}

func TestIntegrationAddOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should add options to DB successfully", func(t *testing.T) {
		db, teardown := setupSubtest(t)
		defer teardown()

		str, err := store.NewStore(db)
		assert.NoError(t, err)

		optionsToAdd := []entities.Option{
			{
				Name:        "option",
				Description: "description",
				Type:        "str",
				Default:     "default",
				Example:     "example",
				Sources:     []string{"http://example.com"},
			},
			{
				Name:        "option1",
				Description: "description",
				Type:        "str",
				Default:     "default",
				Example:     "example",
				Sources:     []string{"http://example1.com"},
			},
		}

		ctx := context.Background()
		err = str.AddOptions(ctx, optionsToAdd)
		assert.NoError(t, err)

		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM options").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count, "Two entries should have been added to table")

		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sources").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 2, count, "Two entries should have been added to table")
	})
}

func TestIntegrationFindOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should get options from DB successfully", func(t *testing.T) {
		db, teardown := setupSubtest(t)
		defer teardown()

		myStore, err := store.NewStore(db)
		assert.NoError(t, err)

		optionsToAdd := []entities.Option{
			{
				Name:        "option",
				Description: "description",
				Type:        "str",
				Default:     "default",
				Example:     "example",
				Sources:     []string{"http://example.com"},
			},
			{
				Name:        "option.enable",
				Description: "description",
				Type:        "str",
				Default:     "default",
				Example:     "example",
				Sources:     []string{"http://example1.com"},
			},
			{
				Name:        "other_name",
				Description: "description",
				Type:        "str",
				Default:     "default",
				Example:     "example",
				Sources:     []string{"http://example2.com"},
			},
		}
		ctx := context.Background()
		err = myStore.AddOptions(ctx, optionsToAdd)
		assert.NoError(t, err)

		options, err := myStore.FindOptions(ctx, "option", 10)
		assert.NoError(t, err)
		assert.Len(t, options, 2, "Two entries with the name `option` should've been found")
	})
}
