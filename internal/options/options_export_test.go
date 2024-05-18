package options

import (
	"context"
)

func (o Opt) ShouldFetch(ctx context.Context) (bool, error) {
	return o.shouldFetch(ctx)
}
