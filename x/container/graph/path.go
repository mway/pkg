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

package graph

import (
	"github.com/mway/pkg/x/container/graph/internal"
)

// A Path is a costed sequence of vertices, where cost is equal to the sum of
// edge costs used to construct the path.
type (
	Path struct {
		Cost     int
		Vertices []Key
	}

	// Paths is an ordered collection of paths.
	Paths = []Path
)

func newPathFromInternal(path internal.Path) Path {
	if len(path.Vertices) == 0 {
		return Path{}
	}

	p := Path{
		Cost:     path.Cost,
		Vertices: make([]Key, len(path.Vertices)),
	}

	for i, key := range path.Vertices {
		p.Vertices[i] = newKey(key)
	}

	return p
}
