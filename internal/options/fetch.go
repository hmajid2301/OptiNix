package options

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/html"
)

var (
	timeout = time.Second * 10
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

func (f Fetcher) Fetch(ctx context.Context, url string) (*html.Node, error) {
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
		return nil, fmt.Errorf("request failed with status code %d", response.StatusCode)
	}

	defer response.Body.Close()

	return html.Parse(response.Body)
}
