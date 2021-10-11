package canceler

import (
	"sync"
)

// Canceler is a safe, idempotent type for consistently propagating cancellation
// signals via the io.Canceler interface. It is intended as a lightweight
// replacement for holding hierarchical context.Context objects. A Canceler may
// not be copied.
//
// The zero value for Canceler is valid and ready for use.
type Canceler struct {
	ch     chan struct{}
	mu     sync.RWMutex
	parent <-chan struct{}
}

// New returns a new Canceler that is a child of parent. If parent is not nil,
// the new Canceler is appended to the parent hierarchy in a chain-like fashion.
// Notably, if parent cancels, the new Canceler will as well; however, if the
// new Canceler cancels, the parent is unaffected.
func New(parent *Canceler) *Canceler {
	var ch <-chan struct{}
	if parent != nil {
		ch = parent.C()
	}

	return &Canceler{
		parent: ch,
	}
}

// C returns a wait channel. Both Canceler and the channel returned by C() may
// be shared across API boundaries and goroutines.
func (c *Canceler) C() <-chan struct{} {
	var ch <-chan struct{}

	c.mu.RLock()
	if c.ch == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		if c.ch == nil {
			c.ch = make(chan struct{})
			ch = c.ch

			go func() {
				c.mu.RLock()
				done := c.ch
				c.mu.RUnlock()

				select {
				case <-c.parent:
					close(done)
				case <-done:
				}
			}()
		}
		c.mu.Unlock()
	} else {
		ch = c.ch
		c.mu.RUnlock()
	}

	return ch
}

// Cancel causes the channel returned by C() to unblock, signaling receivers to
// close.
func (c *Canceler) Cancel() {
	var ch chan struct{}

	c.mu.RLock()
	if c.ch == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		c.ch = make(chan struct{})
		ch = c.ch
		close(ch)
		c.mu.Unlock()
	} else {
		ch = c.ch
		c.mu.RUnlock()
	}

	select {
	case <-ch:
	default:
		close(ch)
	}
}
