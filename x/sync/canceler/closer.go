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
	mtx    sync.Mutex
	parent <-chan struct{}
}

// New returns a new Canceler that is a child of parent. If parent is not nil,
// the new Canceler is appended to the parent hierarchy in a chain-like fashion.
// Notably, if parent cancels, the new Canceler will as well; however, if the
// new Canceler cancels, the parent is unaffected.
func New(parent *Canceler) *Canceler {
	c := &Canceler{}

	if parent != nil {
		c.parent = parent.C()
	}

	return c
}

// C returns a wait channel. Both Canceler and the channel returned by C() may
// be shared across API boundaries and goroutines.
func (c *Canceler) C() <-chan struct{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.ch == nil {
		c.ch = make(chan struct{})

		go func() {
			select {
			case <-c.parent:
				close(c.ch)
			case <-c.ch:
			}
		}()
	}

	return c.ch
}

// Cancel causes the channel returned by C() to unblock, signaling receivers to
// close.
func (c *Canceler) Cancel() {
	c.mtx.Lock()
	if c.ch == nil {
		c.ch = make(chan struct{})
		close(c.ch)
	}
	c.mtx.Unlock()

	select {
	case <-c.ch:
	default:
		close(c.ch)
	}
}
