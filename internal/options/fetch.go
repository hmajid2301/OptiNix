package options

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/html"
)

type Source string

const (
	NixOSSource       Source = "nixos"
	HomeManagerSource Source = "home-manager"
)

var (
	timeout = time.Second * 5
	sources = map[Source]string{
		NixOSSource:       "https://nixos.org/manual/nixos/unstable/options",
		HomeManagerSource: "https://nix-community.github.io/home-manager/options.html",
	}
)

type Fetcher struct {
	Client *http.Client
}

func NewFetcher(maxRetries int) Fetcher {
	client := retryablehttp.NewClient()
	client.RetryMax = maxRetries
	std := client.StandardClient()
	return Fetcher{Client: std}
}

func (f Fetcher) Fetch(ctx context.Context, source Source) (*html.Node, error) {
	url, ok := sources[source]
	if !ok {
		return nil, fmt.Errorf("invalid source name %s", source)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := f.Client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to get options failed with status code %d", response.StatusCode)
	}

	defer response.Body.Close()

	return html.Parse(response.Body)
}
