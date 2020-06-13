// Copyright (c) 2019 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
