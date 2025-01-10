package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

type Runner struct {
	limit       int
	interval    time.Duration
	tokens      chan struct{}
	tokenReturn chan struct{}
	activeJobs  atomic.Int32
}

// NewRunner creates a new Limiter that will start a maximum of limit jobs per interval. Jobs are staggered by stagger duration. Job start order is not guaranteed.
func NewRunner(ctx context.Context, limit int, interval, stagger time.Duration) *Runner {
	tokens := make(chan struct{})
	tokenReturn := make(chan struct{}, limit)

	// fill the token bucket
	for range limit {
		tokenReturn <- struct{}{}
	}

	// stagger the token return
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case tokens <- <-tokenReturn:
				time.Sleep(stagger)
			}
		}
	}()

	return &Runner{
		limit:       limit,
		interval:    interval,
		tokens:      tokens,
		tokenReturn: tokenReturn,
		activeJobs:  atomic.Int32{},
	}
}

// Run runs the given function if the limiting parameters are met. Run blocks, it is designed to run within a go routine.
func (r *Runner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	// block until a token is available
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.tokens:
	}

	// release the token after the duration has passed
	go func() {
		time.Sleep(r.interval)
		r.tokenReturn <- struct{}{}
	}()

	r.activeJobs.Add(1)
	defer func() {
		r.activeJobs.Add(-1)
	}()

	return fn(ctx)
}

func (r *Runner) ActiveJobs() int32 {
	return r.activeJobs.Load()
}
