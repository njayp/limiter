package limiter

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// use NewRunner()
type Runner struct {
	limit       int
	interval    time.Duration
	token       <-chan struct{}
	tokenBucket chan<- struct{}
	close       chan struct{}
	activeJobs  atomic.Int32
}

// NewRunner creates a new Runner that starts a maximum of limit jobs per interval, and staggers the start of jobs
func NewRunner(limit int, interval, stagger time.Duration) *Runner {
	token := make(chan struct{})
	tokenBucket := make(chan struct{}, limit)
	close := make(chan struct{})

	// fill the token bucket
	for range limit {
		tokenBucket <- struct{}{}
	}

	// stagger the release of tokens
	go func() {
		for {
			select {
			case <-close:
				return
			case token <- <-tokenBucket:
				time.Sleep(stagger)
			}
		}
	}()

	return &Runner{
		limit:       limit,
		interval:    interval,
		token:       token,
		tokenBucket: tokenBucket,
		activeJobs:  atomic.Int32{},
		close:       close,
	}
}

// Run runs the given function if the limiting parameters are met. Run is thread-safe. It blocks until the function returns.
func (r *Runner) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	// block until a token is available
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.close:
		return fmt.Errorf("limiter.Runner closed")
	case <-r.token:
	}

	// return the token after the duration has passed
	go func() {
		select {
		case <-r.close:
			return
		case <-time.After(r.interval):
			r.tokenBucket <- struct{}{}
		}
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

// Close cleans up any lingering goroutines. No more jobs can be started after Close is called
func (r *Runner) Close() {
	close(r.close)
}
