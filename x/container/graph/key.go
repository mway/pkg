package graph

import (
	"math/rand"
	"time"

	"github.com/mway/pkg/x/container/graph/internal"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Key is a thin identifier for graph vertices.
type Key struct {
	key internal.Key
}

func newKey(key internal.Key) Key {
	return Key{
		key: key,
	}
}

func newRandomKey() Key {
	return Key{
		key: internal.Key(rand.Uint64()),
	}
}
