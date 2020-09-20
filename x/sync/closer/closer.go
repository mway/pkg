package closer

import (
	"io"
	"sync"
)

var (
	_ io.Closer = (*Closer)(nil)
)

// Closer is a safe, idempotent type for consistently propagating cancellation
// signals via the io.Closer interface. It is intended as a lightweight
// replacement for holding hierarchical context.Context objects. A Closer may
// not be copied.
//
// The zero value for Closer is valid and ready for use.
type Closer struct {
	ch     chan struct{}
	mtx    sync.Mutex
	parent <-chan struct{}
}

// New returns a new Closer
func New(parent *Closer) *Closer {
	c := &Closer{}

	if parent != nil {
		c.parent = parent.C()
	}

	return c
}

// C returns a wait channel. Both Closer and the channel returned by C() may be
// shared across API boundaries and goroutines.
func (c *Closer) C() <-chan struct{} {
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

// Close causes the channel returned by C() to unblock, signaling receivers to
// close.
func (c *Closer) Close() error {
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

	return nil
}
