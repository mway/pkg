package graph

import (
	"github.com/mway/pkg/x/container/graph/internal"
)

// Dijkstra evaluates paths in g spanning from and to and returns the cheapest
// path possible within the graph.
func Dijkstra(g *Graph, from Key, to Key) Path {
	var (
		internalFrom = internal.Key(from.key)
		internalTo   = internal.Key(to.key)
		visited      = make(map[internal.Key]struct{})
		heap         = internal.NewPathHeap(internal.Path{
			Cost:     0,
			Vertices: []internal.Key{internalFrom},
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

		if key == internalTo {
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
