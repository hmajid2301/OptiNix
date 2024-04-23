package options_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

var sources = map[options.Source]string{
	options.NixOSSource:       optionstest.GetHost("/manual/nixos/unstable/options"),
	options.HomeManagerSource: optionstest.GetHost("/home-manager/options.xhtml"),
}

type TestOptionSuite struct {
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
	str, err := store.New(db)
	assert.NoError(t, err)
	opt := options.New(str)

	suite.Run(t, &TestOptionSuite{opt: opt})
}

func (s *TestOptionSuite) SetUpSubTest() {
	ctx := context.Background()
	s.db = optionstest.CreateDB(ctx, s.T())
	s.T().Cleanup(func() {
		s.db.Close()
	})
}

func (s *TestOptionSuite) TestIntegrationSaveOptions() {
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

func (s *TestOptionSuite) TestIntegrationGetOptions() {
	s.Run("Should get option with `appstream` in option name", func() {
		ctx := context.Background()
		err := s.opt.SaveOptions(ctx, sources)
		s.NoError(err)

		nixOpts, err := s.opt.GetOptions(ctx, "appstream")
		s.NoError(err)

		expectedResults := 1
		s.Len(nixOpts, expectedResults)
	})
}
