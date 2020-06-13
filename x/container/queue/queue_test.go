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
