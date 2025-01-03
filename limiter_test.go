package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	ctx := context.Background()
	limit := 5
	interval := time.Second
	l := NewLimiter(limit, interval)

	var wg sync.WaitGroup
	var count int32

	fn := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	start := time.Now()
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Run(ctx, fn)
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	if count != 10 {
		t.Errorf("expected 10 executions, got %d", count)
	}

	if duration < time.Second {
		t.Errorf("expected duration to be at least 1 seconds, got %v", duration)
	}
}

func TestLimiterActiveJobs(t *testing.T) {
	ctx := context.Background()
	limit := 2
	interval := time.Second
	l := NewLimiter(limit, interval)

	var wg sync.WaitGroup
	fn := func(ctx context.Context) error {
		time.Sleep(500 * time.Millisecond)
		return nil
	}

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Run(ctx, fn)
		}()
	}

	time.Sleep(100 * time.Millisecond)
	activeJobs := l.ActiveJobs()
	if activeJobs != 2 {
		t.Errorf("expected 2 active jobs, got %d", activeJobs)
	}

	wg.Wait()
}
