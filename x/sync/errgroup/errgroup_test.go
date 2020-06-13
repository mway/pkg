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
