package options

import (
	"context"

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

func (s Searcher) GetAll(ctx context.Context) ([]entities.Option, error) {
	return s.store.GetAllOptions(ctx)
}

func (s Searcher) Find(ctx context.Context, name string, limit int64) ([]entities.Option, error) {
	return s.store.FindOptions(ctx, name, limit)
}

func (s Searcher) Count(ctx context.Context) (int64, error) {
	return s.store.CountOptions(ctx)
}

func (s Searcher) Save(ctx context.Context, sources entities.Sources) error {
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
