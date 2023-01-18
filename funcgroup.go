package sync

import (
	"context"
	"sync"
)

// TODO: fix documentation

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// FuncGroup is safe to copy.
type FuncGroup interface {
	Wait() error
	Go(f func(context.Context) error)
	Cancel()
}

type funcGroup struct {
	wait       chan bool
	cancelFunc func()
	ctx        context.Context

	mu         sync.Mutex
	cancelOnce sync.Once
	errOnce    sync.Once
	err        error
}

func WithNewContext() (FuncGroup, context.Context) {
	return WithContext(context.Background())
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled if parent context gets cancelled or when first function passed to Go() returns.
func WithContext(ctx context.Context) (FuncGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &funcGroup{wait: make(chan bool), ctx: ctx, cancelFunc: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
// Calling wait without any Go() calls will block forever
// Wait returns nil error when FuncGroup is explictly cancelled by Cancel()
func (g *funcGroup) Wait() error {
	<-g.wait

	// if ok && g.cancelFunc != nil {
	// 	g.cancelFunc()
	// }
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.err
}

// Cancels the context once and unblocks Wait immediately. Goroutines must handle context cancellation
func (g *funcGroup) Cancel() {
	g.cancelOnce.Do(func() {
		g.cancelFunc()
		close(g.wait)
	})
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext.
// returns error from first returning function.
func (g *funcGroup) Go(f func(context.Context) error) {

	go func() {

		err := f(g.ctx)
		g.errOnce.Do(func() {
			g.mu.Lock()
			defer g.mu.Unlock()
			g.err = err
		})

		g.Cancel()
	}()
}
