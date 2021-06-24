package compute // import "github.com/docker/docker/pkg/compute"

import (
	"context"
)

// TODO(rvolosatovs): Refactor into result[T] once support for generics lands.

// result is a result of a computation.
type result struct {
	Value interface{}
	Error error
}

// Func represents a computation producing a value or an error.
type Func func(context.Context) (interface{}, error)

// NewSingleton returns a new singleton f ready to use.
func NewSingleton(f Func) *Singleton {
	return &Singleton{
		callCh:   make(chan struct{}, 1),
		resultCh: make(chan result),
		f:        f,
	}
}

// Singleton ensures Func is executed by at most one goroutine at a time and
// propagates the result to all simultaneous callers.
// The intended use case is performing slow computations, e.g. doing I/O and sharing the results.
type Singleton struct {
	callCh   chan struct{}
	resultCh chan result
	f        Func
}

// Do calls the Func singleton is initialized with and returns the result.
// If the computation is already being performed by another goroutine, it will block
// until either the result is received from the other goroutine or context is done.
func (s Singleton) Do(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case s.callCh <- struct{}{}:
		// Lock acquired - perform computation.
		defer func() {
			<-s.callCh // Release lock.
		}()

	case res := <-s.resultCh:
		// Another goroutine computed the result - return.
		return res.Value, res.Error
	}

	var res result
	v, err := s.f(ctx)
	if err != nil {
		res = result{
			Error: err,
		}
	} else {
		res = result{
			Value: v,
		}
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case s.resultCh <- res:
			// Push computation result to other goroutines calling the function, if any.

		default:
			return res.Value, res.Error
		}
	}
}
