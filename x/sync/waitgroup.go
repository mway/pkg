package sync

import (
	"sync"
	"sync/atomic"
)

// WaitGroup is a wrapper around sync.WaitGroup that adds a length component.
// It is a drop-in replacement that is functionally equivalent in every way
// except that it also tracks the value of the underlying WaitGroup counter.
type WaitGroup struct {
	wg sync.WaitGroup
	n  uint64
}

// Add adds delta, which may be negative, to the WaitGroup counter. If the
// counter becomes zero, all goroutines blocked on Wait are released. If the
// counter goes negative, Add panics.
func (g *WaitGroup) Add(delta int) {
	atomic.AddUint64(&g.n, uint64(delta))
	g.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (g *WaitGroup) Done() {
	atomic.AddUint64(&g.n, ^uint64(0))
	g.wg.Done()
}

// Len returns the current value of the underlying WaitGroup counter.
func (g *WaitGroup) Len() int {
	return int(atomic.LoadUint64(&g.n))
}

// Wait blocks until the WaitGroup counter is zero.
func (g *WaitGroup) Wait() {
	g.wg.Wait()
}
