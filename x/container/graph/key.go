package graph

import (
	"github.com/mway/pkg/x/container/graph/internal"
)

// Key is a thin identifier for graph vertices.
type Key struct {
	key internal.Key
}

func newKey(key internal.Key) Key {
	return Key{
		key: key,
	}
}
