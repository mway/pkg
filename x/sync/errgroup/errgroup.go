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

// Package errgroup is a minor refactor of golang.org/x/sync/errgroup to better
// enable context propagation and short-circuiting of managed goroutines.
package errgroup

import (
	"context"
	"sync"

	"go.uber.org/multierr"
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

// MultiGroup is a Group variant that accumulates errors rather than simply
// returning the first encountered. Semantically, this affords callers the
// ability to determine whether or not launched goroutines should abort
// immediately once the provided context is cancelled. If callers choose not
// to respect the context, it is their responsibility to ensure that goroutines
// do not leak.
type MultiGroup struct {
	group Group
	mtx   sync.Mutex
}

// MultiWithContext returns a new MultiGroup and an associated Context derived
// from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func MultiWithContext(ctx context.Context) (*MultiGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	return &MultiGroup{
		group: Group{
			cancel: cancel,
			ctx:    ctx,
		},
	}, ctx
}

// Go calls the given function in a new goroutine. Callers may choose whether or
// not to respect the context (in order to have multiple launched routines'
// errors combined), however it is the caller's responsibility to avoid
// goroutine leaks in such cases.
func (g *MultiGroup) Go(fn GroupableFunc) {
	g.group.wg.Add(1)

	go func() {
		defer g.group.wg.Done()

		if err := fn(g.group.ctx); err != nil {
			g.mtx.Lock()
			defer g.mtx.Unlock()

			// TODO(mway): Consider using a builtin wrapping convention rather
			//             than multierr.
			g.group.err = multierr.Combine(g.group.err, err)
			g.group.cancel()
		}
	}()
}

// Wait blocks until all function calls from the Go method have returned, then
// returns all non-nil errors (if any) from them, combined via multierr.
func (g *MultiGroup) Wait() error {
	return g.group.Wait()
}
