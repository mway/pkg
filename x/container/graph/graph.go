package graph

import (
	"bytes"
	"math"
	"sync"

	"github.com/mway/pkg/x/container/graph/internal"
)

// TODO(mway): proper A*

var (
	// Any TODOC
	Any = Key{key: internal.Key(math.MaxUint64)}
	// Root is a default key for specifying graph iteration and/or filtering.
	Root = Key{key: internal.Key(math.MaxUint64 - 1)}

	_zeroKey    Key
	_zeroVertex Vertex
)

type (
	// VertexFilterFunc is used by Graph as a predicate function to filter its vertices.
	VertexFilterFunc = func(Vertex) bool
	// VertexVisitorFunc is used by Graph to yield visited vertices to callers.
	VertexVisitorFunc = func(Vertex) bool
	// EdgeFilterFunc is used by Graph as a predicate function to filter its edges.
	EdgeFilterFunc = func(Edge) bool
	// EdgeVisitorFunc is used by Graph to yield visited edges to callers.
	EdgeVisitorFunc = func(Edge) bool
	// FindPathFunc is used by Graph do perform pluggable pathing/costing for a single path.
	FindPathFunc = func(graph *Graph, from Key, to Key) Path
	// FindPathsFunc is used by Graph to perform pluggable pathing/costing for multiple paths.
	FindPathsFunc = func(graph *Graph, from Key, to Key) Paths
)

// A Graph is a basic data structure defined as a set of vertices and a set of edges.
// Graphs may be either directed or undirected.
type Graph struct {
	lastKey  uint64
	mtx      sync.Mutex
	vertices map[internal.Key]Vertex
	edges    map[internal.Key]map[internal.Key]int
	redges   map[internal.Key]map[internal.Key]int

	config struct {
		undirected bool
	}
}

// New constructs a new directed graph.
func New() *Graph {
	return &Graph{
		vertices: make(map[internal.Key]Vertex),
		edges:    make(map[internal.Key]map[internal.Key]int),
		redges:   make(map[internal.Key]map[internal.Key]int),
	}
}

// NewUndirected constructs a new undirected graph.
func NewUndirected() *Graph {
	g := New()
	g.config.undirected = true

	return g
}

// AddEdge adds a new directed edge to the graph, spanning the vertices from and
// to. If the graph is undirected, a second edge is added implicitly in the
// reverse direction. Edges added with AddEdge have an implicit cost of 1.
func (g *Graph) AddEdge(from Key, to Key) bool {
	return g.AddEdgeCost(from, to, 1)
}

// AddEdgeCost behaves identically to AddEdge except using the provided cost.
func (g *Graph) AddEdgeCost(from Key, to Key, cost int) bool {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	if _, exists := g.vertices[from.key]; !exists {
		return false
	}

	if _, exists := g.vertices[to.key]; !exists {
		return false
	}

	edges := g.getEdgesUnsafe(from.key)
	edges[to.key] = cost

	edges = g.getReverseEdgesUnsafe(to.key)
	edges[from.key] = cost

	return true
}

// AddVertex adds a new vertex containing value to the graph and returns its
// corresponding Key for subsequent lookup.
func (g *Graph) AddVertex(value interface{}) Key {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.lastKey++

	var (
		key  = internal.Key(g.lastKey)
		node = Vertex{
			key:   key,
			graph: g,
			value: value,
		}
	)

	g.vertices[key] = node

	return newKey(key)
}

// DeleteVertex deletes the vertex represented by key, if it exists.
//
// Importantly, deleting a vertex will also delete adjacent edges, potentially
// fragmenting the underlying graph or isolating vertices.
func (g *Graph) DeleteVertex(key Key) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.deleteEdgeUnsafe(key, Any)
	g.deleteEdgeUnsafe(Any, key)

	_, ok := g.vertices[key.key]
	if !ok {
		return
	}

	delete(g.vertices, key.key)
}

// DeleteEdge deletes the edge spanning vertices from and to, if such an edge
// exists. If the graph is undirected, the reverse edge will also be deleted.
func (g *Graph) DeleteEdge(from Key, to Key) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.deleteEdgeUnsafe(from, to)
}

// Get gets the Vertex represented by key, if it exists.
func (g *Graph) Get(key Key) (Vertex, bool) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	node, ok := g.vertices[key.key]
	return node, ok
}

// FilterVertices returns a copy of g, filtering the vertices of g based on the
// provided predicate filter: for a given vertex v, if filter(v) returns true, v
// remains part of g; if filter(v) returns false, however, v - as well as any
// edges referencing v - are removed from g.
//
// It is possible for filters to create disjoint (multi-part) sub-graphs, or to
// introduce vertex isolation.
func (g *Graph) FilterVertices(filter VertexFilterFunc) *Graph {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	dup := g.cloneUnsafe()
	for key, vertex := range dup.vertices {
		if filter(vertex) {
			continue
		}

		dup.DeleteVertex(newKey(key))
	}

	return dup
}

