//go:build !integration
// +build !integration

package store_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestAddOption(t *testing.T) {
	t.Run("Should successfully add an option to database", func(t *testing.T) {
		option := store.Option{
			Name:        "home.enable",
			Description: "Whether to enable home",
			Type:        "boolean",
			Default:     "false",
			Sources: []store.Source{
				{
					URL: "https://example.com",
				},
			},
		}
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "options" ("created_at","updated_at","deleted_at","name","description","type","default","example")
		  		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "home.enable", "Whether to enable home", "boolean", "false", "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "sources" ("created_at","updated_at","deleted_at","url") 
				VALUES ($1,$2,$3,$4) ON CONFLICT ("id") DO UPDATE SET "id"="excluded"."id" RETURNING "id"`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "https://example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectCommit()

		s := store.New(gormDB)
		err = s.AddOptions(context.Background(), []*store.Option{&option})
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Should fail add an option with same name to store", func(t *testing.T) {
		option := store.Option{
			Name:        "home.enable",
			Description: "Whether to enable home",
			Type:        "boolean",
			Default:     "false",
			Sources: []store.Source{
				{
					URL: "https://example.com",
				},
			},
		}
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		assert.NoError(t, err)

		s := store.New(gormDB)
		err = s.AddOptions(context.Background(), []*store.Option{&option, &option})
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
