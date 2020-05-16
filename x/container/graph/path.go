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
