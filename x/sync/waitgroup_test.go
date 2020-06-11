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
