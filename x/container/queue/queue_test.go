package queue_test

import (
	"strconv"
	"testing"

	"github.com/mway/pkg/x/container/queue"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
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
				queue queue.Queue
			)

			require.Equal(t, nil, queue.Front())
			require.Equal(t, nil, queue.Back())
			require.Equal(t, nil, queue.Dequeue())

			for i := 0; i < num; i++ {
				queue.Push(tc.values[i])
				require.Equal(t, tc.values[0], queue.Front())
				require.Equal(t, tc.values[i], queue.Back())
			}

			require.Equal(t, len(tc.values), queue.Len())

			sz := queue.Len()
			for i := 0; i < sz; i++ {
				require.Equal(t, tc.values[i], queue.Dequeue())
			}

			require.Equal(t, 0, queue.Len())
		})
	}
}
