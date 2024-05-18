package options

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type Source string

type Fetcherer interface {
	Fetch(ctx context.Context, sources Sources) ([]Option, error)
}

type Opt struct {
	fetcher Fetcherer
	store   store.Store
}

func NewOptions(s store.Store, f Fetcherer) Opt {
	return Opt{store: s, fetcher: f}
}

func (o Opt) SaveOptions(ctx context.Context, sources Sources, forceRefresh bool) error {
	// INFO: If force fresh is passed then we should always fetch the options and update the store.
	// Else we check if should update the options in the store, if they have gone stale i.e. older than a day.
	if !forceRefresh {
		shouldFetch, err := o.shouldFetch(ctx)
		if err != nil {
			return err
		}

		if !shouldFetch {
			return nil
		}
	}

	options, err := o.fetcher.Fetch(ctx, sources)
	if err != nil {
		return err
	}

	storeOptions := getStoreOptions(options)
	batchSize := 100
	maxBatches := (len(storeOptions) / batchSize) + 1

	for i := 0; i < maxBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(storeOptions) {
			end = len(storeOptions)
		}

		opts := storeOptions[start:end]
		err = o.store.AddOptions(ctx, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func getStoreOptions(options []Option) []store.OptionWithSources {
	matchingOptions := []store.OptionWithSources{}
	for _, option := range options {
		storeOption := store.OptionWithSources{
			Name:         option.Name,
			Description:  option.Description,
			Type:         option.Type,
			DefaultValue: option.Default,
			Example:      option.Example,
			Sources:      option.Sources,
		}

		matchingOptions = append(matchingOptions, storeOption)
	}
	return matchingOptions
}

// TODO: Common options struct here
func (o Opt) GetOptions(ctx context.Context, name string) ([]store.OptionWithSources, error) {
	return o.store.FindOptions(ctx, name)
}

// TODO: move this out from save options bit, what if user wants to a force a refresh using CLI
func (o Opt) shouldFetch(ctx context.Context) (bool, error) {
	now := time.Now()
	lastUpdatedDB, err := o.store.GetLastAddedTime(ctx)
	if err != nil {
		// INFO: Likely first time the CLI has ran, as DB is empty
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	// INFO: If its been a day since we fetched the options, lets update the DB again
	x := lastUpdatedDB.AddDate(0, 0, 1)
	if x.After(now) {
		return false, nil
	}
	return true, nil
}
