package limiter

import (
	"context"
	"fmt"
	"time"
)

// use NewLimiter()
type Limiter struct {
	limit       int
	interval    time.Duration
	token       <-chan struct{}
	tokenBucket chan<- struct{}
	close       chan struct{}
}

// NewLimiter impl a token bucket that will allow a maximum number of function calls within a given interval.
func NewLimiter(limit int, interval, stagger time.Duration) *Limiter {
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

	return &Limiter{
		limit:       limit,
		interval:    interval,
		token:       token,
		tokenBucket: tokenBucket,
		close:       close,
	}
}

// Wait blocks until a token is available.
func (r *Limiter) Wait(ctx context.Context) error {
	// block until a token is available
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.close:
		return fmt.Errorf("limiter.Runner closed")
	case <-r.token:
		go r.returnToken()
		return nil
	}
}

// Try returns immediately with an error if a token is not available. If a token is available, it is consumed and nil is returned.
func (r *Limiter) Try(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.close:
		return fmt.Errorf("limiter.Runner closed")
	case <-r.token:
		go r.returnToken()
		return nil
	// return an error if no token is available
	default:
		return fmt.Errorf("no token available")
	}
}

func (r *Limiter) returnToken() {
	select {
	case <-r.close:
		return
	case <-time.After(r.interval):
		r.tokenBucket <- struct{}{}
	}
}

// Close cleans up any lingering goroutines. No more jobs can be started after Close is called
func (r *Limiter) Close() {
	close(r.close)
}
