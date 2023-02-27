package util

import (
	"context"
	"sync"
)

func FanIn[K any](ctx context.Context, fetchers ...<-chan K) chan K {
	out := make(chan K)
	var wg sync.WaitGroup

	for _, f := range fetchers {
		wg.Add(1)
		go func(fetcher <-chan K) {
			defer wg.Done()
			for {
				select {
				case out <- <-fetcher:
				case <-ctx.Done():
					return
				}
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
