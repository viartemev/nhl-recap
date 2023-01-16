package util

import (
	"context"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	ctx := context.Background()

	want := []int{0, 2, 4, 6, 8}
	got := make([]int, 0)
	filtered := Filter(ctx, testGenerator(), evenNumbers)
	for i := range filtered {
		got = append(got, i)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %d is not equal to get %d", want, got)
	}
}

func evenNumbers(element int) bool {
	if element%2 == 0 {
		return true
	} else {
		return false
	}
}

func testGenerator() chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < 10; i++ {
			out <- i
		}
	}()
	return out
}
