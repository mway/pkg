// Copyright (c) 2020 Matt Way
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

package internal

import (
	"container/heap"
)

// Key represents a graph vertex.
type Key uint64

// Path represents an ordered series of contiguously incident vertices.
type Path struct {
	Cost     int
	Vertices []Key // TODO(mway): list
}

// Extend extends p to contain vertex as the latest incident vertex and
// increases p's total cost by cost.
func (p Path) Extend(cost int, vertex Key) Path {
	dup := p
	dup.Cost += cost

	dup.Vertices = make([]Key, len(p.Vertices), len(p.Vertices)+1)
	copy(dup.Vertices, p.Vertices)
	dup.Vertices = append(dup.Vertices, vertex)

	return dup
}

// Paths is a sortable collection of paths.
type Paths []Path

// Len returns the number of paths in p.
func (p Paths) Len() int {
	return len(p)
}

// Less returns true if the cost of the ith Path in p is less than the cost of
// the jth path.
func (p Paths) Less(i int, j int) bool {
	return p[i].Cost < p[j].Cost
}

// Swap swaps the ith and jth paths in p.
func (p Paths) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

// Push pushes path onto p. It should not be called directly and is used to
// satisfy heap.Interface.
func (p *Paths) Push(path interface{}) {
	*p = append(*p, path.(Path))
}

// Pop removes and returns the latest path from p.
func (p *Paths) Pop() interface{} {
	var (
		deref = *p
		n     = len(deref)
		edge  interface{}
	)

	if n > 0 {
		edge, *p = deref[n-1], deref[:n-1]
	}

	return edge
}

// PathHeap is a convenience wrapper around Paths and heap.Interface.
type PathHeap struct {
	paths *Paths
}

// NewPathHeap creates a new PathHeap, initialized to contain paths.
func NewPathHeap(paths ...Path) *PathHeap {
	pheap := &PathHeap{
		paths: &Paths{},
	}

	heap.Init(pheap.paths)

	for _, path := range paths {
		pheap.Push(path)
	}

	return pheap
}

// Len returns the heap's length.
func (p *PathHeap) Len() int {
	return p.paths.Len()
}

// Push pushes path onto the heap.
func (p *PathHeap) Push(path Path) {
	heap.Push(p.paths, path)
}

// Pop removes and returns the latest path in p.
func (p *PathHeap) Pop() Path {
	ipath := heap.Pop(p.paths)
	return ipath.(Path)
}
