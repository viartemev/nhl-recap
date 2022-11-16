package util

type Set[K comparable] struct {
	values map[K]bool
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{
		values: make(map[K]bool),
	}
}

func (s *Set[K]) Add(key K) {
	s.values[key] = true
}

func (s *Set[K]) Delete(key K) {
	delete(s.values, key)
}

// TODO make it concurrent
func (s *Set[K]) Range(fun func(value K)) {
	for key, _ := range s.values {
		fun(key)
	}
}
