package util

import (
	"sync"
)

type Set[K comparable] struct {
	values map[K]bool
	sync.RWMutex
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{values: make(map[K]bool)}
}

func (s *Set[K]) Add(key K) {
	s.Lock()
	defer s.Unlock()
	s.values[key] = true
}

func (s *Set[K]) Exists(key K) bool {
	s.RLock()
	defer s.RUnlock()
	return s.values[key]
}

func (s *Set[K]) Delete(key K) {
	s.Lock()
	defer s.Unlock()
	delete(s.values, key)
}

func (s *Set[K]) Range(fun func(value K)) {
	s.RLock()
	defer s.RUnlock()
	for key, _ := range s.values {
		fun(key)
	}
}
