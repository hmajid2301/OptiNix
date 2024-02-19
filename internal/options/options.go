package options

import (
	"context"
	"time"

	"gitlab.com/majiy00/go/clis/optinix/internal/options/store"
)

type Source string

const (
	NixOSSource       Source = "nixos"
	HomeManagerSource Source = "home-manager"
)

var (
	defaultHTTPRetries = 3
	sources            = map[Source]string{
		NixOSSource:       "https://nixos.org/manual/nixos/unstable/options",
		HomeManagerSource: "https://nix-community.github.io/home-manager/options.html",
	}
)

type Opt struct {
	store store.Store
}

func New(s store.Store) Opt {
	return Opt{store: s}
}

func (o Opt) SaveOptions(ctx context.Context) error {
	shouldFetch, err := o.shouldFetch(ctx)
	if err != nil {
		return err
	}

	if !shouldFetch {
		return nil
	}

	fetcher := NewFetcher(defaultHTTPRetries)
	for source := range sources {
		url := sources[source]
		html, err := fetcher.Fetch(ctx, url)
		if err != nil {
			return err
		}
		options := Parse(html)

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

func (o Opt) GetOptions(ctx context.Context, name string) ([]store.OptionWithSources, error) {
	return o.store.FindOptions(ctx, name)
}

func (o Opt) shouldFetch(ctx context.Context) (bool, error) {
	lastUpdatedDB := time.Now()
	time, err := o.store.GetLastAddedTime(ctx)
	if err != nil {
		return false, err
	}

	nextWeek := lastUpdatedDB.AddDate(0, 1, 0)
	if nextWeek.Before(time) {
		return false, nil
	}
	return true, nil
}
