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

			require.Equal(t, nil, stack.Top())
			require.Equal(t, nil, stack.Pop())

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
