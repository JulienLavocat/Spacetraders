package utils

import "container/list"

type Queue[T any] struct {
	queue *list.List
}

// TODO: Could use https://go.dev/wiki/SliceTricks rather than list.List
func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		queue: list.New(),
	}
}

func (s *Queue[T]) Clear() {
	s.queue.Init()
}

func (s *Queue[T]) Enqueue(v T) {
	s.queue.PushBack(v)
}

func (s *Queue[T]) HasNext() bool {
	return s.queue.Len() > 0
}

func (s *Queue[T]) Dequeue() T {
	element := s.queue.Front()
	s.queue.Remove(element)
	return element.Value.(T)
}

func (s *Queue[T]) QueueAll(values []T) {
	for i := range values {
		s.Enqueue(values[i])
	}
}
