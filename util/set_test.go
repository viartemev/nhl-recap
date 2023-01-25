package util

import (
	"sync"
	"testing"
)

func TestSet_Add_If_Element_Not_Exists(t *testing.T) {
	set := NewSet[int]()
	key := 42
	if !set.Add(key) {
		t.Errorf("key %d should be added", key)
	}
	if !set.Exists(key) {
		t.Errorf("key %d doesn't exist in set", key)
	}
}

func TestSet_Add_If_Element_Exists(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Add(key)
	if set.Add(key) {
		t.Errorf("key %d shouldn't be added", key)
	}
	if !set.Exists(key) {
		t.Errorf("key %d doesn't exist in set", key)
	}
}

func TestSet_Delete_If_Element_Exists(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Add(key)
	if !set.Exists(key) {
		t.Errorf("key %d doesn't exist in set", key)
	}
	if !set.Delete(key) {
		t.Errorf("key %d should be deleted", key)
	}
	if set.Exists(key) {
		t.Errorf("key %d exists in set, but shouldn't", key)
	}
}

func TestSet_Delete_If_Element_Not_Exists(t *testing.T) {
	set := NewSet[int]()
	key := 42
	if set.Delete(key) {
		t.Errorf("key %d should be deleted", key)
	}
}

func TestSet_Delete_Not_Existing_Element(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Delete(key)
	if set.Exists(key) {
		t.Errorf("key %d exists in set, but shouldn't", key)
	}
}

func TestSet_Range(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Add(key)

	counter := 0
	set.Range(func(value int) {
		counter++
	})

	if counter != 1 {
		t.Errorf("Set should contains only one element And counter should be 1")
	}
}

func TestSet_Concurrent(t *testing.T) {
	set := NewSet[int]()
	key := 42

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		set.Add(key)
	}()

	go func() {
		defer wg.Done()
		set.Delete(key)
	}()

	go func() {
		defer wg.Done()
		set.Exists(key)
	}()

	wg.Wait()
}
