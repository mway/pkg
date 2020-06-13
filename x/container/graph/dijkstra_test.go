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
	"bytes"
	"testing"

	"github.com/mway/pkg/x/container/graph"
	"github.com/stretchr/testify/require"
)

func TestDijkstraPath(t *testing.T) {
	var (
		g    = graph.New()
		keyA = g.AddVertex('A')
		keyB = g.AddVertex('B')
		keyC = g.AddVertex('C')
		keyD = g.AddVertex('D')
		keyE = g.AddVertex('E')
		keyF = g.AddVertex('F')
		keyG = g.AddVertex('G')
		keyH = g.AddVertex('H')
		keyI = g.AddVertex('I')
		keyJ = g.AddVertex('J')
		keyK = g.AddVertex('K')
		keyL = g.AddVertex('L')
	)

	g.AddEdgeCost(keyA, keyB, 1)

	// Short but expensive
	g.AddEdgeCost(keyB, keyC, 10)
	g.AddEdgeCost(keyC, keyD, 10)

	// Short but medium
	g.AddEdgeCost(keyB, keyK, 4)
	g.AddEdgeCost(keyK, keyL, 4)
	g.AddEdgeCost(keyL, keyD, 4)

	// Long but cheap
	g.AddEdgeCost(keyB, keyE, 1)
	g.AddEdgeCost(keyE, keyF, 1)
	g.AddEdgeCost(keyF, keyG, 1)
	g.AddEdgeCost(keyG, keyH, 1)
	g.AddEdgeCost(keyH, keyI, 1)
	g.AddEdgeCost(keyI, keyJ, 1)
	g.AddEdgeCost(keyJ, keyD, 1)

	expected := graph.Path{
		Cost: 8,
		Vertices: []graph.Key{
			keyA, keyB, keyE, keyF,
			keyG, keyH, keyI, keyJ,
			keyD,
		},
	}

	actual := g.FindPath(graph.Dijkstra, keyA, keyD)

	require.Equal(t, expected.Cost, actual.Cost)
	require.Equal(t, expected.Vertices, actual.Vertices)

	var buf bytes.Buffer
	for i, vertex := range actual.Vertices {
		if i > 0 {
			buf.WriteString(" -> ")
		}

		node, ok := g.Get(vertex)
		require.True(t, ok)

		buf.WriteRune(node.Value().(rune))
	}

	// fmt.Println("Graph traversal path:")
	// fmt.Println(buf.String())
}
