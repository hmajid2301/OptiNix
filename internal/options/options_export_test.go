package options

import (
	"context"
)

func (s Searcher) ShouldFetch(ctx context.Context) (bool, error) {
	return s.shouldFetch(ctx)
}
