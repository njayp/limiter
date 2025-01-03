package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

type Limiter struct {
	limit      int
	interval   time.Duration
	tokens     chan struct{}
	activeJobs *atomic.Int32
}

// NewLimiter creates a new Limiter that will start a maximum of limit jobs per interval. Job start order is not guaranteed.
func NewLimiter(limit int, interval time.Duration) *Limiter {
	tokens := make(chan struct{}, limit)

	// fill the token bucket
	for range limit {
		tokens <- struct{}{}
	}

	return &Limiter{
		limit:      limit,
		interval:   interval,
		tokens:     tokens,
		activeJobs: new(atomic.Int32),
	}
}

// Run runs the given function if the limit has not been reached. Run blocks, it is designed to run within a go routine.
func (l *Limiter) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	// block until a token is available
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
	}

	// release the token after the duration has passed
	go func() {
		time.Sleep(l.interval)
		l.tokens <- struct{}{}
	}()

	l.activeJobs.Add(1)
	defer func() {
		l.activeJobs.Add(-1)
	}()

	return fn(ctx)
}

func (l *Limiter) ActiveJobs() int32 {
	return l.activeJobs.Load()
}
