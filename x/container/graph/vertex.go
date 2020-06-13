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

package graph

import (
	"fmt"

	"github.com/mway/pkg/x/container/graph/internal"
)

// Vertex is a point node in a graph.
type Vertex struct {
	key   internal.Key
	graph *Graph
	value interface{}
}

// Key provides the key for v.
func (v *Vertex) Key() Key {
	if v == nil {
		return _zeroKey
	}

	return newKey(v.key)
}

// String returns a string representation of v.
func (v *Vertex) String() string {
	if v == nil {
		return ""
	}

	if str, ok := v.value.(fmt.Stringer); ok {
		return str.String()
	}

	return fmt.Sprintf("%v", v.value)
}

// Value returns the value held in v.
func (v *Vertex) Value() interface{} {
	if v == nil {
		return nil
	}

	return v.value
}
