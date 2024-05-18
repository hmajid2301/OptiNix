package options_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"gitlab.com/hmajid2301/optinix/internal/options"
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

func (n NixCmdExecutor) Executor(expression string) (string, error) {
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

type TestFetchSuite struct {
	suite.Suite
	store store.Store
	opt   options.Opt
}

func TestIntegrationStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, &TestFetchSuite{opt: options.Opt{}})
}

func (s *TestFetchSuite) SetupSubTest() {
	ctx := context.Background()
	db := optionstest.CreateDB(ctx, s.T())
	dbStore, err := store.NewStore(db)
	s.NoError(err)

	fetcher := options.NewFetcher(NixCmdExecutor{}, NixReader{})
	s.store = dbStore
	s.opt = options.NewOptions(dbStore, fetcher)

	s.T().Cleanup(func() {
		db.Close()
	})
}

func (s *TestFetchSuite) TestIntegrationSaveOptions() {
	sources := options.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}
	forceRefresh := false

	s.Run("Should save options", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources, forceRefresh)
		s.NoError(err)
	})

	s.Run("Should not fetch new options unless they are a day old", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources, forceRefresh)
		s.NoError(err)

		shouldFetch, err := options.ShouldFetch(ctx, s.opt)
		s.False(shouldFetch)
		s.NoError(err)
	})

	s.Run("Should fetch new options if force refresh argument is passed", func() {
		ctx := context.Background()

		err := s.opt.SaveOptions(ctx, sources, forceRefresh)
		s.NoError(err)
		lastUpdated, err := s.store.GetLastAddedTime(ctx)
		s.NoError(err)

		// TODO: find a nicer to do this, time is only accurate to the second
		time.Sleep(1 * time.Second)
		shouldForceRefresh := true
		err = s.opt.SaveOptions(ctx, sources, shouldForceRefresh)
		s.NoError(err)
		lastUpdated2, err := s.store.GetLastAddedTime(ctx)
		s.NoError(err)

		s.NotEqual(lastUpdated, lastUpdated2, "check that new rows were added to the database")
	})
}

func (s *TestFetchSuite) TestIntegrationGetOptions() {
	sources := options.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}
	forceRefresh := false

	s.Run("Should get option with `vdirsyncer` in option name", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources, forceRefresh)
		s.NoError(err)

		nixOpts, err := s.opt.GetOptions(ctx, "vdirsyncer enable")
		s.NoError(err)

		expectedResults := 2
		s.Len(nixOpts, expectedResults)

		opt := store.OptionWithSources{
			Name:         "accounts.calendar.accounts.<name>.vdirsyncer.enable",
			Description:  "Whether to enable synchronization using vdirsyncer.",
			Type:         "boolean",
			DefaultValue: "false",
			Example:      "true",
			Sources: []string{
				"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
			},
		}
		s.Contains(nixOpts, opt)
	})
}
