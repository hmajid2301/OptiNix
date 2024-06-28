package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	sqlc "gitlab.com/hmajid2301/optinix/internal/options/store/db"
)

type Store struct {
	db      *sql.DB
	queries *sqlc.Queries
}

var SearchLimit = 10

func NewStore(db *sql.DB) (Store, error) {
	queries := sqlc.New(db)
	store := Store{
		db:      db,
		queries: queries,
	}

	return store, nil
}

func (s Store) AddOptions(ctx context.Context, options []entities.Option) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()

	for _, option := range options {
		addOption := sqlc.AddOptionParams{
			OptionName:   option.Name,
			Description:  option.Description,
			OptionType:   option.Type,
			OptionFrom:   option.OptionFrom,
			DefaultValue: option.Default,
			Example:      option.Example,
		}

		newOption, err := s.queries.WithTx(tx).AddOption(ctx, addOption)
		if err != nil {
			return err
		}

		for _, source := range option.Sources {
			newSource, err := s.queries.WithTx(tx).AddSource(ctx, source)
			if err != nil {
				return err
			}

			addSourceOption := sqlc.AddSourceOptionParams{
				SourceID: newSource.ID,
				OptionID: newOption.ID,
			}
			_, err = s.queries.WithTx(tx).AddSourceOption(ctx, addSourceOption)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s Store) GetAllOptions(ctx context.Context) ([]entities.Option, error) {
	options := []entities.Option{}
	opts, err := s.queries.GetAllOptions(ctx)
	if err != nil {
		return options, err
	}

	for _, opt := range opts {
		sources := ""
		sources, _ = opt.SourceList.(string)
		sourceList := strings.Split(sources, ",")

		option := entities.Option{
			Name:        opt.OptionName,
			Description: opt.Description,
			Type:        opt.OptionType,
			OptionFrom:  opt.OptionFrom,
			Default:     opt.DefaultValue,
			Example:     opt.Example,
			Sources:     sourceList,
		}
		options = append(options, option)
	}
	return options, nil
}

func (s Store) FindOptions(ctx context.Context, name string, limit int64) ([]entities.Option, error) {
	options := []entities.Option{}
	opts, err := s.queries.FindOptions(ctx, sqlc.FindOptionsParams{
		OptionName: fmt.Sprintf("\"%s\"", name),
		Limit:      limit,
	})
	if err != nil {
		return options, err
	}

	for _, opt := range opts {
		sources := ""
		sources, _ = opt.SourceList.(string)
		sourceList := strings.Split(sources, ",")

		option := entities.Option{
			Name:        opt.OptionName,
			Description: opt.Description,
			Type:        opt.OptionType,
			OptionFrom:  opt.OptionFrom,
			Default:     opt.DefaultValue,
			Example:     opt.Example,
			Sources:     sourceList,
		}
		options = append(options, option)
	}
	return options, nil
}

func (s Store) GetLastAddedTime(ctx context.Context) (time.Time, error) {
	result, err := s.queries.GetLastUpdated(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return result.Time, nil
}

func GetDB(dbFolder string) (*sql.DB, error) {
	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		permissions := 0755
		err = os.Mkdir(dbFolder, fs.FileMode(permissions))
		if err != nil {
			return nil, err
		}
	}

	dbPath := filepath.Join(dbFolder, "options.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	return db, err
}
