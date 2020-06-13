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
