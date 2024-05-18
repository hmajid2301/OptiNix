package options

import (
	"context"
)

func ShouldFetch(ctx context.Context, o Opt) (bool, error) {
	return o.shouldFetch(ctx)
}
