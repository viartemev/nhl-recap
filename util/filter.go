package util

import (
	"context"
)

type predicate[K any] func(element K) bool

// And returns predicate. Predicate returns true if all passed predicates returned true
func And[K any](predicates ...predicate[K]) predicate[K] {
	return func(element K) bool {
		var result = true
		for _, p := range predicates {
			result = result && p(element)
		}
		return result
	}
}

func Filter[K any](ctx context.Context, in chan K, predicate predicate[K]) chan K {
	out := make(chan K)

	go func() {
		defer close(out)
		for {
			select {
			case element, ok := <-in:
				if !ok {
					return
				}
				if !predicate(element) {
					continue
				}
				select {
				case out <- element:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
