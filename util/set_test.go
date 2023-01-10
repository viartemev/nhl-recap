package util

import "testing"

func TestSet_Add_Exists(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Add(key)
	if !set.Exists(key) {
		t.Errorf("key %d doesn't exist in set", key)
	}
}

func TestSet_Delete(t *testing.T) {
	set := NewSet[int]()
	key := 42
	set.Add(key)
	if !set.Exists(key) {
		t.Errorf("key %d doesn't exist in set", key)
	}
	set.Delete(key)
	if set.Exists(key) {
		t.Errorf("key %d exists in set, but shouldn't", key)
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
		t.Errorf("Set should contains only one element and counter should be 1")
	}
}
