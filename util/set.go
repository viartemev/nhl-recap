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

func (s *Set[K]) Add(key K) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.values[key]
	if !ok {
		s.values[key] = true
		return true
	}
	return false
}

func (s *Set[K]) Exists(key K) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.values[key]
	return ok
}

func (s *Set[K]) Delete(key K) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.values[key]
	if ok {
		delete(s.values, key)
		return true
	}
	return false
}

func (s *Set[K]) Range(fun func(value K)) {
	s.RLock()
	defer s.RUnlock()
	for key, _ := range s.values {
		fun(key)
	}
}
