//go:build integration
// +build integration

package store_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
	"gorm.io/gorm"
)

func TestIntegrationAddOptionsToDB(t *testing.T) {
	t.Run("Should add multiple options to DB", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		db.AutoMigrate(&store.Option{}, &store.Source{})
		s := store.New(db)

		options := []*store.Option{
			{
				Name:        "home.enable",
				Description: "Whether to enable home",
				Type:        "boolean",
				Default:     "false",
				Sources: []store.Source{
					{
						URL: "https://example.com",
					},
				},
			},
			{
				Name:        "accounts.calendar.accounts",
				Description: "List of calendars.",
				Type:        "attribute set of (submodule)",
				Default:     "{}",
				Sources: []store.Source{
					{
						URL: "https://github.com/nix-community/home-manager/blob/master/modules/programs/qcal.nix",
					},
					{
						URL: "https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
					},
				},
			},
		}

		err = s.AddOptions(context.Background(), options)
		assert.NoError(t, err)
	})
}
