package limiter

import (
	"context"
	"time"

	"github.com/njayp/limiter/concurrent"
	"github.com/njayp/limiter/rate"
)

type Limiter struct {
	rateLimiter *rate.Limiter
	semaphore   *concurrent.Semaphore
}

func NewLimiter(limit int, interval time.Duration, concurrentLimit int) *Limiter {
	return &Limiter{
		rateLimiter: rate.NewLimiter(limit, interval),
		semaphore:   concurrent.NewSemaphore(concurrentLimit),
	}
}

func (l *Limiter) Wait(ctx context.Context) (*concurrent.Token, error) {
	token, err := l.semaphore.Wait(ctx)
	if err != nil {
		return nil, err
	}

	err = l.rateLimiter.Wait(ctx)
	if err != nil {
		token.Release()
		return nil, err
	}

	return token, nil
}

func (l *Limiter) Close() {
	l.rateLimiter.Close()
}
