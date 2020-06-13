package stack

import (
	"sync"
)

// A Stack is a LIFO queue.
type Stack struct {
	items *item
	mtx   sync.Mutex
	size  int
}

// Len returns the current number of items in the stack.
func (s *Stack) Len() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.size
}

// Top returns the topmost item in the stack. The item is not removed.
func (s *Stack) Top() interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.items == nil {
		return nil
	}

	return s.items.value
}

// Pop removes the topmost item from the stack and returns it.
func (s *Stack) Pop() interface{} {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.items == nil {
		return nil
	}

	s.size--
	val := s.items.value
	s.items = s.items.prev

	return val
}

// Push pushes value onto the top of the stack.
func (s *Stack) Push(value interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.size++
	s.items = &item{
		value: value,
		prev:  s.items,
	}
}

type item struct {
	value interface{}
	prev  *item
}
