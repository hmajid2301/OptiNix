//go:build integration
// +build integration

package store_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gitlab.com/majiy00/go/clis/optinix/internal/options/store"
)

func TestIntegrationAddOptionsToDB(t *testing.T) {
	t.Run("Should add multiple options to DB", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

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

func TestIntegrationFindOptions(t *testing.T) {
	t.Run("Should add retrieve multiple options", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		assert.NoError(t, err)

		s, err := store.New(db)
		assert.NoError(t, err)

		startingOpts := []*store.Option{
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

		err = s.AddOptions(context.Background(), startingOpts)
		assert.NoError(t, err)

		options, err := s.FindOptions(context.Background(), "home.enable")
		assert.NoError(t, err)

		for _, option := range options {
			assert.Contains(t, option.Name, "home.enable")
		}
		assert.Equal(t, len(options), 2)
	})
}
