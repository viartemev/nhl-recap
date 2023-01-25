package util

import (
	"context"
	"testing"
)

func BenchmarkFanIn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FanIn[int, int](context.Background(), []int{1, 2, 3}, func(element int) int {
			return element * 3
		})
	}
}
