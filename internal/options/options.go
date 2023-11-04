package options

import (
	"context"

	"gorm.io/gorm"

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

func SaveOptions(ctx context.Context, db *gorm.DB) error {
	fetcher := NewFetcher(defaultHTTPRetries)
	for source := range sources {
		url := sources[source]
		html, err := fetcher.Fetch(ctx, url)
		if err != nil {
			return err
		}

		options := Parse(html)
		s, err := store.New(db)
		if err != nil {
			return err
		}

		storeOptions := getStoreOptions(options)
		batchSize := 100
		for i := 0; i < len(storeOptions)/batchSize; i++ {
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(storeOptions) {
				end = len(storeOptions)
			}

			opts := storeOptions[start:end]
			err = s.AddOptions(ctx, opts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getStoreOptions(options []Option) []*store.Option {
	matchingOptions := []*store.Option{}
	for _, option := range options {
		storeSources := []store.Source{}
		for _, source := range option.Sources {
			storeSource := store.Source{
				URL: source,
			}
			storeSources = append(storeSources, storeSource)
		}

		storeOption := store.Option{
			Name:        option.Name,
			Description: option.Description,
			Type:        option.Type,
			Default:     option.Default,
			Example:     option.Example,
			Sources:     storeSources,
		}

		matchingOptions = append(matchingOptions, &storeOption)
	}
	return matchingOptions
}
