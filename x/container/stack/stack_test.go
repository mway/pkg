package stack_test

import (
	"strconv"
	"testing"

	"github.com/mway/pkg/x/container/stack"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	cases := []struct {
		values []interface{}
	}{
		{
			values: []interface{}{1, 2, 3, 4, 5},
		},
	}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				num   = len(tc.values)
				stack stack.Stack
			)

			for i := 0; i < num; i++ {
				stack.Push(tc.values[i])
				require.Equal(t, tc.values[i], stack.Top())
			}

			require.Equal(t, len(tc.values), stack.Len())

			for i := stack.Len() - 1; i >= 0; i-- {
				require.Equal(t, tc.values[i], stack.Pop())
			}

			require.Equal(t, 0, stack.Len())
		})
	}
}
