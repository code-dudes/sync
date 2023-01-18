# sync

[![Go Report Card](https://goreportcard.com/badge/github.com/code-dudes/sync?style=flat-square)](https://goreportcard.com/report/github.com/code-dude/sync) ![Build](https://github.com/code-dudes/sync/actions/workflows/main.yml/badge.svg?branch=main&event=push)  [![PkgGoDev](https://pkg.go.dev/badge/github.com/code-dudes/sync)](https://pkg.go.dev/github.com/code-dudes/sync)   [![Release](https://img.shields.io/github/v/release/code-dudes/sync.svg?style=flat-square)](https://github.com/code-dudes/sync/releases/latest)

Sync go package provides enhanced version of some constructs available in golang's sync package.

</br>

Available Constructs:
1. **Once**  
Copy safe, stateful implementation of builtin sync.Once[idempotent]. Can be passed to functions and shared by go-routines.

2. **FuncGroup**
Launch a bunch of go routines which get cancelled once first worker returns. Slightly different from ErrGroup:
- Calls all goroutines along with a cancellable context.
- Context gets cancelled when first goroutine returns.
- Goroutines should return on context cancellation.

	```
	type FuncGroup interface {
		Wait() error
		Go(f func(context.Context) error)
		Cancel()
	}
	func WithNewContext() (FuncGroup, context.Context)
	func WithContext(ctx context.Context) (FuncGroup, context.Context)
	```
	**When to use**  
	Running parallel workers handling a single task, where if any of the worker fails or completes, other handlers should also exit.
