package util

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type Fetcher[K any] interface {
	Fetch(ctx context.Context) (chan K, error)
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
			result, err := s.fetcher.Fetch(ctx)
			if err != nil {
				log.WithError(err).Errorf("Got an error from fetch method")
				break
			}
			for res := range result {
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
