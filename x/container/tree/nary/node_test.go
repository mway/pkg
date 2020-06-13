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
