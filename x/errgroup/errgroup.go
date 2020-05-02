// Package errgroup is a minor refactor of golang.org/x/sync/errgroup to better
// enable context propagation and short-circuiting of managed goroutines.
package errgroup

import (
	"context"
	"sync"
)

// A GroupableFunc is a function that is able to be processed as part of a
// Group. Such functions must respect the live-ness of the provided Context.
type GroupableFunc = func(context.Context) error

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
type Group struct {
	cancel context.CancelFunc
	ctx    context.Context
	err    error
	once   sync.Once
	wg     sync.WaitGroup
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	return &Group{
		cancel: cancel,
		ctx:    ctx,
	}, ctx
}

// Go calls the given function in a new goroutine. Critically, functions passed
// to Go must respect the provided Context and abort their execution if it is
// canceled.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(fn GroupableFunc) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := fn(g.ctx); err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()

	return g.err
}
