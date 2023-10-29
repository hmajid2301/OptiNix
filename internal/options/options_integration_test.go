//go:build integration
// +build integration

package options_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

func TestIntegrationOptions(t *testing.T) {
	t.Run("Should save options", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		err = options.SaveOptions(context.Background(), db)
		assert.NoError(t, err)

		var count int64
		db.Model(&store.Option{}).Group("name").Count(&count)
		assert.Greater(t, count, int64(1))
	})
}
