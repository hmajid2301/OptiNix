package fetcher_test

import (
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/fetcher"
	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
)

func TestIntegrationFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		fetch := fetcher.NewFetcher(3)
		nixosHost := optionstest.GetHost("/manual/nixos/unstable/options")
		html, err := fetch.Fetch(context.Background(), nixosHost)
		assert.NoError(t, err)

		dtCount := 0
		doc := goquery.NewDocumentFromNode(html)
		doc.Find("dt").Each(
			func(_ int, _ *goquery.Selection) {
				dtCount++
			},
		)

		assert.Equal(t, dtCount, 3)
	})

	t.Run("Should fetch home-manager HTML", func(t *testing.T) {
		fetch := fetcher.NewFetcher(3)
		nixosHost := optionstest.GetHost("/home-manager/options.xhtml")
		html, err := fetch.Fetch(context.Background(), nixosHost)
		assert.NoError(t, err)

		dtCount := 0
		doc := goquery.NewDocumentFromNode(html)
		doc.Find("dt").Each(
			func(_ int, _ *goquery.Selection) {
				dtCount++
			},
		)

		assert.Equal(t, dtCount, 3)
	})

	t.Run("Should not find HMTL page", func(t *testing.T) {
		fetch := fetcher.NewFetcher(3)
		nixosHost := optionstest.GetHost("/home-manager/invalid.xhtml")
		_, err := fetch.Fetch(context.Background(), nixosHost)
		assert.ErrorContains(t, err, "request failed with status code 404")
	})
}
