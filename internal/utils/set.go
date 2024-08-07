package utils

// Source (adapted to be able to use generics): https://gist.github.com/bgadrian/cb8b9344d9c66571ef331a14eb7a2e80

type Set[T comparable] struct {
	list map[T]struct{}
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set[T]) Add(v T) {
	s.list[v] = struct{}{}
}

func (s *Set[T]) AddValues(values []T) {
	for i := range values {
		s.Add(values[i])
	}
}

func (s *Set[T]) Remove(v T) {
	delete(s.list, v)
}

func (s *Set[T]) Clear() {
	s.list = make(map[T]struct{})
}

func (s *Set[T]) Size() int {
	return len(s.list)
}

func (s *Set[T]) Values() []T {
	keys := make([]T, len(s.list))

	i := 0
	for k := range s.list {
		keys[i] = k
		i++
	}

	return keys
}

func NewSet[T comparable]() *Set[T] {
	s := &Set[T]{
		list: map[T]struct{}{},
	}
	return s
}

func NewSetFrom[T comparable](values []T) *Set[T] {
	set := NewSet[T]()
	set.AddValues(values)
	return set
}

func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		res.Add(v)
	}

	for v := range s2.list {
		res.Add(v)
	}
	return res
}

func (s *Set[T]) Intersect(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		if !s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}

// Difference returns the subset from s, that doesn't exists in s2 (param)
func (s *Set[T]) Difference(s2 *Set[T]) *Set[T] {
	res := NewSet[T]()
	for v := range s.list {
		if s2.Has(v) {
			continue
		}
		res.Add(v)
	}
	return res
}
