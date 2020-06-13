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

package graph_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/mway/pkg/x/container/graph"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	cases := []struct {
		graph    *graph.Graph
		expected []interface{}
	}{
		{
			graph: func() *graph.Graph {
				g := graph.New()

				key1 := g.AddVertex(1)
				key2 := g.AddVertex(2)

				require.Equal(t, 2, g.Order())
				require.True(t, g.AddEdge(key1, key2))
				require.False(t, g.AddEdge(graph.Key{}, key2))
				require.False(t, g.AddEdge(key1, graph.Key{}))

				return g
			}(),
			expected: []interface{}{1, 2},
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var actual []interface{}

			tc.graph.VisitVertices(graph.Root, func(node graph.Vertex) bool {
				actual = append(actual, node.Value())
				return true
			})
			require.ElementsMatch(t, tc.expected, actual)
		})
	}
}

func TestGraphSubgraph(t *testing.T) {
	g := graph.New()

	k1 := g.AddVertex(1)
	k2 := g.AddVertex(2)
	g.AddEdge(k1, k2)

	k3 := g.AddVertex(3)
	k4 := g.AddVertex(4)
	g.AddEdge(k2, k3)
	g.AddEdge(k3, k4)

	var actual []interface{}
	g.VisitVertices(k3, func(node graph.Vertex) bool {
		actual = append(actual, node.Value())
		return true
	})

	require.ElementsMatch(t, []interface{}{3, 4}, actual)
}

func TestGraphDeleteVertex(t *testing.T) {
	var (
		g        = graph.New()
		vertices = []graph.Key{
			g.AddVertex(1),
			g.AddVertex(2),
			g.AddVertex(3),
		}
	)

	numVertices := func() (n int) {
		g.VisitVertices(graph.Root, func(graph.Vertex) bool {
			n++
			return true
		})
		return
	}

	for i := 0; i < len(vertices); i++ {
		require.Equal(t, len(vertices[i:]), numVertices())
		g.DeleteVertex(vertices[i])
	}

	require.Equal(t, 0, numVertices())
}

func TestGraphDeleteEdge(t *testing.T) {
	var (
		g        = graph.New()
		vertices = []graph.Key{
			g.AddVertex(1),
			g.AddVertex(2),
			g.AddVertex(3),
		}
	)

	for i := 0; i < len(vertices)-1; i++ {
		require.True(t, g.AddEdge(vertices[i], vertices[i+1]))
	}

	numEdges := func() (n int) {
		g.VisitEdges(graph.Root, func(graph.Edge) bool {
			n++
			return true
		})
		return
	}

	for i := 0; i < len(vertices)-1; i++ {
		require.Equal(t, len(vertices[i:])-1, numEdges())
		g.DeleteEdge(vertices[i], vertices[i+1])
	}

	require.Equal(t, 0, numEdges())
}

func TestGraphFilterVertices(t *testing.T) {
	var (
		g        = graph.New()
		k1       = g.AddVertex(1)
		k2       = g.AddVertex(2)
		k3       = g.AddVertex(3)
		k4       = g.AddVertex(4)
		k5       = g.AddVertex(5)
		vertices int
		edges    int
	)

	g.AddEdge(k1, k2)
	g.AddEdge(k2, k3)
	g.AddEdge(k3, k4)
	g.AddEdge(k4, k5)

	count := func() {
		g.VisitVertices(graph.Root, func(vertex graph.Vertex) bool {
			vertices++
			return true
		})

		g.VisitEdges(graph.Root, func(edge graph.Edge) bool {
			edges++
			return true
		})
	}

	vertices, edges = 0, 0
	count()

	require.Equal(t, 5, vertices)
	require.Equal(t, 4, edges)

	g = g.FilterVertices(func(vertex graph.Vertex) bool {
		if k := vertex.Key(); k == k1 || k == k2 {
			return true
		}

		return false
	})

	vertices, edges = 0, 0
	count()

	require.Equal(t, 2, vertices)
	require.Equal(t, 1, edges)
}

