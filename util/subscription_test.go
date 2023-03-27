package util

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockFetcher[K any] struct {
	results chan K
	err     error
}

func (m *mockFetcher[K]) Fetch(ctx context.Context) (chan K, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.results, nil
}

func TestNewSubscription(t *testing.T) {
	ctx := context.Background()
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	fetcher := &mockFetcher[int]{
		results: make(chan int),
	}
	s := NewSubscription[int](ctx, fetcher, ticker)
	if s == nil {
		t.Errorf("NewSubscription() returned nil")
	}
}

func TestSubscription_Updates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fetcher := &mockFetcher[int]{
		results: make(chan int),
	}
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	s := NewSubscription[int](ctx, fetcher, ticker)

	go func() {
		for i := 0; i < 5; i++ {
			fetcher.results <- i
		}
		close(fetcher.results)
	}()

	var count int
	for i := range s.Updates() {
		if i != count {
			t.Errorf("expected %v, but got %v", count, i)
		}
		count++
		if count == 5 {
			break
		}
	}

	if count != 5 {
		t.Errorf("expected 5 updates, but got %v", count)
	}
}

func TestSubscription_UpdatesWithContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	fetcher := &mockFetcher[int]{
		results: make(chan int),
	}
	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()

	s := NewSubscription[int](ctx, fetcher, ticker)

	cancel()

	select {
	case _, ok := <-s.Updates():
		if ok {
			t.Errorf("received unexpected update")
		}
	case <-time.After(time.Second * 3):
	}
}

func TestSubscription_UpdatesWithFetcherError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fetcher := &mockFetcher[int]{
		err: errors.New("fetcher error"),
	}
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	s := NewSubscription[int](ctx, fetcher, ticker)

	select {
	case <-s.Updates():
		t.Errorf("received unexpected update")
	case <-time.After(time.Second):
	}
}
