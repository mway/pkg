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
