package store_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type TestStoreSuite struct {
	suite.Suite
	db *sql.DB
}

func TestIntegrationStore(t *testing.T) {
	ctx := context.Background()
	db := optionstest.CreateDB(ctx, t)

	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, &TestStoreSuite{db: db})
}

func (s *TestStoreSuite) SetUpSubTest() {
	ctx := context.Background()
	s.db = optionstest.CreateDB(ctx, s.T())
	s.T().Cleanup(func() {
		s.db.Close()
	})
}

func (s *TestStoreSuite) TestIntegrationAddOptions() {
	s.Run("Should add options to DB successfully", func() {
		ctx := context.Background()

		str, err := store.NewStore(s.db)
		s.NoError(err)

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
				Sources:      []string{"http://example1.com"},
			},
		}
		err = str.AddOptions(ctx, optionsToAdd)
		s.NoError(err)

		var count int
		err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM options").Scan(&count)
		s.NoError(err)
		s.Equal(2, count, "Two entries should have been added to table")

		err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sources").Scan(&count)
		s.NoError(err)
		s.Equal(2, count, "Two entries should have been added to table")
	})
}

func (s *TestStoreSuite) TestIntegrationFindOptions() {
	s.Run("Should get options from DB successfully", func() {
		ctx := context.Background()

		str, err := store.NewStore(s.db)
		s.NoError(err)

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
				Sources:      []string{"http://example1.com"},
			},
			{
				Name:         "other_name",
				Description:  "description",
				Type:         "str",
				DefaultValue: "default",
				Example:      "example",
				Sources:      []string{"http://example2.com"},
			},
		}
		err = str.AddOptions(ctx, optionsToAdd)
		s.NoError(err)

		options, err := str.FindOptions(ctx, "option")
		s.NoError(err)
		s.Len(options, 2, "Two entries with the name `option` should've been found")
	})
}
