// Copyright (c) 2020 Matt Way
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

package nary

import (
	"fmt"
	"sync/atomic"

	"github.com/mway/pkg/x/container/tree"
)

var (
	_ tree.Tree = (*Node)(nil)
)

// Node is an n-ary tree node.
type Node struct {
	tree.Node

	idx    map[tree.Key]*Node
	k      uint
	lastID uint64

	root     *Node
	prev     *Node
	next     *Node
	parent   *Node
	children *Node
}

// NewTree creates a new n-ary tree of order k with the given root value.
func NewTree(k uint, value interface{}) *Node {
	node := &Node{
		Node: tree.Node{
			Key:   tree.NewKey(1),
			Value: value,
		},
		k:      k,
		lastID: 1,
	}

	node.idx = map[tree.Key]*Node{
		node.Node.Key: node,
	}

	return node
}

// Delete deletes key from the current subtree.
func (n *Node) Delete(key tree.Key) {
	subtree, ok := n.getNode(key)
	if !ok {
		return
	}

	subtree.Iterate(tree.PostOrder, func(cur tree.Node) bool {
		node, ok := n.getNode(cur.Key)
		if !ok {
			return false
		}

		n.unindex(node.Key)

		if node.prev != nil {
			node.prev.next = node.next
		}

		if node.next != nil {
			node.next.prev = node.prev
		}

		if node.parent != nil && node.parent.children == node {
			node.parent.children = node.next
		}

		return false
	})
}

// Insert inserts value into the tree. Inserts attempt to fill the tree, thus
// the lowest-level (closest to root), leftmost child will be used for appending
// such that that child will have no more than k children.
func (n *Node) Insert(value interface{}) (key tree.Key) {
	if uint(n.numChildren()) < n.k {
		key = n.addChild(value)
		return
	}

	n.Iterate(tree.LevelOrder, func(cur tree.Node) bool {
		node, ok := n.getNode(cur.Key)
		if !ok {
			return false
		}

		if uint(node.numChildren()) >= n.k {
			return true
		}

		key = node.addChild(value)
		return false
	})

	return
}

// Iterate iterates over n using the given order and handler.
func (n *Node) Iterate(
	order tree.TraversalOrder,
	handler tree.NodeHandlerFunc,
) {
	switch order {
	case tree.LevelOrder:
		n.traverseLevelOrder(handler)
	case tree.PostOrder:
		n.traversePostOrder(handler)
	case tree.PreOrder:
		n.traversePreOrder(handler)
	}
}

// K returns the k order of the tree.
func (n *Node) K() uint {
	return n.k
}

// Root returns the root of the tree.
func (n *Node) Root() tree.Key {
	return n.getRoot().Node.Key
}

func (n *Node) addChild(value interface{}) (key tree.Key) {
	if children := n.numChildren(); uint(children) >= n.k {
		panic(fmt.Sprintf("cannot add child to node with %d children", children))
	}

	key = n.newKey()
	node := &Node{
		Node: tree.Node{
			Key:   key,
			Value: value,
		},
		k:      n.k,
		parent: n,
		root:   n.getRoot(),
	}

	n.index(key, node)

	if n.children == nil {
		n.children = node
		return
	}

	child := n.children
	for {
		if child.next == nil {
			node.prev = child
			child.next = node
			break
		}

		child = child.next
	}

	return
}

func (n *Node) getRoot() *Node {
	if n.root == nil {
		return n
	}

	return n.root
}

func (n *Node) getNode(key tree.Key) (*Node, bool) {
	node, ok := n.getRoot().idx[key]
	return node, ok
}

func (n *Node) index(key tree.Key, node *Node) {
	n.getRoot().idx[key] = node
}

func (n *Node) newKey() tree.Key {
	root := n.getRoot()
	return tree.NewKey(atomic.AddUint64(&root.lastID, 1))
}

func (n *Node) numChildren() (num int) {
	child := n.children
	for child != nil {
		num++
		child = child.next
	}

	return
}

func (n *Node) unindex(key tree.Key) {
	delete(n.getRoot().idx, key)
}

func (n *Node) traverseLevelOrder(handler tree.NodeHandlerFunc) {
	queue := []*Node{n}
	for len(queue) > 0 {
		num := len(queue)

		for i := 0; i < num; i++ {
			node := queue[0]
			queue = queue[1:]

			if node == nil {
				continue
			}

			if !handler(node.Node) {
				return
			}

			for child := node.children; child != nil; child = child.next {
				queue = append(queue, child)
			}
		}
	}
}

func (n *Node) traversePostOrder(handler tree.NodeHandlerFunc) {
	for child := n.children; child != nil; child = child.next {
		child.traversePostOrder(handler)
	}

	handler(n.Node)
}

func (n *Node) traversePreOrder(handler tree.NodeHandlerFunc) {
	handler(n.Node)

	for child := n.children; child != nil; child = child.next {
		child.traversePostOrder(handler)
	}
}
