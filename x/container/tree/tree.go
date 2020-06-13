// Copyright (c) 2019 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