func TestGraphFilterVerticesMiddle(t *testing.T) {
	var (
		g        = graph.New()
		k1       = g.AddVertex(1)
		k2       = g.AddVertex(2)
		k3       = g.AddVertex(3)
		k4       = g.AddVertex(4)
		k5       = g.AddVertex(5)
		vertices int
		edges    int
	)

	g.AddEdge(k1, k2)
	g.AddEdge(k2, k3)
	g.AddEdge(k3, k4)
	g.AddEdge(k4, k5)

	count := func() {
		g.VisitVertices(graph.Root, func(vertex graph.Vertex) bool {
			vertices++
			return true
		})

		g.VisitEdges(graph.Root, func(edge graph.Edge) bool {
			edges++
			return true
		})
	}

	vertices, edges = 0, 0
	count()

	require.Equal(t, 5, vertices)
	require.Equal(t, 4, edges)

	g = g.FilterVertices(func(vertex graph.Vertex) bool {
		return vertex.Key() != k3
	})

	vertices, edges = 0, 0
	count()

	require.Equal(t, 4, vertices)
	require.Equal(t, 2, edges)
}

func TestGraphFilterEdges(t *testing.T) {
	var (
		g        = graph.New()
		k1       = g.AddVertex(1)
		k2       = g.AddVertex(2)
		k3       = g.AddVertex(3)
		k4       = g.AddVertex(4)
		k5       = g.AddVertex(5)
		vertices int
		edges    int
	)

	g.AddEdge(k1, k2)
	g.AddEdge(k2, k3)
	g.AddEdge(k3, k4)
	g.AddEdge(k4, k5)

	count := func() {
		g.VisitVertices(graph.Root, func(vertex graph.Vertex) bool {
			vertices++
			return true
		})

		g.VisitEdges(graph.Root, func(edge graph.Edge) bool {
			edges++
			return true
		})
	}

	vertices, edges = 0, 0
	count()

	require.Equal(t, 5, vertices)
	require.Equal(t, 4, edges)

	g = g.FilterEdges(func(edge graph.Edge) bool {
		return edge.Start.Key() == k1
	})

	vertices, edges = 0, 0
	count()

	require.Equal(t, 5, vertices)
	require.Equal(t, 1, edges)
}

func TestGraphEdges(t *testing.T) {
	g := graph.New()

	k1 := g.AddVertex(1)
	k2 := g.AddVertex(2)
	g.AddEdge(k1, k2)

	k3 := g.AddVertex(3)
	k4 := g.AddVertex(4)
	g.AddEdge(k2, k3)
	g.AddEdge(k3, k4)
	g.AddEdgeCost(k1, k4, 2)

	get := func(key graph.Key) graph.Vertex {
		node, ok := g.Get(key)
		if !ok {
			panic(fmt.Sprintf("could not get node with key %v", key))
		}

		return node
	}

	t.Run("root", func(t *testing.T) {
		var (
			expected = []graph.Edge{
				{Start: get(k1), End: get(k2), Cost: 1},
				{Start: get(k2), End: get(k3), Cost: 1},
				{Start: get(k3), End: get(k4), Cost: 1},
				{Start: get(k1), End: get(k4), Cost: 2},
			}
			actual []graph.Edge
		)

		g.VisitEdges(graph.Root, func(edge graph.Edge) bool {
			actual = append(actual, edge)
			return true
		})

		require.ElementsMatch(t, expected, actual)
	})

	t.Run("subgraph", func(t *testing.T) {
		var (
			expected = []graph.Edge{
				{Start: get(k3), End: get(k4), Cost: 1},
			}
			actual []graph.Edge
		)

		g.VisitEdges(k3, func(edge graph.Edge) bool {
			actual = append(actual, edge)
			return true
		})

		require.ElementsMatch(t, expected, actual)
	})
}

func TestGraphString(t *testing.T) {
	g := graph.New()

	k1 := g.AddVertex(1)
	k2 := g.AddVertex(2)
	g.AddEdge(k1, k2)

	parts := strings.Split(strings.TrimSpace(g.String()), "\n")
	require.Equal(t, 2, len(parts))
	require.Contains(t, parts, "2")
	require.Contains(t, parts, "1\t2")
}
