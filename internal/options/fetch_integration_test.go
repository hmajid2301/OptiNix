package options_test

import (
	"testing"
)

// TODO: mock out using docker container to capture request?
// TODO: Move this logic to E2E
func TestIntegrationFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		t.Skip("skipping integration test")
		// fetcher := options.NewFetcher(3)
		// html, err := fetcher.Fetch(context.Background(), "https://nixos.org/manual/nixos/unstable/options")
		// assert.NoError(t, err)
		//
		// dtCount := 0
		// doc := goquery.NewDocumentFromNode(html)
		// doc.Find("dt").Each(
		// 	func(_ int, _ *goquery.Selection) {
		// 		dtCount++
		// 	},
		// )
		//
		// assert.Greater(t, dtCount, 0)
	})
}
