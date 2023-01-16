package util

import (
	"context"
	"time"
)

type Fetcher[K any] interface {
	Fetch(ctx context.Context) chan K
}

type Subscription[K any] interface {
	Updates() <-chan K
}

type sub[K any] struct {
	fetcher Fetcher[K]
	updates chan K
}

func (s *sub[K]) Updates() <-chan K {
	return s.updates
}

func NewSubscription[K any](ctx context.Context, fetcher Fetcher[K], ticker *time.Ticker) Subscription[K] {
	s := &sub[K]{
		fetcher: fetcher,
		updates: make(chan K),
	}
	go s.serve(ctx, ticker)
	return s
}

func (s *sub[K]) serve(ctx context.Context, ticker *time.Ticker) {
	defer close(s.updates)
	for {
		select {
		case <-ticker.C:
			for res := range s.fetcher.Fetch(ctx) {
				select {
				case s.updates <- res:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
