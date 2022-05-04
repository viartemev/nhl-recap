package util

import (
	"sync"
)

type ConcurrentMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		m: make(map[K]V),
	}
}

func (c *ConcurrentMap[K, V]) Put(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.m[key] = value
}

func (c *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *ConcurrentMap[K, V]) Length() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.m)
}

func (c *ConcurrentMap[K, V]) Range(fun func(value V)) {
	for _, info := range c.m {
		fun(info)
	}
}
