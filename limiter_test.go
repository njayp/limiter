package limiter

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	ctx := context.Background()
	limit := 2
	interval := time.Second
	stagger := 0 * time.Millisecond
	l := NewLimiter(limit, interval, stagger)
	wg := &sync.WaitGroup{}
	start := time.Now()

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Wait(ctx)
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	expectedDuration := 1 * time.Second
	if duration < expectedDuration {
		t.Errorf("expected duration to be at least %v, got %v", expectedDuration, duration)
	}
}

func TestRunnerStagger(t *testing.T) {
	ctx := context.Background()
	limit := 10
	interval := time.Millisecond
	stagger := time.Second
	l := NewLimiter(limit, interval, stagger)
	a := atomic.Int32{}

	for range 10 {
		go func() {
			l.Wait(ctx)
		}()
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	got := a.Load()
	expectedCount := int32(2)
	if got != expectedCount {
		t.Errorf("expected %d executions, got %d", expectedCount, got)
	}
}
