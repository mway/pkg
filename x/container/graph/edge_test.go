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
	"strings"
	"testing"

	"github.com/mway/pkg/x/container/graph"
	"github.com/stretchr/testify/require"
)

func TestEdge(t *testing.T) {
	var (
		g        = graph.New()
		key1     = g.AddVertex(1)
		key2     = g.AddVertex(2)
		vertices = make(map[graph.Key]graph.Vertex)
	)

	g.VisitVertices(graph.Root, func(v graph.Vertex) bool {
		vertices[v.Key()] = v
		return true
	})

	require.Contains(t, vertices, key1)
	require.Contains(t, vertices, key2)

	var edge graph.Edge
	g.AddEdge(key1, key2)
	g.VisitEdges(graph.Root, func(e graph.Edge) bool {
		edge = e
		return true
	})

	formatted := edge.String()
	require.True(t, strings.Contains(formatted, key1.String()))
	require.True(t, strings.Contains(formatted, key2.String()))
}
