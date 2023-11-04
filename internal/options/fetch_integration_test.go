//go:build integration
// +build integration

package options_test

import (
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"gitlab.com/majiy00/go/clis/optinix/internal/options"
)

// TODO: mock out using docker container to capture request?
// TODO: Move this logic to E2E
func TestIntegrationFetch(t *testing.T) {
	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		fetcher := options.NewFetcher(3)
		html, err := fetcher.Fetch(context.Background(), "https://nixos.org/manual/nixos/unstable/options")
		assert.NoError(t, err)

		dtCount := 0
		doc := goquery.NewDocumentFromNode(html)
		doc.Find("dt").Each(
			func(_ int, _ *goquery.Selection) {
				dtCount += 1
			},
		)

		assert.Greater(t, dtCount, 0)
	})
}
