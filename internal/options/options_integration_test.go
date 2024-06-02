package options_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type NixReader struct{}

func (f NixReader) Read(pathToExpression string) ([]byte, error) {
	if pathToExpression == "../../testdata/hm-options.json/share/doc/home-manager/options.json" {
		pathToExpression = "../../testdata/hm-options.json"
	}
	nixExpression, err := os.ReadFile(pathToExpression)
	return nixExpression, err
}

type NixCmdExecutor struct{}

func (n NixCmdExecutor) Execute(_ context.Context, expression string) (string, error) {
	switch expression {
	case "./nix/nixos-options.nix":
		return "../../testdata/nixos-options.json", nil
	case "./nix/hm-options.nix":
		return "../../testdata/hm-options.json", nil
	case "./nix/darwin-options.nix":
		return "../../testdata/darwin-options.json", nil
	}

	return "", nil
}

type Updater struct{}

func (Updater) SendMessage(_ string) {
}

func setupSubTest(t *testing.T) (options.Searcher, store.Store, func()) {
	ctx := context.Background()
	db := optionstest.CreateDB(ctx, t)
	dbStore, err := store.NewStore(db)
	assert.NoError(t, err)

	fetcher := fetch.NewFetcher(NixCmdExecutor{}, NixReader{}, Updater{})
	opt := options.NewSearcher(dbStore, fetcher)

	return opt, dbStore, func() {
		db.Close()
	}
}

func TestIntegrationSaveOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	sources := entities.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}
	forceRefresh := false

	t.Run("Should save options", func(t *testing.T) {
		opt, _, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.SaveOptions(ctx, sources, forceRefresh)
		assert.NoError(t, err)
	})

	t.Run("Should not fetch new options unless they are a day old", func(t *testing.T) {
		opt, _, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.SaveOptions(ctx, sources, forceRefresh)
		assert.NoError(t, err)

		shouldFetch, err := opt.ShouldFetch(ctx)
		assert.False(t, shouldFetch)
		assert.NoError(t, err)
	})

	t.Run("Should fetch new options if force refresh argument is passed", func(t *testing.T) {
		opt, store, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.SaveOptions(ctx, sources, forceRefresh)
		assert.NoError(t, err)
		lastUpdated, err := store.GetLastAddedTime(ctx)
		assert.NoError(t, err)

		// TODO: find a nicer to do this, time is only accurate to the second
		time.Sleep(1 * time.Second)
		shouldForceRefresh := true
		err = opt.SaveOptions(ctx, sources, shouldForceRefresh)
		assert.NoError(t, err)
		lastUpdated2, err := store.GetLastAddedTime(ctx)
		assert.NoError(t, err)

		assert.NotEqual(t, lastUpdated, lastUpdated2, "check that new rows were added to the database")
	})
}

func TestIntegrationGetOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	sources := entities.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}
	forceRefresh := false

	t.Run("Should get option with `vdirsyncer` in option name", func(t *testing.T) {
		opt, _, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.SaveOptions(ctx, sources, forceRefresh)
		assert.NoError(t, err)

		nixOpts, err := opt.GetOptions(ctx, "vdirsyncer enable", 10)
		assert.NoError(t, err)

		expectedResults := 2
		assert.Len(t, nixOpts, expectedResults)

		expectedOpt := entities.Option{
			Name:        "accounts.calendar.accounts.<name>.vdirsyncer.enable",
			Description: "Whether to enable synchronization using vdirsyncer.",
			Type:        "boolean",
			Default:     "false",
			Example:     "true",
			Sources: []string{
				"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
			},
			OptionFrom: "Home Manager",
		}
		assert.Contains(t, nixOpts, expectedOpt)
	})
}
