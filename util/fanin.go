package util

import (
	"context"
	"sync"
)

func FanIn[K any, V any](ctx context.Context, arr []K, fn func(games K) *V) chan *V {
	out := make(chan *V)
	var wg sync.WaitGroup

	for _, game := range arr {
		wg.Add(1)
		go func(games K) {
			defer wg.Done()
			//TODO rewrite to channels or it's too much?
			select {
			case out <- fn(games):
			case <-ctx.Done():
				return
			}
		}(game)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}