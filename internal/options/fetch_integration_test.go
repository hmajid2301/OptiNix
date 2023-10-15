//go:build integration
// +build integration

package options_test

import (
	"context"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"gitlab.com/hmajid2301/optinix/internal/options"
)

func TestIntegrationFetch(t *testing.T) {
	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		ctx := context.Background()
		doc, err := options.Fetch(ctx, options.NixOSSource)
		assert.NoError(t, err)

		dtCount := 0
		doc_ := goquery.NewDocumentFromNode(doc)
		doc_.Find("dt").Each(
			func(_ int, _ *goquery.Selection) {
				dtCount += 1
			},
		)

		assert.Greater(t, dtCount, 0)
	})
}
