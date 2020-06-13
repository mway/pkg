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
	"github.com/mway/pkg/x/container/graph/internal"
)

// Dijkstra evaluates paths in g spanning from and to and returns the cheapest
// path possible within the graph.
func Dijkstra(g *Graph, from Key, to Key) Path {
	var (
		visited = make(map[internal.Key]struct{})
		heap    = internal.NewPathHeap(internal.Path{
			Cost:     0,
			Vertices: []internal.Key{from.key},
		})
	)

	for heap.Len() > 0 {
		var (
			path = heap.Pop()
			key  = path.Vertices[len(path.Vertices)-1]
		)

		if _, seen := visited[key]; seen {
			continue
		}

		if key == to.key {
			return newPathFromInternal(path)
		}

		g.VisitEdges(newKey(key), func(edge Edge) bool {
			if _, seen := visited[edge.End.key]; !seen {
				heap.Push(path.Extend(edge.Cost, edge.End.key))
			}

			return true
		})
	}

	return Path{}
}
