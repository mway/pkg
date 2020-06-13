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

package errgroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mway/pkg/x/sync/errgroup"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	var (
		calls uint32 = 0
		funcs        = []errgroup.GroupableFunc{
			func(ctx context.Context) error {
				atomic.AddUint32(&calls, 1)
				return nil
			},
			func(ctx context.Context) error {
				atomic.AddUint32(&calls, 1)
				return nil
			},
		}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	for _, fn := range funcs {
		group.Go(fn)
	}

	select {
	case <-ctx.Done():
		require.Fail(t, "context must not be done before Wait()")
	default:
	}

	err := group.Wait()
	require.NoError(t, err)
	require.Equal(t, uint32(2), calls)

	select {
	case <-ctx.Done():
	default:
		require.Fail(t, "context must be done after Wait()")
	}
}

func TestGroupErrors(t *testing.T) {
	var (
		calls uint32 = 0
		funcs        = []errgroup.GroupableFunc{
			func(ctx context.Context) error {
				return errors.New("error")
			},
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(time.Hour):
				}

				atomic.AddUint32(&calls, 1)
				return nil
			},
		}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	for _, fn := range funcs {
		group.Go(fn)
	}

	err := group.Wait()

	select {
	case <-ctx.Done():
	default:
		require.Fail(t, "context must be done after Wait()")
	}

	require.Error(t, err)
	require.Equal(t, "error", err.Error())
	require.Equal(t, uint32(0), calls)
}
