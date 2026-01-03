package options

import (
	"context"
	"fmt"
	"log/slog"

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
	logger    *slog.Logger
}

func NewSearcher(s store.Store, f Fetcherer, m entities.Messager, logger *slog.Logger) Searcher {
	return Searcher{store: s, fetcher: f, messenger: m, logger: logger}
}

func (s Searcher) GetAll(ctx context.Context) ([]entities.Option, error) {
	return s.store.GetAllOptions(ctx)
}

func (s Searcher) Find(ctx context.Context, name string) ([]entities.Option, error) {
	return s.store.FindOptions(ctx, name)
}

func (s Searcher) Count(ctx context.Context) (int64, error) {
	return s.store.CountOptions(ctx)
}

func (s Searcher) Save(ctx context.Context, sources entities.Sources) error {
	s.messenger.Send("Clearing existing options from DB")
	err := s.store.ClearAllOptions(ctx)
	if err != nil {
		return err
	}

	options, err := s.fetcher.Fetch(ctx, sources)
	if err != nil {
		return err
	}

	s.logger.Info("Fetched options from sources", "count", len(options))
	batchSize := 1000
	maxBatches := (len(options) / batchSize) + 1

	for i := 0; i < maxBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(options) {
			end = len(options)
		}

		opts := options[start:end]
		s.messenger.Send(fmt.Sprintf("Saving batch %d/%d to database...", i+1, maxBatches))
		s.logger.Debug("Saving batch", "batch", i+1, "total_batches", maxBatches, "options_count", len(opts))
		err = s.store.AddOptions(ctx, opts)
		if err != nil {
			return err
		}
	}

	s.logger.Info("Successfully saved options to database", "count", len(options))
	return nil
}
