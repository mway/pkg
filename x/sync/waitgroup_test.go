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

package sync_test

import (
	"testing"

	"github.com/mway/pkg/x/sync"
	"github.com/stretchr/testify/require"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		require.Equal(t, 1, wg.Len())

		wg.Done()
		require.Equal(t, 0, wg.Len())
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		require.Equal(t, i+1, wg.Len())
	}

	var (
		ch   = make(chan struct{})
		done = func() bool {
			select {
			case <-ch:
				return true
			default:
			}
			return false
		}
	)

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	require.False(t, done())

	for i := wg.Len(); i > 0; i-- {
		wg.Done()
		require.Equal(t, i-1, wg.Len())
	}

	<-ch
	require.True(t, done())
}
