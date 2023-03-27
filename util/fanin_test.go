package util

import (
	"context"
	"testing"
	"time"
)

func BenchmarkFanIn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetchers := func() <-chan int {
			out := make(chan int)
			go func() {
				out <- 42
			}()
			return out
		}
		FanIn[int](context.Background(), fetchers())
	}
}

func TestFanIn_WithTwoChannelsWithOneItemEach(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
	}()

	go func() {
		ch2 <- 2
	}()

	out := FanIn(ctx, ch1, ch2)

	var result []int
	for i := 0; i < 2; i++ {
		select {
		case val := <-out:
			result = append(result, val)
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}

	// Check that each expected value is present in the result
	expectedValues := map[int]bool{1: true, 2: true}
	for _, val := range result {
		if _, ok := expectedValues[val]; !ok {
			t.Fatalf("unexpected value: %d", val)
		}
		delete(expectedValues, val)
	}

	// Check that all expected values were found
	if len(expectedValues) != 0 {
		t.Fatalf("missing expected values: %v", expectedValues)
	}
}

func TestFanIn_WithTwoChannelsWithDifferentNumbersOfItems(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
		ch1 <- 2
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- 3
	}()

	out := FanIn(ctx, ch1, ch2)

	var result []int
	for i := 0; i < 3; i++ {
		select {
		case val := <-out:
			result = append(result, val)
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}

	if result[0] != 1 || result[1] != 2 || result[2] != 3 {
		t.Fatalf("expected [1, 2, 3], got %v", result)
	}
}

func TestFanIn_WithNoChannels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out := FanIn[int](ctx)

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("unexpected value")
		}
	case <-time.After(time.Second):
		// expected timeout
	}
}
