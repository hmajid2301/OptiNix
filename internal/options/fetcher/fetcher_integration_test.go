package fetcher_test

import (
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/fetcher"
)

// TODO: mock out using docker container to capture request?
// TODO: Move this logic to E2E
func TestIntegrationFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		fetch := fetcher.NewFetcher(3)
		html, err := fetch.Fetch(context.Background(), "https://nixos.org/manual/nixos/unstable/options")
		assert.NoError(t, err)

		dtCount := 0
		doc := goquery.NewDocumentFromNode(html)
		doc.Find("dt").Each(
			func(_ int, _ *goquery.Selection) {
				dtCount++
			},
		)

		assert.Greater(t, dtCount, 0)
	})
}
