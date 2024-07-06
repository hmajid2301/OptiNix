package options

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type Source string

type Fetcherer interface {
	Fetch(ctx context.Context, sources entities.Sources) ([]entities.Option, error)
}

type Searcher struct {
	fetcher   Fetcherer
	store     store.Store
	messenger entities.Messager
}

func NewSearcher(s store.Store, f Fetcherer, m entities.Messager) Searcher {
	return Searcher{store: s, fetcher: f, messenger: m}
}

func (s Searcher) GetAllOptions(ctx context.Context) ([]entities.Option, error) {
	return s.store.GetAllOptions(ctx)
}

func (s Searcher) FindOptions(ctx context.Context, name string, limit int64) ([]entities.Option, error) {
	return s.store.FindOptions(ctx, name, limit)
}

func (s Searcher) SaveOptions(ctx context.Context, sources entities.Sources, forceRefresh bool) error {
	// INFO: If force fresh is passed then we should always fetch the options and update the store.
	// Else we check if should update the options in the store, if they have gone stale i.e. older than a day.
	if !forceRefresh {
		shouldFetch, err := s.shouldFetch(ctx)
		if err != nil {
			return err
		}

		if !shouldFetch {
			return nil
		}
	}

	options, err := s.fetcher.Fetch(ctx, sources)
	if err != nil {
		return err
	}

	s.messenger.Send("Trying to save options to DB")
	batchSize := 1000
	maxBatches := (len(options) / batchSize) + 1

	for i := 0; i < maxBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(options) {
			end = len(options)
		}

		opts := options[start:end]
		err = s.store.AddOptions(ctx, opts)
		if err != nil {
			return err
		}
	}

	s.messenger.Send("Successfully to save options to DB")
	return nil
}

func (s Searcher) shouldFetch(ctx context.Context) (bool, error) {
	now := time.Now()
	lastUpdatedDB, err := s.store.GetLastAddedTime(ctx)
	if err != nil {
		// INFO: Likely first time the CLI has ran, as DB is empty.
		// So we should fetch the options.
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	// INFO: If its been a day since we fetched the options, lets update the DB again.
	x := lastUpdatedDB.AddDate(0, 0, 1)
	if x.After(now) {
		return false, nil
	}
	return true, nil
}
