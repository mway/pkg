package tree

type (
	// A NodeTraversalFunc is used to walk a Node subtree with a
	// NodeHandlerFunc.
	NodeTraversalFunc = func(Node, NodeHandlerFunc)

	// A NodeHandlerFunc is used to evaluate a Node while walking a tree.
	NodeHandlerFunc = func(Node) bool
)

// TraversalOrder determines tree traversal order.
type TraversalOrder string

// Available traversal orders.
const (
	LevelOrder   TraversalOrder = "level"
	PreOrder     TraversalOrder = "pre"
	PostOrder    TraversalOrder = "post"
	DefaultOrder                = LevelOrder
)

// Key is an identity type used to reference tree nodes.
type Key struct {
	n uint64
}

// NewKey creates a new Key based on n.
func NewKey(n uint64) Key {
	return Key{
		n: n,
	}
}

// Node is a tree node.
type Node struct {
	Key   Key
	Value interface{}
}

// Tree is a tree-like interfaces.
type Tree interface {
	Delete(Key)
	Insert(interface{}) Key
	Iterate(TraversalOrder, NodeHandlerFunc)
	K() uint
	Root() Key
}
