package util

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestFilter_Even_Numbers(t *testing.T) {
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

func TestFilter_Empty_In_Channel_Cancelation(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	got := make([]int, 0)
	filtered := Filter(ctx, make(chan int), evenNumbers)
	go func() {
		time.Sleep(300 * time.Millisecond)
		cancelFunc()
	}()
	for i := range filtered {
		got = append(got, i)
	}
	if len(got) != 0 {
		t.Errorf("Filtered slice should be empty")
	}
	_, ok := <-filtered
	if ok {
		t.Errorf("Channel should be closed")
	}
}

func TestFilter_Empty_Out_Channel_Cancelation(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	got := make([]int, 0)
	filtered := Filter(ctx, testGenerator(), evenNumbers)
	for i := 0; i < 1; i++ {
		value := <-filtered
		got = append(got, value)
	}
	cancelFunc()
	if len(got) != 1 {
		t.Errorf("Filtered slice should be empty")
	}
	_, ok := <-filtered
	if ok {
		t.Errorf("Channel should be closed")
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
