package options_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	db  *sql.DB
	opt options.Opt
}

func TestIntegrationStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	db := optionstest.CreateDB(ctx, t)
	dbStore, err := store.NewStore(db)
	assert.NoError(t, err)

	fetcher := options.NewFetcher(NixCmdExecutor{}, NixReader{})
	opt := options.NewOptions(dbStore, fetcher)

	suite.Run(t, &TestFetchSuite{opt: opt})
}

func (s *TestFetchSuite) SetUpSubTest() {
	ctx := context.Background()
	s.db = optionstest.CreateDB(ctx, s.T())
	s.T().Cleanup(func() {
		s.db.Close()
	})
}

func (s *TestFetchSuite) TestIntegrationSaveOptions() {
	sources := options.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}

	s.Run("Should save options", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources)
		s.NoError(err)
	})

	// TODO: Make this an actual
	s.Run("Should not save options because latest in db not a week old", func() {
		s.True(true)
	})
}

func (s *TestFetchSuite) TestIntegrationGetOptions() {
	sources := options.Sources{
		NixOS:       "./nix/nixos-options.nix",
		HomeManager: "./nix/hm-options.nix",
		Darwin:      "./nix/darwin-options.nix",
	}

	s.Run("Should get option with `vdirsyncer` in option name", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources)
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
