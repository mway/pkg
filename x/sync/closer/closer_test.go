package closer_test

import (
	"context"
	"testing"
	"time"

	"github.com/mway/pkg/x/sync/closer"
	"github.com/stretchr/testify/require"
)

func TestCloser(t *testing.T) {
	var (
		closer closer.Closer

		ch1 = closer.C()
		ch2 = closer.C()
		ch3 = closer.C()
	)

	require.Equal(t, ch1, ch2, "C() returned different channels")
	require.Equal(t, ch2, ch3, "C() returned different channels")

	require.NoError(t, closer.Close())

	for _, ch := range []<-chan struct{}{ch1, ch2, ch3} {
		select {
		case <-ch:
		default:
			require.Fail(t, "signal channel not closed")
		}
	}
}

func TestCloserEarlyClose(t *testing.T) {
	var c closer.Closer
	require.NoError(t, c.Close())

	select {
	case <-c.C():
	default:
		require.Fail(t, "channel not closed after early close")
	}
}

func TestCloserInheritance(t *testing.T) {
	var (
		closer1 = closer.New(nil)
		closer2 = closer.New(closer1)
		closer3 = closer.New(closer2)
		closer4 = closer.New(closer1)
	)

	require.NoError(t, closer4.Close())

	select {
	case <-closer1.C():
		require.Fail(t, "child closer closed parent")
	default:
	}

	select {
	case <-closer2.C():
		require.Fail(t, "sibling closer closed other sibling")
	default:
	}

	require.NoError(t, closer1.Close())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	closers := []*closer.Closer{closer1, closer2, closer3, closer4}
	for _, closer := range closers {
		select {
		case <-closer.C():
		case <-ctx.Done():
			require.Fail(t, "ancestor closer did not close descendant")
		}
	}
}
