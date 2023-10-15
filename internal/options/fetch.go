package options

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type Source string

const (
	NixOSSource       Source = "nixos"
	HomeManagerSource Source = "home-manager"
)

var sources = map[Source]string{
	NixOSSource:       "https://nixos.org/manual/nixos/unstable/options",
	HomeManagerSource: "https://nix-community.github.io/home-manager/options.html",
}

func Fetch(ctx context.Context, source Source) (*html.Node, error) {
	url, ok := sources[source]
	if !ok {
		return nil, fmt.Errorf("invalid source name %s", source)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to get options failed with status code %d", response.StatusCode)
	}

	defer response.Body.Close()

	return html.Parse(response.Body)
}
