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
	}

	return keys
}

func NewSet[T comparable]() *Set[T] {
	s := &Set[T]{
		list: map[T]struct{}{},
	}
	return s
}
