package graph

import "fmt"

// Edge is a connection between two incident vertices in a graph. An edge is
// always directed, but for undirected graphs, can be assumed to be invertable
// with the same cost.
type Edge struct {
	Start Vertex
	End   Vertex
	Cost  int
}

// String returns a string representation of e.
func (e *Edge) String() string {
	return fmt.Sprintf(
		"%s->%s(%d)",
		e.Start.Key().String(),
		e.End.Key().String(),
		e.Cost,
	)
}
