package graph_test

import (
	"testing"

	"github.com/mway/pkg/x/container/graph"
	"github.com/stretchr/testify/require"
)

type testStringer struct{}

func (t testStringer) String() string { return "test" }

func TestVertex(t *testing.T) {
	var v *graph.Vertex
	require.Equal(t, graph.Key{}, v.Key())
	require.Equal(t, "", v.String())
	require.Equal(t, nil, v.Value())

	var (
		g      = graph.New()
		value  = testStringer{}
		vertex graph.Vertex
	)

	g.AddVertex(value)
	g.VisitVertices(graph.Root, func(v graph.Vertex) bool {
		vertex = v
		return false
	})

	require.Equal(t, value.String(), vertex.String())
	require.Equal(t, value, vertex.Value())
}
