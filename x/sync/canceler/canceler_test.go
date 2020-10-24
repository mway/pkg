package canceler_test

import (
	"context"
	"testing"
	"time"

	"github.com/mway/pkg/x/sync/canceler"
	"github.com/stretchr/testify/require"
)

func TestCanceler(t *testing.T) {
	var (
		c canceler.Canceler

		ch1 = c.C()
		ch2 = c.C()
		ch3 = c.C()
	)

	require.Equal(t, ch1, ch2, "C() returned different channels")
	require.Equal(t, ch2, ch3, "C() returned different channels")

	c.Cancel()

	for _, ch := range []<-chan struct{}{ch1, ch2, ch3} {
		select {
		case <-ch:
		default:
			require.Fail(t, "cancellation channel not closed")
		}
	}
}

func TestCancelerEarlyCancel(t *testing.T) {
	var c canceler.Canceler
	c.Cancel()

	select {
	case <-c.C():
	default:
		require.Fail(t, "channel not closed after early cancellation")
	}
}

func TestCancelerInheritance(t *testing.T) {
	var (
		canceler1 = canceler.New(nil)
		canceler2 = canceler.New(canceler1)
		canceler3 = canceler.New(canceler2)
		canceler4 = canceler.New(canceler1)
	)

	canceler4.Cancel()

	select {
	case <-canceler1.C():
		require.Fail(t, "child canceler canceled parent")
	default:
	}

	select {
	case <-canceler2.C():
		require.Fail(t, "sibling canceler canceled other sibling")
	default:
	}

	canceler1.Cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cancelers := []*canceler.Canceler{canceler1, canceler2, canceler3, canceler4}
	for _, canceler := range cancelers {
		select {
		case <-canceler.C():
		case <-ctx.Done():
			require.Fail(t, "ancestor canceler did not cancel descendant")
		}
	}
}
