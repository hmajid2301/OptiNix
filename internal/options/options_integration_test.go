//go:build integration
// +build integration

package options_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gitlab.com/majiy00/go/clis/optinix/internal/options"
	"gitlab.com/majiy00/go/clis/optinix/internal/options/store"
)

func TestIntegrationSaveOptions(t *testing.T) {
	t.Run("Should save options", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)

		var count int64
		db.Model(&store.Option{}).Group("name").Count(&count)
		assert.Greater(t, count, int64(1))
	})
	t.Run("Should not save options because latest in db not a week old", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(context.Background())
		assert.NoError(t, err)

		var count int64
		db.Model(&store.Option{}).Group("name").Count(&count)
		assert.Greater(t, count, int64(1))
	})
}

func TestIntegrationGetOptions(t *testing.T) {
	t.Run("Should get option with `name` in option name", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
