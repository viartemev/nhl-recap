package util

import (
	"context"
	"testing"
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
