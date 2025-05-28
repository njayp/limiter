package concurrent

import (
	"context"
)

type Semaphore struct {
	bucket chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		bucket: make(chan struct{}, limit),
	}
}

func (s *Semaphore) Wait(ctx context.Context) (*Token, error) {
	// wait for a slot in the semaphore
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.bucket <- struct{}{}:
		// give the caller a token that can be released when done
		return NewToken(func() {
			<-s.bucket
		}), nil
	}
}

func (s *Semaphore) Len() int {
	return len(s.bucket)
}