// FilterEdges returns a copy of g, filtering the edges of g based on the
// provided predicate filter: for a given edge e, if filter(e) returns true, e
// remains part of g; if filter(e) returns false, however, e is removed from g.
//
// A critical difference between FilterVertices and FilterEdges is that the
// former will prune edges (there is no such thing as non-incident/adjacent
// edges) whereas the latter will remove only edges (potentially introducing
// disjoint/multi-part sub-graphs or vertex isolation).
func (g *Graph) FilterEdges(filter EdgeFilterFunc) *Graph {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	dup := g.cloneUnsafe()
	for start, edges := range g.edges {
		for end, cost := range edges {
			keep := filter(Edge{
				Start: g.vertices[start],
				End:   g.vertices[end],
				Cost:  cost,
			})
			if keep {
				continue
			}

			dup.DeleteEdge(newKey(start), newKey(end))
		}
	}

	return dup
}

// FindPath uses algo to find a path spanning vertices from and to.
func (g *Graph) FindPath(algo FindPathFunc, from Key, to Key) Path {
	return algo(g, from, to)
}

// FindPaths uses algo to find paths spanning vertices from and to. The number
// and order of paths is determined by algo.
func (g *Graph) FindPaths(algo FindPathsFunc, from Key, to Key) Paths {
	return algo(g, from, to)
}

// Order returns the order of the graph (the number of vertices).
func (g *Graph) Order() int {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	return len(g.vertices)
}

// String provides a string representation of the graph.
func (g *Graph) String() string {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	var buf bytes.Buffer

	for key, node := range g.vertices {
		buf.WriteString(node.String())

		newline := false
		for dst := range g.edges[key] {
			vertex := g.vertices[dst]
			buf.WriteRune('\t')
			buf.WriteString(vertex.String())
			buf.WriteRune('\n')
			newline = true
		}

		if !newline {
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}

// VisitEdges uses fn to visit each edge, starting at the vertex key.
func (g *Graph) VisitEdges(key Key, fn EdgeVisitorFunc) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	visit := func(start internal.Key) {
		for end, cost := range g.edges[start] {
			res := fn(Edge{
				Start: g.vertices[start],
				End:   g.vertices[end],
				Cost:  cost,
			})
			if !res {
				break
			}
		}
	}

	if key != Root {
		visit(key.key)
		return
	}

	for key := range g.edges {
		visit(key)
	}
}

// VisitVertices uses fn to visit each vertex, starting at the vertex key.
func (g *Graph) VisitVertices(key Key, fn VertexVisitorFunc) {
	var (
		visited = make(map[internal.Key]struct{})
		queue   []Vertex
	)

	g.mtx.Lock()
	defer g.mtx.Unlock()

	if key != Root {
		queue = append(queue, g.vertices[key.key])
	} else {
		for _, node := range g.vertices {
			queue = append(queue, node)
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if _, skip := visited[node.key]; skip {
			continue
		}
		visited[node.key] = struct{}{}

		if !fn(node) {
			break
		}

		for edge := range g.edges[node.key] {
			if _, skip := visited[edge]; skip {
				continue
			}

			queue = append(queue, g.vertices[edge])
		}
	}
}

func (g *Graph) cloneUnsafe() *Graph {
	clone := &Graph{
		lastKey:  g.lastKey,
		vertices: make(map[internal.Key]Vertex, len(g.vertices)),
		edges:    make(map[internal.Key]map[internal.Key]int, len(g.edges)),
		redges:   make(map[internal.Key]map[internal.Key]int, len(g.redges)),
	}

	for key, vertex := range g.vertices {
		clone.vertices[key] = vertex
	}

	for src, edges := range g.edges {
		newedges := make(map[internal.Key]int, len(edges))
		for key, cost := range edges {
			newedges[key] = cost
		}
		clone.edges[src] = newedges
	}

	for src, edges := range g.redges {
		newedges := make(map[internal.Key]int, len(edges))
		for key, cost := range edges {
			newedges[key] = cost
		}
		clone.redges[src] = newedges
	}

	clone.config.undirected = g.config.undirected

	return clone
}

func (g *Graph) deleteEdgeUnsafe(from Key, to Key) {
	if from == to {
		return
	}

	if from != Any && to != Any {
		delete(g.edges[from.key], to.key)
		delete(g.redges[to.key], from.key)
		return
	}

	if from == Any {
		if _, ok := g.vertices[to.key]; !ok {
			return
		}

		for from := range g.redges[to.key] {
			delete(g.edges[from], to.key)
			delete(g.redges[to.key], from)
		}
	} else {
		if _, ok := g.vertices[from.key]; !ok {
			return
		}

		for to := range g.edges[from.key] {
			delete(g.edges[from.key], to)
			delete(g.redges[to], from.key)
		}
	}
}

func (g *Graph) getEdgesUnsafe(key internal.Key) map[internal.Key]int {
	edges, ok := g.edges[key]
	if !ok {
		edges = make(map[internal.Key]int)
		g.edges[key] = edges
	}

	return edges
}

func (g *Graph) getReverseEdgesUnsafe(
	key internal.Key,
) map[internal.Key]int {
	edges, ok := g.redges[key]
	if !ok {
		edges = make(map[internal.Key]int)
		g.redges[key] = edges
	}

	return edges
}
