package options

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/hmajid2301/optinix/internal/options/fetcher"
	"gitlab.com/hmajid2301/optinix/internal/options/parser"
	"gitlab.com/hmajid2301/optinix/internal/options/store"
)

type Source string

const (
	NixOSSource       Source = "nixos"
	HomeManagerSource Source = "home-manager"
)

type Opt struct {
	store store.Store
	// TODO: should this be in save option func
	httpRetries int
}

func New(s store.Store) Opt {
	return Opt{store: s}
}

func (o Opt) SaveOptions(ctx context.Context, sources map[Source]string) error {
	shouldFetch, err := o.shouldFetch(ctx)
	if err != nil {
		return err
	}

	if !shouldFetch {
		return nil
	}

	fetch := fetcher.NewFetcher(o.httpRetries)
	for source := range sources {
		url := sources[source]
		html, err := fetch.Fetch(ctx, url)
		if err != nil {
			return err
		}
		options := parser.Parse(html)

		storeOptions := getStoreOptions(options)
		batchSize := 100
		for i := 0; i < len(storeOptions)/batchSize; i++ {
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
	}

	return nil
}

func getStoreOptions(options []parser.Option) []store.OptionWithSources {
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

func (o Opt) GetOptions(ctx context.Context, name string) ([]store.OptionWithSources, error) {
	return o.store.FindOptions(ctx, name)
}

// TODO: move this out from save options bit, what if user wants to a force a refresh using CLI
func (o Opt) shouldFetch(ctx context.Context) (bool, error) {
	lastUpdatedDB := time.Now()
	time, err := o.store.GetLastAddedTime(ctx)
	if err != nil {
		// Likely first time the CLI has ran, as DB is empty
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	nextWeek := lastUpdatedDB.AddDate(0, 1, 0)
	if nextWeek.Before(time) {
		return false, nil
	}
	return true, nil
}
