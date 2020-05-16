package graph

import (
	"fmt"

	"github.com/mway/pkg/x/container/graph/internal"
)

// Vertex is a point node in a graph.
type Vertex struct {
	key   internal.Key
	graph *Graph
	value interface{}
}

// Key provides the key for v.
func (v *Vertex) Key() Key {
	if v == nil {
		return _zeroKey
	}

	return newKey(v.key)
}

// String returns a string representation of v.
func (v *Vertex) String() string {
	if v == nil {
		return ""
	}

	if str, ok := v.value.(fmt.Stringer); ok {
		return str.String()
	}

	return fmt.Sprintf("%v", v.value)
}

// Value returns the value held in v.
func (v *Vertex) Value() interface{} {
	if v == nil {
		return nil
	}

	return v.value
}
