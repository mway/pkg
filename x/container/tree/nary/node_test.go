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

package nary_test

import (
	"strconv"
	"testing"

	"github.com/mway/pkg/x/container/tree"
	"github.com/mway/pkg/x/container/tree/nary"
	"github.com/stretchr/testify/require"
)

func TestNodeTraversal(t *testing.T) {
	cases := []struct {
		k        uint
		values   []interface{}
		order    tree.TraversalOrder
		expected []interface{}
	}{
		{
			k:        3,
			values:   []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			order:    tree.LevelOrder,
			expected: []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			k:        3,
			values:   []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			order:    tree.PostOrder,
			expected: []interface{}{5, 6, 7, 2, 8, 9, 10, 3, 4, 1},
		},
		{
			k:        3,
			values:   []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			order:    tree.PreOrder,
			expected: []interface{}{1, 5, 6, 7, 2, 8, 9, 10, 3, 4},
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				root   = newTree(tc.k, tc.values)
				values []interface{}
			)

			root.Iterate(tc.order, func(node tree.Node) bool {
				values = append(values, node.Value)
				return true
			})

			require.Equal(t, tc.expected, values)
		})
	}
}

func TestNodeDeletion(t *testing.T) {
	var (
		root  = nary.NewTree(3, 1)
		_     = root.Insert(2)
		key3  = root.Insert(3)
		_     = root.Insert(4)
		_     = root.Insert(5)
		_     = root.Insert(6)
		_     = root.Insert(7)
		_     = root.Insert(8)
		_     = root.Insert(9)
		key10 = root.Insert(10)
		total = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	)

	getValues := func() []interface{} {
		var values []interface{}
		root.Iterate(tree.LevelOrder, func(node tree.Node) bool {
			values = append(values, node.Value)
			return true
		})
		return values
	}

	require.Equal(t, total, getValues())

	// n.b. No effect.
	root.Delete(tree.Key{})

	root.Delete(key10)
	require.Equal(t, total[:9], getValues())

	root.Delete(key3)
	require.Equal(t, []interface{}{1, 2, 4, 5, 6, 7}, getValues())
}

func newTree(k uint, values []interface{}) *nary.Node {
	if len(values) == 0 {
		return nil
	}

	root := nary.NewTree(k, values[0])
	for _, value := range values[1:] {
		root.Insert(value)
	}

	return root
}
