package queue

import (
	"sync"
)

// A Queue is a FIFO queue.
type Queue struct {
	items *item
	last  *item
	mtx   sync.Mutex
	size  int
}

// Len returns the number of items in the queue.
func (s *Queue) Len() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.size
}

// Front returns the item at the front of the queue.
func (s *Queue) Front() interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.items == nil {
		return nil
	}

	return s.items.value
}

// Back returns the item at the back of the queue.
func (s *Queue) Back() interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.last == nil {
		return nil
	}

	return s.last.value
}

// Dequeue removes the item at the front of the queue and returns it.
func (s *Queue) Dequeue() interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.items == nil {
		return nil
	}

	s.size--
	val := s.items.value
	s.items = s.items.next

	return val
}

// Push pushes value onto the back of the queue.
func (s *Queue) Push(value interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.size++

	ob := &item{
		value: value,
	}

	if s.items == nil {
		s.items = ob
	}

	if s.last != nil {
		s.last.next = ob
	}

	s.last = ob
}

type item struct {
	value interface{}
	next  *item
}
