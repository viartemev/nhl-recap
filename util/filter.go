package util

import (
	"context"
)

func Filter[K any](ctx context.Context, in chan K, predicate func(element K) bool) chan K {
	out := make(chan K)

	go func() {
		defer close(out)
		for element := range in {
			if !predicate(element) {
				continue
			}
			select {
			case out <- element:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
