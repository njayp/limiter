package rate

import (
	"context"
	"time"
)

// use NewLimiter()
type Limiter struct {
	interval    time.Duration
	token       <-chan struct{}
	tokenBucket chan<- struct{}
	close       chan struct{}
}

// NewLimiter impl a token bucket that will allow a maximum number of function calls within a given interval.
func NewLimiter(count int, interval time.Duration) *Limiter {
	token := make(chan struct{})
	tokenBucket := make(chan struct{}, count)
	close := make(chan struct{})

	// fill the token bucket
	for range count {
		tokenBucket <- struct{}{}
	}

	// stagger the release of tokens
	go func() {
		for {
			select {
			case <-close:
				return
			case token <- <-tokenBucket:
				// wait a short time before releasing the next token to prevent DDOS
				// 50ms seems good for most use cases
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	return &Limiter{
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
		return LimiterClosedError
	case <-r.token:
		go r.returnToken()
		return nil
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
