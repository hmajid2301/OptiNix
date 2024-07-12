package options_test

import (
	"context"
	"os"
	"testing"

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

func (Updater) Send(_ string) {}

func setupSubTest(t *testing.T) (options.Searcher, func()) {
	ctx := context.Background()
	db := optionstest.CreateDB(ctx, t)
	dbStore, err := store.NewStore(db)
	assert.NoError(t, err)

	fetcher := fetch.NewFetcher(NixCmdExecutor{}, NixReader{}, Updater{})
	opt := options.NewSearcher(dbStore, fetcher, Updater{})

	return opt, func() {
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

	t.Run("Should save options", func(t *testing.T) {
		opt, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.Save(ctx, sources)
		assert.NoError(t, err)
		count, err := opt.Count(ctx)
		assert.Greater(t, count, int64(0))
		assert.NoError(t, err)
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

	t.Run("Should get option with `vdirsyncer` in option name", func(t *testing.T) {
		opt, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.Save(ctx, sources)
		assert.NoError(t, err)

		nixOpts, err := opt.Find(ctx, "vdirsyncer enable", 10)
		assert.NoError(t, err)

		expectedResults := 1
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
		assert.Equal(t, nixOpts[0], expectedOpt)
	})

	t.Run("Should get all options", func(t *testing.T) {
		opt, teardown := setupSubTest(t)
		defer teardown()

		ctx := context.Background()
		err := opt.Save(ctx, sources)
		assert.NoError(t, err)

		nixOpts, err := opt.GetAll(ctx)
		assert.NoError(t, err)

		expectedResults := 168
		assert.Len(t, nixOpts, expectedResults)
	})
}
