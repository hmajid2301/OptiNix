package fetcher_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/fetcher"
)

func TestFetch(t *testing.T) {
	t.Run("Should fetch NixOS HTML", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://nixos.org").
			Get("/manual/nixos/unstable/options").
			Reply(http.StatusOK).
			File("../../../testdata/nixos_options.html")

		ctx := context.Background()

		fetch := fetcher.NewFetcher(0)
		gock.InterceptClient(fetch.Client)

		doc, err := fetch.Fetch(ctx, "https://nixos.org/manual/nixos/unstable/options")
		assert.NoError(t, err)

		assert.NotNil(t, doc)
		assert.Equal(t, gock.IsDone(), true)
	})

	t.Run("Should fetch Home Manager HTML", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://nix-community.github.io").
			Get("/home-manager/options.html").
			Reply(http.StatusOK).
			File("../../../testdata/home_manager_options.html")

		ctx := context.Background()
		fetch := fetcher.NewFetcher(0)
		gock.InterceptClient(fetch.Client)
		doc, err := fetch.Fetch(ctx, "https://nix-community.github.io/home-manager/options.html")
		assert.NoError(t, err)

		assert.NotNil(t, doc)
		assert.Equal(t, gock.IsDone(), true)
	})

	t.Run("Should throw error because of http error", func(t *testing.T) {
		defer gock.Off()

		gock.New("https://nix-community.github.io").
			Get("/home-manager/options.html").
			Reply(http.StatusInternalServerError)

		ctx := context.Background()
		fetch := fetcher.NewFetcher(0)
		gock.InterceptClient(fetch.Client)
		_, err := fetch.Fetch(ctx, "https://nix-community.github.io/home-manager/options.html")
		assert.Error(t, err)
	})
}
