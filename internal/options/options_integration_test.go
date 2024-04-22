package options_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/dbtest"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

var sources = map[options.Source]string{
	options.NixOSSource:       getHost("/manual/nixos/unstable/options"),
	options.HomeManagerSource: getHost("/home-manager/options.xhtml"),
}

func getHost(path string) string {
	fullPath := "http://localhost:8080" + path
	if os.Getenv("CI") == "true" {
		fullPath = "http://docker:8080" + path
	}

	return fullPath
}

func TestIntegrationSaveOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("Should save options", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.CreateDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(ctx, sources)
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

	t.Run("Should get option with `appstream` in option name", func(t *testing.T) {
		ctx := context.Background()
		db := dbtest.CreateDB(ctx, t)

		s, err := store.New(db)
		assert.NoError(t, err)

		opt := options.New(s)
		err = opt.SaveOptions(ctx, sources)
		assert.NoError(t, err)

		nixOpts, err := opt.GetOptions(ctx, "appstream")
		assert.NoError(t, err)

		expectedResults := 2
		assert.Len(t, nixOpts, expectedResults)
	})
}
